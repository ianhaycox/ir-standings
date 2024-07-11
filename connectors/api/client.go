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

// APIClient manages communication over HTTP
type APIClient struct {
	cfg *Configuration
}

type APIClientInterface interface {
	CallAPI(request *http.Request) (*http.Response, error)
	PrepareRequest(ctx context.Context, path string, method string, queryParams url.Values, postBody interface{}) (request *http.Request, err error)
	Decode(v interface{}, b []byte, contentType string) (err error)
	ReportError(response *http.Response, body []byte) error
}

// NewAPIClient creates a new API client. Requires a userAgent string describing your application.
// optionally, a custom http.Client to allow for advanced features such as caching.
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	return &APIClient{
		cfg: cfg,
	}
}

// CallAPI do the request.
func (c *APIClient) CallAPI(request *http.Request) (*http.Response, error) {
	return c.cfg.HTTPClient.Do(request)
}

// PrepareRequest build the request
func (c *APIClient) PrepareRequest(ctx context.Context, path string, method string, queryParams url.Values, postBody interface{},
) (request *http.Request, err error) {
	var body *bytes.Buffer

	// Setup path and query parameters, path should have a leading '/', e.g. /bookings
	parsedURL, err := url.Parse(c.cfg.BasePath + path)
	if err != nil {
		return nil, err
	}

	// Adding Query Param
	query := parsedURL.Query()

	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	parsedURL.RawQuery = query.Encode()

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

	// Generate a new request
	if body != nil {
		request, err = http.NewRequest(method, parsedURL.String(), body)
	} else {
		request, err = http.NewRequest(method, parsedURL.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	// Add the user agent to the request.
	request.Header.Add("User-Agent", c.cfg.UserAgent)

	for header, value := range c.cfg.DefaultHeader {
		request.Header.Add(header, value)
	}

	return request, nil
}

func (c *APIClient) Decode(v interface{}, b []byte, contentType string) (err error) {
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

func (c *APIClient) ReportError(response *http.Response, body []byte) error {
	var v APIErrorResponse

	if len(body) > 0 {
		err := c.Decode(&v, body, response.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("server returned an error, body %s: err %w", string(body), err)
		}
	}

	return fmt.Errorf("server returned non-200 http code: %v, response '%s'", response.StatusCode, string(body))
}
