//go:generate mockgen -package iracing -destination client_mock.go -source client.go

package iracing

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

	if ctx != nil {
		request, err = setHeadersFromContext(ctx, request)
		if err != nil {
			return request, err
		}
	}

	for header, value := range c.cfg.DefaultHeader {
		request.Header.Add(header, value)
	}

	return request, nil
}

func setHeadersFromContext(ctx context.Context, request *http.Request) (*http.Request, error) {
	if ctx == nil {
		return nil, fmt.Errorf("no context")
	}

	// add context to the request
	request = request.WithContext(ctx)

	// Walk through any authentication.
	// Usage: auth := context.WithValue(ctx, ContextAPIKey, APIKey{Key: "foo"})

	// Basic HTTP Authentication
	if auth, ok := ctx.Value(ContextBasicAuth).(BasicAuth); ok {
		request.SetBasicAuth(auth.UserName, auth.Password)
	}

	// AccessToken Authentication
	if auth, ok := ctx.Value(ContextAccessToken).(string); ok {
		request.Header.Add("Authorization", "Bearer "+auth)
	}

	// Basic API Key Authentication
	if auth, ok := ctx.Value(ContextBasicAPIKey).(BasicAPIKey); ok {
		request.Header.Add("Authorization", "Basic "+auth.Key)
	}

	// API Key Authentication
	if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
		var key string
		if auth.Prefix != "" {
			key = auth.Prefix + " " + auth.Key
		} else {
			key = auth.Key
		}

		request.Header.Add("X-API-KEY", key)
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
