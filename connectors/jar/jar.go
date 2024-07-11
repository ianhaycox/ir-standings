// Package cookiejar a barrel of cookies
package cookiejar

import (
	"net/http"
	"net/url"
	"sync"
)

type Cookies struct {
	*sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewCookieJar() http.CookieJar {
	return Cookies{
		cookies: make(map[string][]*http.Cookie),
	}
}

func (c Cookies) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.Lock()
	if _, ok := c.cookies[u.Host]; ok {
		c.cookies[u.Host] = append(c.cookies[u.Host], cookies...)
	} else {
		c.cookies[u.Host] = cookies
	}
	c.Unlock()
}

func (c Cookies) Cookies(u *url.URL) []*http.Cookie {
	return c.cookies[u.Host]
}
