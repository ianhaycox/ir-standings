package iracing

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSeasonBroadcastResults(t *testing.T) {
	t.Run("For a list of series results filter and get the full results for broadcast races", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var (
			link1 results.ResultLink
			link3 results.ResultLink
			link5 results.ResultLink
			res1  results.Result
			res3  results.Result
			res5  results.Result
		)

		linkResponse1 := results.ResultLink{Link: "https://cdn.com/result/1"}
		linkResponse3 := results.ResultLink{Link: "https://cdn.com/result/3"}
		linkResponse5 := results.ResultLink{Link: "https://cdn.com/result/5"}

		mockDataAPI := NewMockIracingDataService(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &link1, Endpoint+"/data/results/get", url.Values{"subsession_id": {"1"}}).Return(nil).SetArg(1, linkResponse1)
		mockDataAPI.EXPECT().CDN(ctx, "https://cdn.com/result/1", &res1).Return(nil).SetArg(2, results.Result{SubsessionID: 1})
		mockDataAPI.EXPECT().Get(ctx, &link3, Endpoint+"/data/results/get", url.Values{"subsession_id": {"3"}}).Return(nil).SetArg(1, linkResponse3)
		mockDataAPI.EXPECT().CDN(ctx, "https://cdn.com/result/3", &res3).Return(nil).SetArg(2, results.Result{SubsessionID: 3})
		mockDataAPI.EXPECT().Get(ctx, &link5, Endpoint+"/data/results/get", url.Values{"subsession_id": {"5"}}).Return(nil).SetArg(1, linkResponse5)
		mockDataAPI.EXPECT().CDN(ctx, "https://cdn.com/result/5", &res5).Return(nil).SetArg(2, results.Result{SubsessionID: 5})

		found := []searchseries.SearchSeriesResult{
			{
				SubsessionID: 1,
				StartTime:    time.Date(2024, 3, 16, 17, 0, 0, 0, time.UTC), // Saturday 17:00
			},
			{
				SubsessionID: 2,
				StartTime:    time.Date(2024, 3, 17, 17, 0, 0, 0, time.UTC), // Sunday
			},
			{
				SubsessionID: 3,
				StartTime:    time.Date(2024, 3, 23, 17, 0, 0, 0, time.UTC), // Saturday 17:00
			},
			{
				SubsessionID: 4,
				StartTime:    time.Date(2024, 3, 23, 19, 0, 0, 0, time.UTC), // 19:00
			},
			{
				SubsessionID: 5,
				StartTime:    time.Date(2024, 3, 30, 17, 0, 0, 0, time.UTC), // Saturday 17:00
			},
		}

		ir := NewIracingService(nil, mockDataAPI, nil)

		actual, err := ir.SeasonBroadcastResults(ctx, found)
		assert.NoError(t, err)
		assert.Len(t, actual, 3)
		assert.Equal(t, results.Result{SubsessionID: 1}, actual[0])
		assert.Equal(t, results.Result{SubsessionID: 3}, actual[1])
		assert.Equal(t, results.Result{SubsessionID: 5}, actual[2])
	})
}

func TestSeasonBroadcastResultsErrors(t *testing.T) {
	t.Run("Should return an error if we fail to get the results link", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var link1 results.ResultLink

		mockDataAPI := NewMockIracingDataService(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &link1, Endpoint+"/data/results/get", url.Values{"subsession_id": {"1"}}).Return(fmt.Errorf("failed link"))

		found := []searchseries.SearchSeriesResult{
			{
				SubsessionID: 1,
				StartTime:    time.Date(2024, 3, 16, 17, 0, 0, 0, time.UTC), // Saturday 17:00
			},
		}

		ir := NewIracingService(nil, mockDataAPI, nil)

		_, err := ir.SeasonBroadcastResults(ctx, found)
		assert.ErrorContains(t, err, "failed link")
	})

	t.Run("Should return an error if we fail to get the results data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var (
			link1 results.ResultLink
			res1  results.Result
		)

		linkResponse1 := results.ResultLink{Link: "https://cdn.com/result/1"}

		mockDataAPI := NewMockIracingDataService(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &link1, Endpoint+"/data/results/get", url.Values{"subsession_id": {"1"}}).Return(nil).SetArg(1, linkResponse1)
		mockDataAPI.EXPECT().CDN(ctx, "https://cdn.com/result/1", &res1).Return(fmt.Errorf("failed data"))
		found := []searchseries.SearchSeriesResult{
			{
				SubsessionID: 1,
				StartTime:    time.Date(2024, 3, 16, 17, 0, 0, 0, time.UTC), // Saturday 17:00
			},
		}

		ir := NewIracingService(nil, mockDataAPI, nil)

		_, err := ir.SeasonBroadcastResults(ctx, found)
		assert.ErrorContains(t, err, "failed data")
	})
}
