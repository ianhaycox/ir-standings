package iracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

type ResultsLink struct {
	Link string `json:"link,omitempty"`
}

// ResultLink https://members-ng.iracing.com/data/results/get?subsession_id=38280997
func (ir *IracingService) ResultLink(ctx context.Context, subsessionID string) (*ResultsLink, error) {
	queryParams := url.Values{}
	queryParams.Add("subsession_id", subsessionID)

	r, err := ir.client.PrepareRequest(ctx, Endpoint+"/data/results/get", http.MethodGet, queryParams, nil)
	if err != nil {
		return nil, err
	}

	response, err := ir.client.CallAPI(r) //nolint:bodyclose // ok
	if err != nil || response == nil {
		return nil, err
	}

	defer api.BodyClose(response)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		var apiError APIErrorResponse

		return nil, ir.client.ReportError(&apiError, response, body)
	}

	var link ResultsLink

	err = ir.client.Decode(&link, body, response.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("could not decode body:%s, err:%w", string(body), err)
	}

	return &link, nil
}
