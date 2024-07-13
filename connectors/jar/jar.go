// Package cookiejar a barrel of cookies
package cookiejar

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

const mode = 0600

type Cookies struct {
	sync.Mutex
	cookies  map[string][]*http.Cookie
	filename string
}

func NewCookieJar(filename string) http.CookieJar {
	return &Cookies{
		cookies:  nil,
		filename: filename,
	}
}

func (c *Cookies) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.cookies[u.Host]; ok {
		c.cookies[u.Host] = c.merge(u.Host, cookies)
	} else {
		c.cookies[u.Host] = cookies
	}

	b, err := json.Marshal(c.cookies)
	if err != nil {
		log.Printf("broken cookie jar:%s", err)
	}

	err = os.WriteFile(c.filename, b, mode)
	if err != nil {
		log.Printf("can't close cookie jar:%s", err)
	}
}

func (c *Cookies) Cookies(u *url.URL) []*http.Cookie {
	c.Lock()
	defer c.Unlock()

	if c.cookies == nil {
		c.cookies = make(map[string][]*http.Cookie)

		b, err := os.ReadFile(c.filename)
		if err != nil {
			return c.cookies[u.Host]
		}

		err = json.Unmarshal(b, &c.cookies)
		if err != nil {
			log.Printf("can't open cookie jar:%s", err)
		}
	}

	return c.cookies[u.Host]
}

func (c *Cookies) merge(host string, cookies []*http.Cookie) []*http.Cookie {
	unique := make(map[string]*http.Cookie)

	for i := range c.cookies[host] {
		unique[c.cookies[host][i].Name] = c.cookies[host][i]
	}

	for i := range cookies {
		unique[cookies[i].Name] = cookies[i]
	}

	merged := make([]*http.Cookie, 0)

	for _, v := range unique {
		merged = append(merged, v)
	}

	return merged
}
