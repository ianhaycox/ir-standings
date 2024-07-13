package api

import (
	"net/http"
)

const UserAgent = "github.com/ianhaycox/ir-standings/1.0.0/go"

type Configuration struct {
	DefaultHeader map[string]string `json:"defaultHeader,omitempty"`
	UserAgent     string            `json:"userAgent,omitempty"`
	HTTPClient    *http.Client
}

func NewConfiguration(httpClient *http.Client, userAgent string) *Configuration {
	cfg := &Configuration{
		HTTPClient:    httpClient,
		DefaultHeader: make(map[string]string),
		UserAgent:     userAgent,
	}

	return cfg
}

func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
