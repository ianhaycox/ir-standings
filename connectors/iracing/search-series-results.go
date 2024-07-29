package iracing

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
)

const EventRace = 5
const KamelSeriesID model.SeriesID = 285

func (ir *IracingAPI) SearchSeriesResults(ctx context.Context, seasonYear, seasonQuarter, seriesID int) ([]searchseries.SearchSeriesResult, error) {
	queryParams := url.Values{}
	queryParams.Add("season_year", fmt.Sprintf("%d", seasonYear))
	queryParams.Add("season_quarter", fmt.Sprintf("%d", seasonQuarter))
	queryParams.Add("series_id", fmt.Sprintf("%d", seriesID))
	queryParams.Add("official_only", "true")
	queryParams.Add("event_types", fmt.Sprintf("%d", EventRace))

	var seriesResults searchseries.SearchSeriesResults

	err := ir.data.Get(ctx, &seriesResults, Endpoint+"/data/results/search_series", queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get series ID:%d results for year:%d, quarter:%d, err:%w", seriesID, seasonYear, seasonQuarter, err)
	}

	var ssResults []searchseries.SearchSeriesResult

	if seriesResults.Data.Success {
		for i := range seriesResults.Data.ChunkInfo.ChunkFileNames {
			var ssr []searchseries.SearchSeriesResult

			url := seriesResults.Data.ChunkInfo.BaseDownloadURL + seriesResults.Data.ChunkInfo.ChunkFileNames[i]

			err := ir.data.CDN(ctx, url, &ssr)
			if err != nil {
				return nil, fmt.Errorf("can not get search series result:%s, err:%w", url, err)
			}

			ssResults = append(ssResults, ssr...)
		}
	}

	return ssResults, nil
}
