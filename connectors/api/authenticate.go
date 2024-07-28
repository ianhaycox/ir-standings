//go:generate mockgen -package api -destination authenticate_mock.go -source authenticate.go
package api

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// Credentials provides post body of encoded credentials
type Credentials struct {
	Email           string `json:"email,omitempty"`
	EncodedPassword string `json:"password,omitempty"`
}

type AuthenticationService struct {
	credentials *Credentials
	email       string
	password    string // ENCODEDPW=$(echo -n $PASSWORD$EMAILLOWER | openssl dgst -binary -sha256 | openssl base64)
}

type Authenticator interface {
	Credentials() (*Credentials, error)
}

func NewAuthenticationService(email string, password string) *AuthenticationService {
	return &AuthenticationService{
		email:    email,
		password: password,
	}
}

func (a *AuthenticationService) Credentials() (*Credentials, error) {
	if a.credentials != nil {
		return a.credentials, nil
	}

	if a.email == "" || a.password == "" {
		return nil, fmt.Errorf("username:password combo can not be blank")
	}

	credentials := Credentials{
		Email:           a.email,
		EncodedPassword: a.encode(),
	}

	a.credentials = &credentials

	return a.credentials, nil
}

func (a *AuthenticationService) encode() string {
	h := sha256.New()
	h.Write([]byte(a.password + strings.ToLower(a.email)))
	hash := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}
