package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	ipn "github.com/jackdoe/gin-ipn"
)

const ROOT = "/mnt/home_attached"

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	pid := os.Getpid()
	counter := uint64(0)

	r := gin.Default()

	r.POST("/register/:user", func(c *gin.Context) {
		key, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		u := c.Param("user")
		err = addUser(u, key)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		requestDump, err := httputil.DumpRequest(c.Request, true)
		if err != nil {
			panic(err)
		}
		err = appendUserLog(u, "register.txt", []byte(requestDump))
		if err != nil {
			panic(err)
		}

		c.String(200, fmt.Sprintf(AFTER_REGISTER, u, u))
	})

	r.GET("/~:user", func(c *gin.Context) {
		c.Redirect(302, "/~"+c.Param("user")+"/")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(200, SLASH)
	})

	r.GET("/tos", func(c *gin.Context) {
		c.String(200, LICENSE)
	})

	r.GET("/thanks_for_paying", func(c *gin.Context) {
		c.String(200, THANKS_FOR_PAYING)
	})

	r.GET("/~:user/*path", func(c *gin.Context) {
		u := c.Param("user")
		rp := c.Param("path")

		err := userIsValid(u)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// cleanup golang http.ServeFile special handling of index.html
		if strings.HasSuffix(rp, "/index.html") {
			c.Redirect(302, "/~"+u+"/"+strings.TrimSuffix(rp, "/index.html"))
			return
		}

		local := path.Join(ROOT, u, "public_html")
		p := path.Join(local, filepath.Clean(rp))

		l, err := os.Readlink(p)
		if err == nil {
			p = l
		}

		// dont allow symlinks leading outside of home
		if !strings.HasPrefix(p, local) {
			c.String(418, "out of home")
			return
		}

		c.File(p)
	})

	r.POST("/mail/:user", func(c *gin.Context) {
		u := strings.Trim(c.Param("user"), "~")

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		uid, _, err := uidgid(u)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if uid < 1000 {
			c.JSON(400, gin.H{"error": "invalid user"})
			return
		}
		_, err = mail.ReadMessage(bytes.NewBuffer(body))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error(), "body": string(body)})
			return
		}

		basename := fmt.Sprintf("%v.M%vP%v_%v.%v", time.Now().Unix(), time.Now().Nanosecond()/1000, pid, atomic.AddUint64(&counter, 1), hostname)

		fn := path.Join(ROOT, u, "Maildir", "tmp", basename)
		f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		size, err := f.Write(body)
		if err != nil {
			f.Close()
			_ = os.Remove(fn)

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		f.Close()

		err = chown(u, fn)
		if err != nil {
			_ = os.Remove(fn)

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		newname := path.Join(ROOT, u, "Maildir", "new", fmt.Sprintf("%v,S=%v", basename, size))
		err = os.Rename(fn, newname)
		if err != nil {
			_ = os.Remove(fn)

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.String(200, "OK")
	})

	r.GET("/sub/:user", func(c *gin.Context) {
		u := c.Param("user")
		err := userIsValid(u)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		prefix := "https://www.paypal.com/cgi-bin/webscr"

		url := prefix + "?cmd=_xclick-subscriptions&business=jack%40baxx.dev&a3=1&p3=1&t3=M&item_name=berserk.red+-+personal+website&return=https%3A%2F%2Fberserk.red%2Fthanks_for_paying&a1=0.1&no_shipping=1&p1=1&t1=M&src=1&sra=1&no_note=1&no_note=1&currency_code=EUR&lc=GB&notify_url=https%3A%2F%2Fberserk.red%2Fipn%2F" + u
		c.Redirect(http.StatusFound, url)
	})

	ipn.Listener(r, "/ipn/:user", func(c *gin.Context, err error, body string, n *ipn.Notification) error {
		u := c.Param("user")
		if userIsValid(u) != nil {
			return err
		}

		// FIXME: verify actual payment value, now you can pay 0.1 forever

		var b []byte
		if err != nil {
			b = []byte(err.Error())
		} else {
			b = []byte(body)
		}
		err = appendUserLog(u, "ipn.txt", b)
		if err != nil {
			panic(err)
		}
		if n != nil {
			j, err := json.MarshalIndent(n, "", "\t")
			if err != nil {
				panic(err)
			}

			err = appendUserLog(u, "ipn.txt", []byte(j))
			if err != nil {
				panic(err)
			}

			//if n.TestIPN {
			// FIXME: allowing test, lets see how many people will scam
			//}

			dirs := []string{"tmp", "Mail", "public_html", "priv", "Maildir", path.Join("Maildir", "cur"), path.Join("Maildir", "new"), path.Join("Maildir", "tmp")}
			if n.TxnType == "subscr_signup" || n.TxnType == "cart" {
				for _, dir := range dirs {
					p := path.Join(ROOT, u, dir)
					_ = appendUserLog(u, "status.txt", []byte(fmt.Sprintf("chown %s %s", u, p)))
					_ = chown(u, p)
				}
			} else if n.TxnType == "subscr_cancel" {
				for _, dir := range dirs {
					p := path.Join(ROOT, u, dir)
					_ = appendUserLog(u, "status.txt", []byte(fmt.Sprintf("chown %s %s", u, p)))
					_ = chown("root", p)
				}
			}
		}
		return nil
	})

	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
