package iracing

import (
	"context"
	"fmt"

	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
)

func (ir *IracingAPI) SeasonBroadcastResults(ctx context.Context, ssResults []searchseries.SearchSeriesResult) ([]results.Result, error) {
	seasonResults := make([]results.Result, 0)

	for j := range ssResults {
		if !ssResults[j].IsBroadcast() {
			continue
		}

		link, err := ir.ResultLink(ctx, ssResults[j].SubsessionID)
		if err != nil {
			return nil, fmt.Errorf("can not get result link for sub session ID:%d, err:%w", ssResults[j].SubsessionID, err)
		}

		var res results.Result

		err = ir.data.CDN(ctx, link.Link, &res)
		if err != nil {
			return nil, fmt.Errorf("can not get result:%s, err:%w", link.Link, err)
		}

		seasonResults = append(seasonResults, res)
	}

	return seasonResults, nil
}
