package iracing

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ianhaycox/ir-standings/model/iracing/results"
)

// ResultLink https://members-ng.iracing.com/data/results/get?subsession_id=38280997
func (ir *IracingService) ResultLink(ctx context.Context, subsessionID int) (*results.ResultLink, error) {
	queryParams := url.Values{}
	queryParams.Add("subsession_id", fmt.Sprintf("%d", subsessionID))

	var link results.ResultLink

	err := ir.data.Get(ctx, &link, Endpoint+"/data/results/get", queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get result link for subsession ID:%d, err:%w", subsessionID, err)
	}

	return &link, nil
}
