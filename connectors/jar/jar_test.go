package cookiejar

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCookieJar(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "example")
	assert.NoError(t, err)

	defer os.RemoveAll(tempDir)

	cookiesFile := filepath.Join(tempDir, "cookies")

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
		store := NewStore(cookiesFile)
		jar := NewCookieJar(store)

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
		require.Len(t, cookies, 2)

		actual := make(map[string]string, 0)
		actual[cookies[0].Name] = cookies[0].Value
		actual[cookies[1].Name] = cookies[1].Value
		assert.Equal(t, map[string]string{"Flavor": "Oatmeal Raisin", "Third": "Request"}, actual)
	})

	t.Run("Gets and sets a cookie with no persistence", func(t *testing.T) {
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
		jar := NewCookieJar(nil)

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
		require.Len(t, cookies, 2)

		actual := make(map[string]string, 0)
		actual[cookies[0].Name] = cookies[0].Value
		actual[cookies[1].Name] = cookies[1].Value
		assert.Equal(t, map[string]string{"Flavor": "Oatmeal Raisin", "Third": "Request"}, actual)
	})
}

func TestCookieJarErrors(t *testing.T) {
	t.Run("Gets and sets a cookie", func(t *testing.T) {
	})
}
