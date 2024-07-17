package main

import (
	"context"
	"os"
	"testing"

	"github.com/ianhaycox/ir-standings/connectors/iracing"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStandings(t *testing.T) {
	t.Run("For a list of series results filter and get the full results for broadcast races", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		ssResults := []searchseries.SearchSeriesResult{}
		results := []results.Result{}

		ir := iracing.NewMockIracingService(ctrl)
		ir.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return(ssResults, nil)
		ir.EXPECT().SeasonBroadcastResults(ctx, ssResults).Return(results, nil)

		_, err := standings(ctx, ir, 2024, 2)
		assert.NoError(t, err)
	})

	t.Run("args should return season info", func(t *testing.T) {
		os.Args = []string{"", "1", "2"}
		y, q, err := args()
		assert.NoError(t, err)
		assert.Equal(t, 1, y)
		assert.Equal(t, 2, q)
	})
}
