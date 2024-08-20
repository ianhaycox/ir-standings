package main

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	"github.com/stretchr/testify/assert"
)

func TestCookie(t *testing.T) {
	t.Skip()

	cookieStore := cookiejar.NewStore(iracing.CookiesFile)
	jar := cookiejar.NewCookieJar(cookieStore)

	cookies := jar.Cookies(&url.URL{Host: "members-ng.iracing.com"})

	fmt.Println(cookies)

	assert.NotNil(t, cookies)
}
