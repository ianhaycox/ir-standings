package main

import (
	"context"
	"testing"

	"github.com/ianhaycox/ir-standings/connectors/iracing"
	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
	"github.com/ianhaycox/ir-standings/model/data/seasons"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApp(t *testing.T) {
	t.Run("Login OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		irAPI := iracing.NewMockIracingService(ctrl)
		irAPI.EXPECT().Authenticate(ctx, "test@example.com", "pass").Return(nil)
		irAPI.EXPECT().Cars(ctx).Return([]cars.Car{}, nil)
		irAPI.EXPECT().CarClasses(ctx).Return([]cars.CarClass{}, nil)
		irAPI.EXPECT().Seasons(ctx).Return([]seasons.Season{{SeriesID: 99, SeasonYear: 2023, SeasonQuarter: 2}}, nil)
		irAPI.EXPECT().SearchSeriesResults(ctx, 2023, 2, 99).Return([]searchseries.SearchSeriesResult{}, nil)
		irAPI.EXPECT().SeasonBroadcastResults(ctx, []searchseries.SearchSeriesResult{}).Return([]results.Result{}, nil)

		a := NewApp(nil, irAPI, nil, 1, 1, 99, 1)
		a.startup(ctx)

		response := a.Login("test@example.com", "pass")
		assert.True(t, response)
	})

	t.Run("Fake Login OK", func(t *testing.T) {
		ctx := context.TODO()
		a := NewApp(nil, nil, nil, 1, 1, 99, 1)
		a.startup(ctx)

		response := a.Login("test", "pass")
		assert.True(t, response)
	})
}
