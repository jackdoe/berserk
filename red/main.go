package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	ipn "github.com/jackdoe/gin-ipn"
	gemini "github.com/jackdoe/net-gemini"
	"github.com/mitchellh/go-finger"
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
		u, err := CreateSystemUser(c.Param("user"), key)
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
		c.Header("Content-Type", "text/html; charset=utf-8")

		files, _ := ioutil.ReadDir(ROOT)

		type homedir struct {
			dir   string
			t     time.Time
			count int
		}
		available := []*homedir{}
		for _, dir := range files {
			ph := path.Join(ROOT, dir.Name(), "public_html")
			if dirExists(ph) {
				ds, err := os.Stat(ph)
				if err != nil {
					continue
				}

				if ds.Mode().Perm()&4 == 0 {
					// no permissions
					continue
				}

				hd := &homedir{dir: dir.Name()}

				_ = filepath.Walk(ph, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}

					t := info.ModTime()
					if hd.t.Before(t) {
						hd.t = t
					}
					hd.count++
					return nil
				})

				if hd.count > 0 {
					available = append(available, hd)
				}
			}
		}

		sort.Slice(available, func(i, j int) bool {
			return available[j].t.Before(available[i].t)
		})

		var out strings.Builder
		out.WriteString("<html><head><title>berserk.red</title></head><body><pre>")
		out.WriteString("users with websites:\n\n")
		for _, a := range available {
			href := fmt.Sprintf("https://berserk.red/~%s/", a.dir)
			out.WriteString(fmt.Sprintf("<a href='%s'>%s</a>\nupdated %s\n\n", href, href, humanize.Time(a.t)))
		}

		out.WriteString(SLASH)

		out.WriteString("</pre></body></html>")
		c.String(200, out.String())
	})

	r.GET("/tos", func(c *gin.Context) {
		c.String(200, LICENSE)
	})

	r.GET("/thanks_for_paying", func(c *gin.Context) {
		c.String(200, THANKS_FOR_PAYING)
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

	go func() {
		gemini.HandleFunc("/~", func(w *gemini.Response, r *gemini.Request) {
			p := strings.TrimPrefix(r.URL.Path, "/~")
			if len(p) == 0 {
				w.SetStatus(gemini.StatusNotFound, "Not Found")
				return
			}

			splitted := strings.SplitN(p, "/", 2)

			u, err := NewUser(splitted[0])
			if err != nil {
				w.SetStatus(gemini.StatusTemporaryFailure, err.Error())
				return
			}

			local := path.Join(u.Home, "public_html")
			p = local
			if len(splitted) > 1 {
				p = path.Join(p, filepath.Clean(splitted[1]))
			}

			l, err := os.Readlink(p)
			if err == nil {
				p = l
			}

			// dont allow symlinks leading outside of home/public_html
			if !strings.HasPrefix(p, local) {
				w.SetStatus(gemini.StatusTemporaryFailure, "out of home")
				return
			}
			gemini.ServeFilePath(p, w, r)
		})

		gemini.HandleFunc("/", func(w *gemini.Response, r *gemini.Request) {
			w.Write([]byte(SLASH))
		})

		log.Fatal(gemini.ListenAndServeTLS(":"+os.Getenv("GEMINI_PORT"), os.Getenv("GEMINI_CRT"), os.Getenv("GEMINI_KEY")))
	}()

	go func() {
		log.Fatal(finger.Serve(finger.HandlerFunc(func(ctx context.Context, w io.Writer, q *finger.Query) {
			home := path.Join(ROOT, q.Username, "public_html")
			ds, err := os.Stat(home)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("%q not found\n", q.Username)))
				return
			}

			if ds.Mode().Perm()&4 == 0 {
				w.Write([]byte(fmt.Sprintf("%q not found\n", q.Username)))
				return
			}

			w.Write([]byte(fmt.Sprintf("printing https://berserk.red/~%s/finger.txt\n%s\n\n", q.Username, strings.Repeat("-", 8))))

			b, err := ioutil.ReadFile(path.Join(home, "finger.txt"))
			if err != nil {
				w.Write([]byte(fmt.Sprintf("nop, finger.txt is missing\ncheck out https://berserk.red/~%s though\n", q.Username)))
			} else {
				w.Write(b)
			}
		})))
	}()

	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
