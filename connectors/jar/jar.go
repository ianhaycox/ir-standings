// Package cookiejar a barrel of cookies
package cookiejar

import (
	"net/http"
	"net/url"
	"sync"
)

type Cookies struct {
	sync.Mutex
	cookies   map[string][]*http.Cookie
	persister Persister
}

func NewCookieJar(persister Persister) http.CookieJar {
	return &Cookies{
		cookies:   nil,
		persister: persister,
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

	if c.persister != nil {
		_ = c.persister.Write(c.cookies)
	}
}

func (c *Cookies) Cookies(u *url.URL) []*http.Cookie {
	c.Lock()
	defer c.Unlock()

	if c.cookies == nil {
		c.cookies = make(map[string][]*http.Cookie)

		if c.persister != nil {
			c.cookies, _ = c.persister.Read()
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
