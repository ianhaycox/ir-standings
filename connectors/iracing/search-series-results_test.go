package iracing

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSearchSeriesResults(t *testing.T) {
	t.Run("Returns no error if successfully gets series results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var ssr searchseries.SearchSeriesResults
		ssrResponse := searchseries.SearchSeriesResults{
			Data: searchseries.Data{
				Success: true,
				ChunkInfo: searchseries.ChunkInfo{
					NumChunks:      2,
					ChunkFileNames: []string{"res1.json", "res2.json"},
				},
			},
		}

		queryParams := url.Values{}
		queryParams.Add("season_year", fmt.Sprintf("%d", 2024))
		queryParams.Add("season_quarter", fmt.Sprintf("%d", 2))
		queryParams.Add("series_id", fmt.Sprintf("%d", 285))
		queryParams.Add("official_only", "true")
		queryParams.Add("event_types", fmt.Sprintf("%d", 5))

		var ssrs []searchseries.SearchSeriesResult
		mockDataAPI := NewMockIracingDataService(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &ssr, Endpoint+"/data/results/search_series", queryParams).Return(nil).SetArg(1, ssrResponse)
		mockDataAPI.EXPECT().CDN(ctx, "res1.json", &ssrs).Return(nil).SetArg(2, []searchseries.SearchSeriesResult{{SubsessionID: 1}, {SubsessionID: 2}})
		mockDataAPI.EXPECT().CDN(ctx, "res2.json", &ssrs).Return(nil).SetArg(2, []searchseries.SearchSeriesResult{{SubsessionID: 3}})

		ir := NewIracingService(nil, mockDataAPI, nil)

		actual, err := ir.SearchSeriesResults(ctx, 2024, 2, 285)
		assert.NoError(t, err)
		assert.Len(t, actual, 3)
		assert.Equal(t, 1, actual[0].SubsessionID)
		assert.Equal(t, 2, actual[1].SubsessionID)
		assert.Equal(t, 3, actual[2].SubsessionID)
	})
}

func TestSearchSeriesResultsErrors(t *testing.T) {
	t.Run("Returns error if Get series results fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockDataAPI := NewMockIracingDataService(ctrl)
		mockDataAPI.EXPECT().Get(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("get failed"))

		ir := NewIracingService(nil, mockDataAPI, nil)

		actual, err := ir.SearchSeriesResults(ctx, 2024, 2, 285)
		assert.ErrorContains(t, err, "get failed")
		assert.Nil(t, actual)
	})

	t.Run("Returns error if CDN series results fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var ssr searchseries.SearchSeriesResults
		ssrResponse := searchseries.SearchSeriesResults{
			Data: searchseries.Data{
				Success: true,
				ChunkInfo: searchseries.ChunkInfo{
					NumChunks:      1,
					ChunkFileNames: []string{"res.json"},
				},
			},
		}

		mockDataAPI := NewMockIracingDataService(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &ssr, gomock.Any(), gomock.Any()).Return(nil).SetArg(1, ssrResponse)
		mockDataAPI.EXPECT().CDN(ctx, "res.json", gomock.Any()).Return(nil).Return(fmt.Errorf("cdn failed"))

		ir := NewIracingService(nil, mockDataAPI, nil)

		actual, err := ir.SearchSeriesResults(ctx, 2024, 2, 285)
		assert.ErrorContains(t, err, "cdn failed")
		assert.Nil(t, actual)
	})
}
