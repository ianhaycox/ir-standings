package iracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

func (ids *IracingDataAPI) Get(ctx context.Context, v interface{}, path string, queryParams url.Values) error {
	r, err := ids.client.PrepareRequest(ctx, path, http.MethodGet, queryParams, nil)
	if err != nil {
		return err
	}

	response, err := ids.client.CallAPI(r) //nolint:bodyclose // ok
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

		return ids.client.ReportError(&apiError, response, body)
	}

	err = ids.client.Decode(v, body, response.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("could not decode body:%s, err:%w", string(body), err)
	}

	return nil
}
