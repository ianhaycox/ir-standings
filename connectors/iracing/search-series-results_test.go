package iracing

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/ianhaycox/ir-standings/model/iracing/results/searchseries"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSearchSeriesResults(t *testing.T) {
	t.Run("Returns no error if successfully gets series results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var ssr searchseries.SearchSeriesResults
		ssrResponse := searchseries.SearchSeriesResults{Data: searchseries.Data{Success: true, ChunkInfo: searchseries.ChunkInfo{NumChunks: 2}}}

		queryParams := url.Values{}
		queryParams.Add("season_year", fmt.Sprintf("%d", 2024))
		queryParams.Add("season_quarter", fmt.Sprintf("%d", 2))
		queryParams.Add("series_id", fmt.Sprintf("%d", 285))
		queryParams.Add("official_only", "true")
		queryParams.Add("event_types", fmt.Sprintf("%d", 5))
		mockDataAPI := NewMockIracingDataAPI(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &ssr, Endpoint+"/data/results/search_series", queryParams).Return(nil).SetArg(1, ssrResponse)

		ir := NewIracingService(mockDataAPI, nil)

		actual, err := ir.SearchSeriesResults(ctx, 2024, 2, 285)
		assert.NoError(t, err)
		assert.True(t, actual.Data.Success)
		assert.Equal(t, 2, actual.Data.ChunkInfo.NumChunks)
	})
}

func TestSearchSeriesResultsErrors(t *testing.T) {
	t.Run("Returns error if Get series results fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockDataAPI := NewMockIracingDataAPI(ctrl)
		mockDataAPI.EXPECT().Get(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("get failed"))

		ir := NewIracingService(mockDataAPI, nil)

		actual, err := ir.SearchSeriesResults(ctx, 2024, 2, 285)
		assert.ErrorContains(t, err, "get failed")
		assert.Nil(t, actual)
	})
}
