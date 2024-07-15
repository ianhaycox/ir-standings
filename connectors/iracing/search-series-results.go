package iracing

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ianhaycox/ir-standings/model/iracing/results/searchseries"
)

const EventRace = 5
const KamelSeriesID = 285

func (ir *IracingService) SearchSeriesResults(ctx context.Context, seasonYear, seasonQuarter, seriesID int) (*searchseries.SearchSeriesResults, error) {
	queryParams := url.Values{}
	queryParams.Add("season_year", fmt.Sprintf("%d", seasonYear))
	queryParams.Add("season_quarter", fmt.Sprintf("%d", seasonQuarter))
	queryParams.Add("series_id", fmt.Sprintf("%d", seriesID))
	queryParams.Add("official_only", "true")
	queryParams.Add("event_types", fmt.Sprintf("%d", EventRace))

	var results searchseries.SearchSeriesResults

	err := ir.data.Get(ctx, &results, Endpoint+"/data/results/search_series", queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get series ID:%d results for year:%d, quarter:%d, err:%w", seriesID, seasonYear, seasonQuarter, err)
	}

	return &results, nil
}
