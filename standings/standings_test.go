package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ianhaycox/ir-standings/connectors/iracing"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStandings(t *testing.T) {
	t.Run("should call the SearchSeriesResults and SeasonBroadcastResults services without error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		irServiceMock := iracing.NewMockIracingService(ctrl)
		irServiceMock.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return([]searchseries.SearchSeriesResult{{SessionID: 123}}, nil)
		irServiceMock.EXPECT().SeasonBroadcastResults(ctx, []searchseries.SearchSeriesResult{{SessionID: 123}}).Return([]results.Result{{SeasonID: 285}}, nil)

		actual, err := standings(ctx, irServiceMock, 2024, 2)
		assert.NoError(t, err)
		assert.Equal(t, []results.Result{{SeasonID: 285}}, actual)
	})

	t.Run("should return error if SearchSeriesResults fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		irServiceMock := iracing.NewMockIracingService(ctrl)
		irServiceMock.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return(nil, fmt.Errorf("opps"))

		_, err := standings(ctx, irServiceMock, 2024, 2)
		assert.ErrorContains(t, err, "opps")
	})

	t.Run("should return error if SeasonBroadcastResults fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		irServiceMock := iracing.NewMockIracingService(ctrl)
		irServiceMock.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return([]searchseries.SearchSeriesResult{{SessionID: 123}}, nil)
		irServiceMock.EXPECT().SeasonBroadcastResults(ctx, []searchseries.SearchSeriesResult{{SessionID: 123}}).Return(nil, fmt.Errorf("failed"))

		_, err := standings(ctx, irServiceMock, 2024, 2)
		assert.ErrorContains(t, err, "failed")
	})
}

func TestStandingsArgs(t *testing.T) {
	t.Run("args should return year and quarter", func(t *testing.T) {
		os.Args = []string{"", "2024", "2"}

		y, q, err := args()
		assert.NoError(t, err)
		assert.Equal(t, 2024, y)
		assert.Equal(t, 2, q)
	})

	t.Run("insufficient args should return error", func(t *testing.T) {
		os.Args = []string{""}

		_, _, err := args()
		assert.ErrorContains(t, err, "insufficient args")
	})

	t.Run("invalid year should return error", func(t *testing.T) {
		os.Args = []string{"", "aaa", "2"}

		_, _, err := args()
		assert.ErrorContains(t, err, "season year should be numeric")
	})

	t.Run("invalid quarter should return error", func(t *testing.T) {
		os.Args = []string{"", "2024", "bbb"}

		_, _, err := args()
		assert.ErrorContains(t, err, "season quarter should be numeric")
	})
}
