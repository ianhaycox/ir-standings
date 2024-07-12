package cookiejar

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieJar(t *testing.T) {
	t.Run("Gets and sets a cookie", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cookie, err := r.Cookie("Flavor"); err != nil {
				http.SetCookie(w, &http.Cookie{Name: "Flavor", Value: "Chocolate Chip"})
			} else {
				if cookie.Value == "Oatmeal Raisin" {
					http.SetCookie(w, &http.Cookie{Name: "Third", Value: "Request"})
				}

				cookie.Value = "Oatmeal Raisin"
				http.SetCookie(w, cookie)
			}
		}))
		defer ts.Close()

		ctx := context.Background()
		jar := NewCookieJar()

		httpClient := http.DefaultClient
		httpClient.Jar = jar
		u, err := url.Parse(ts.URL)
		assert.NoError(t, err)

		request, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
		assert.NoError(t, err)

		response, err := httpClient.Do(request)
		assert.NoError(t, err)

		response.Body.Close()

		cookies := jar.Cookies(u)
		assert.Len(t, cookies, 1)
		assert.Equal(t, "Flavor", cookies[0].Name)
		assert.Equal(t, "Chocolate Chip", cookies[0].Value)

		// Next request
		request, err = http.NewRequestWithContext(ctx, "GET", u.String(), nil)
		assert.NoError(t, err)

		response, err = httpClient.Do(request)
		assert.NoError(t, err)

		response.Body.Close()

		cookies = jar.Cookies(u)
		assert.Len(t, cookies, 1)
		assert.Equal(t, "Flavor", cookies[0].Name)
		assert.Equal(t, "Oatmeal Raisin", cookies[0].Value)

		// Third request
		request, err = http.NewRequestWithContext(ctx, "GET", u.String(), nil)
		assert.NoError(t, err)

		response, err = httpClient.Do(request)
		assert.NoError(t, err)

		response.Body.Close()

		cookies = jar.Cookies(u)
		assert.Len(t, cookies, 2)
		assert.Equal(t, "Flavor", cookies[0].Name)
		assert.Equal(t, "Oatmeal Raisin", cookies[0].Value)
		assert.Equal(t, "Third", cookies[1].Name)
		assert.Equal(t, "Request", cookies[1].Value)
	})
}
