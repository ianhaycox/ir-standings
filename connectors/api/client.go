//go:generate mockgen -package api -destination client_mock.go -source client.go

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient manages communication over HTTP
type HTTPClient struct {
	cfg     *Configuration
	lastURL string
}

type API interface {
	CallAPI(request *http.Request) (*http.Response, error)
	PrepareRequest(ctx context.Context, path string, method string, queryParams url.Values, postBody interface{}) (request *http.Request, err error)
	Decode(v interface{}, b []byte, contentType string) (err error)
	ReportError(v interface{}, response *http.Response, body []byte) error
}

// NewHTTPClient creates a new API client. Requires a userAgent string describing your application.
// optionally, a custom http.Client to allow for advanced features such as caching.
func NewHTTPClient(cfg *Configuration) *HTTPClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	return &HTTPClient{
		cfg: cfg,
	}
}

// CallAPI do the request.
func (c *HTTPClient) CallAPI(request *http.Request) (*http.Response, error) {
	return c.cfg.HTTPClient.Do(request)
}

// PrepareRequest build the request
func (c *HTTPClient) PrepareRequest(ctx context.Context, path string, method string, queryParams url.Values, postBody interface{},
) (request *http.Request, err error) {
	var body *bytes.Buffer

	// Setup path and query parameters, path is the full URL
	parsedURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Adding Query Param
	if len(queryParams) > 0 {
		query := parsedURL.Query()

		for k, v := range queryParams {
			for _, iv := range v {
				query.Add(k, iv)
			}
		}

		// Encode the parameters.
		parsedURL.RawQuery = query.Encode()
	}

	// Encode body
	if postBody != nil {
		body = &bytes.Buffer{}

		if reader, ok := postBody.(io.Reader); ok {
			_, err = body.ReadFrom(reader)
		} else {
			err = json.NewEncoder(body).Encode(postBody)
		}

		if err != nil {
			return nil, err
		}
	}

	c.lastURL = parsedURL.String() // For error reporting

	// Generate a new request
	if body != nil {
		request, err = http.NewRequestWithContext(ctx, method, c.lastURL, body)
	} else {
		request, err = http.NewRequestWithContext(ctx, method, c.lastURL, nil)
	}

	if err != nil {
		return nil, err
	}

	// Add the user agent to the request.
	if c.cfg.UserAgent != "" {
		request.Header.Add("User-Agent", c.cfg.UserAgent)
	}

	for header, value := range c.cfg.DefaultHeader {
		request.Header.Add(header, value)
	}

	return request, nil
}

func (c *HTTPClient) Decode(v interface{}, b []byte, contentType string) (err error) {
	if strings.Contains(contentType, "application/xml") {
		if err = xml.Unmarshal(b, v); err != nil {
			return err
		}

		return nil
	} else if strings.Contains(contentType, "application/json") {
		if err = json.Unmarshal(b, v); err != nil {
			return err
		}

		return nil
	}

	return errors.New("undefined Content-Type in response")
}

func (c *HTTPClient) ReportError(v interface{}, response *http.Response, body []byte) error {
	if len(body) > 0 {
		err := c.Decode(&v, body, response.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("server %s returned an error, body %s: err %w", c.lastURL, string(body), err)
		}
	}

	return fmt.Errorf("server %s returned non-200 http code: %v, response '%s'", c.lastURL, response.StatusCode, string(body))
}
