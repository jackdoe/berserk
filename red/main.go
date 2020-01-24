package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	ipn "github.com/jackdoe/gin-ipn"
)

const ROOT = "/mnt/home_attached"

func main() {

	r := gin.Default()

	r.POST("/register/:user", func(c *gin.Context) {
		key, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = CreateSystemUser(c.Param("user"), key)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		u, err := NewUser(c.Param("user"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		requestDump, err := httputil.DumpRequest(c.Request, true)
		if err != nil {
			panic(err)
		}
		u.LogP("register.txt", []byte(requestDump))

		c.String(200, fmt.Sprintf(AFTER_REGISTER, u.Name, u.Name))
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
		rp := c.Param("path")
		u, err := NewUser(c.Param("user"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// cleanup golang http.ServeFile special handling of index.html
		if strings.HasSuffix(rp, "/index.html") {
			c.Redirect(302, "/~"+u.Name+"/"+strings.TrimSuffix(rp, "/index.html"))
			return
		}

		local := path.Join(u.Home, "public_html")
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
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		u, err := NewUser(strings.Trim(c.Param("user"), "~"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if u.Uid < 1000 {
			c.JSON(400, gin.H{"error": "invalid user"})
			return
		}

		err = u.Mail(body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.String(200, "OK")
	})

	r.GET("/sub/:user", func(c *gin.Context) {
		u, err := NewUser(c.Param("user"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		prefix := "https://www.paypal.com/cgi-bin/webscr"

		url := prefix + "?cmd=_xclick-subscriptions&business=jack%40baxx.dev&a3=1&p3=1&t3=M&item_name=berserk.red+-+personal+website&return=https%3A%2F%2Fberserk.red%2Fthanks_for_paying&a1=0.1&no_shipping=1&p1=1&t1=M&src=1&sra=1&no_note=1&no_note=1&currency_code=EUR&lc=GB&notify_url=https%3A%2F%2Fberserk.red%2Fipn%2F" + u.Name
		c.Redirect(http.StatusFound, url)
	})

	ipn.Listener(r, "/ipn/:user", func(c *gin.Context, err error, body string, n *ipn.Notification) error {
		u, errx := NewUser(c.Param("user"))
		if errx != nil {
			return errx
		}

		// FIXME: verify actual payment value, now you can pay 0.1 forever

		var b []byte
		if err != nil {
			b = []byte(err.Error())
		} else {
			b = []byte(body)
		}
		u.LogP("ipn.txt", b)
		if n != nil {
			j, err := json.MarshalIndent(n, "", "\t")
			if err != nil {
				panic(err)
			}

			u.LogP("ipn.txt", []byte(j))

			//if n.TestIPN {
			// FIXME: allowing test, lets see how many people will scam
			//}

			if n.TxnType == "subscr_signup" || n.PayerID == "TESTBUYERID01_ENABLE" {
				u.LogP("status.txt", []byte(fmt.Sprintf("ENABLE %v", u)))
				_ = u.Enable()
			} else if n.TxnType == "subscr_cancel" || n.PayerID == "TESTBUYERID01_CANCEL" {
				u.LogP("status.txt", []byte(fmt.Sprintf("DISABLE %v", u)))
				_ = u.Disable()
			}
		}
		return nil
	})

	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
