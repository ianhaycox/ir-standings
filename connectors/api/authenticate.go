//go:generate mockgen -package api -destination authenticate_mock.go -source authenticate.go
package api

import (
	"fmt"
)

type AuthenticationService struct {
	basicAuth *BasicAuth
	email     string
	password  string // ENCODEDPW=$(echo -n $PASSWORD$EMAILLOWER | openssl dgst -binary -sha256 | openssl base64)
}

type Authenticator interface {
	BasicAuth() (*BasicAuth, error)
}

func NewAuthenticationService(email string, password string) *AuthenticationService {
	return &AuthenticationService{
		email:    email,
		password: password,
	}
}

func (a *AuthenticationService) BasicAuth() (*BasicAuth, error) {
	if a.basicAuth != nil {
		return a.basicAuth, nil
	}

	basicAuth := BasicAuth{
		Email:    a.email,
		Password: a.password,
	}

	if basicAuth.Email == "" || basicAuth.Password == "" {
		return nil, fmt.Errorf("username:password combo can not be blank")
	}

	a.basicAuth = &basicAuth

	return a.basicAuth, nil
}
