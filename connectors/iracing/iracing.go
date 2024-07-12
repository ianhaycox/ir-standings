// Package iracing API
package iracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

const (
	Endpoint string = "https://members-ng.iracing.com" // No trailing /
)

type IracingService struct {
	client api.APIClientInterface
	auth   api.Authenticator
}

type APIOptions struct {
}

type APIErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// "{\"authcode\":0,\"inactive\":false,\"message\":\"Invalid email address or password. Please try again.\",\"verificationRequired\":false}"

type AuthenticationGoodResponse struct {
	AuthCode   string `json:"authcode,omitempty"`
	CustomerID int    `json:"custId,omitempty"`
	Email      string `json:"email,omitempty"`
}

type AuthenticationBadResponse struct {
	AuthCode int    `json:"authcode,omitempty"`
	Message  string `json:"message,omitempty"`
}

type ResultsResponse struct {
}

func NewIracingService(client api.APIClientInterface, auth api.Authenticator) *IracingService {
	return &IracingService{
		client: client,
		auth:   auth,
	}
}

type IracingAPI interface {
	Authenticate(ctx context.Context)
	GetResults(ctx context.Context, opts *APIOptions) (ResultsResponse, error)
}

func (ir *IracingService) ReportError(response *http.Response, body []byte) error {
	var v APIErrorResponse

	if len(body) > 0 {
		err := ir.client.Decode(&v, body, response.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("server returned an error, body %s: err %w", string(body), err)
		}
	}

	return fmt.Errorf("server returned non-200 http code: %v, response '%+v'", response.StatusCode, v)
}

func (ir *IracingService) Authenticate(ctx context.Context) error {
	auth, err := ir.auth.BasicAuth()
	if err != nil {
		return fmt.Errorf("unable to apply authentication to context, err:%w", err)
	}

	r, err := ir.client.PrepareRequest(ctx, "/auth", http.MethodPost, url.Values{}, auth)
	if err != nil {
		return err
	}

	response, err := ir.client.CallAPI(r) //nolint:bodyclose // ok
	if err != nil || response == nil {
		return err
	}

	defer api.BodyClose(response)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return ir.ReportError(response, body)
	}

	var ok AuthenticationGoodResponse

	err = ir.client.Decode(&ok, body, response.Header.Get("Content-Type"))
	if err != nil {
		var bad AuthenticationBadResponse

		err = ir.client.Decode(&bad, body, response.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("failed to decode response, %w. Body:%s", err, string(body))
		}

		if bad.AuthCode == 0 {
			return fmt.Errorf("failed to authenticate: %s", bad.Message)
		}
	}

	return nil
}

func (ir *IracingService) GetResults(ctx context.Context, opts *APIOptions) (ResultsResponse, error) {
	return ResultsResponse{}, nil
}
