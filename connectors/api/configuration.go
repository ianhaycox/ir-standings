package api

import (
	"net/http"
)

type Configuration struct {
	BasePath      string            `json:"basePath,omitempty"` // no trailing '/'
	DefaultHeader map[string]string `json:"defaultHeader,omitempty"`
	UserAgent     string            `json:"userAgent,omitempty"`
	HTTPClient    *http.Client
}

func NewConfiguration(basePath string, httpClient *http.Client) *Configuration {
	cfg := &Configuration{
		BasePath:      basePath,
		HTTPClient:    httpClient,
		DefaultHeader: map[string]string{"Accept": "application/json", "Content-Type": "application/json"},
		UserAgent:     "github.com/ianhaycox/ir-standings/1.0.0/go",
	}

	return cfg
}

func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
