package iracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

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

func (ir *IracingAPI) Authenticate(ctx context.Context, username, password string) error {
	auth, err := ir.auth.Credentials(username, password)
	if err != nil {
		return fmt.Errorf("unable to retrieve credentials, err:%w", err)
	}

	r, err := ir.client.PrepareRequest(ctx, Endpoint+"/auth", http.MethodPost, url.Values{}, auth)
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
		var apiError APIErrorResponse

		return ir.client.ReportError(&apiError, response, body)
	}

	// Two different responses are received after an authentication request
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
