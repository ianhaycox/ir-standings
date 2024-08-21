package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticationServiceHappyPath(t *testing.T) {
	t.Parallel()

	t.Run("Returns an Basic Auth username:password pair and no error", func(t *testing.T) {
		svc := NewAuthenticationService()
		require.NotNil(t, svc)

		basicAuth, err := svc.Credentials("user", "pass")
		assert.NoError(t, err)
		assert.Equal(t, &Credentials{Email: "user", EncodedPassword: "ViSRJ1toBu+Css9dtqRDMFOBGz4gUPQHJdNEYcnuXfc="}, basicAuth)
	})
}

func TestErrorPaths(t *testing.T) {
	t.Parallel()

	t.Run("returns error if username or password blank", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewAuthenticationService()
		_, err := svc.Credentials("", "blankuser")
		assert.ErrorContains(t, err, "email:password combo can not be blank")

		svc2 := NewAuthenticationService()
		_, err = svc2.Credentials("user", "")
		assert.ErrorContains(t, err, "email:password combo can not be blank")

		svc3 := NewAuthenticationService()
		_, err = svc3.Credentials("", "")
		assert.ErrorContains(t, err, "email:password combo can not be blank")
	})
}

func TestEncode(t *testing.T) {
	// echo -n "barfoo" | openssl dgst -binary -sha256 | openssl base64
	// iOzekl2jxvjsPRQGg9qdKkIvJsGuHZIS2h5aU0FtzIg=

	a := NewAuthenticationService()
	assert.Equal(t, "iOzekl2jxvjsPRQGg9qdKkIvJsGuHZIS2h5aU0FtzIg=", a.encode("foo", "bar"))
}
