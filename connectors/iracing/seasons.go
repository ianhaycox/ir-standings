package iracing

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/seasons"
)

func (ir *IracingAPI) Seasons(ctx context.Context) ([]seasons.Season, error) {
	var link results.ResultLink

	queryParams := url.Values{}

	err := ir.data.Get(ctx, &link, Endpoint+"/data/series/seasons", queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get seasons, err:%w", err)
	}

	var seasons []seasons.Season

	err = ir.data.CDN(ctx, link.Link, &seasons)
	if err != nil {
		return nil, fmt.Errorf("can not get seasons result:%s, err:%w", link.Link, err)
	}

	return seasons, nil
}
