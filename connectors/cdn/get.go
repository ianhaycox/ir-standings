package cdn

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

func (cdn *CDNService) CDN(ctx context.Context, link string, v interface{}) error {
	r, err := cdn.client.PrepareRequest(ctx, link, http.MethodGet, nil, nil)
	if err != nil {
		return err
	}

	response, err := cdn.client.CallAPI(r) //nolint:bodyclose // ok
	if err != nil || response == nil {
		return err
	}

	defer api.BodyClose(response)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		var errorResponse AWSErrorResponse

		return cdn.client.ReportError(&errorResponse, response, body)
	}

	err = cdn.client.Decode(v, body, response.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("failed to decode result response, %w. Body:%s", err, string(body))
	}

	return nil
}
