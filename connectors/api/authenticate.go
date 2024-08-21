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
}

type Authenticator interface {
	Credentials(email string, password string) (*Credentials, error)
}

func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{}
}

func (a *AuthenticationService) Credentials(email string, password string) (*Credentials, error) {
	if a.credentials != nil {
		return a.credentials, nil
	}

	if email == "" || password == "" {
		return nil, fmt.Errorf("email:password combo can not be blank")
	}

	credentials := Credentials{
		Email:           email,
		EncodedPassword: a.encode(email, password),
	}

	a.credentials = &credentials

	return a.credentials, nil
}

func (a *AuthenticationService) encode(email, password string) string {
	h := sha256.New()
	h.Write([]byte(password + strings.ToLower(email)))
	hash := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}
