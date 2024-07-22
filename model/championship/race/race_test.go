package race

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/stretchr/testify/assert"
)

func TestRaceWinnerLapComplete(t *testing.T) {
	t.Run("Should return empty for no results", func(t *testing.T) {
		race := NewRace(1, 2, []model.Result{})

		actual := race.WinnerLapsComplete()
		assert.Equal(t, make(map[model.CarClassID]int), actual)
	})

	t.Run("Should return max laps complete per class", func(t *testing.T) {
		results := []model.Result{
			{CarClassID: 1, LapsComplete: 8},
			{CarClassID: 2, LapsComplete: 30},
			{CarClassID: 1, LapsComplete: 10},
			{CarClassID: 3, LapsComplete: 12},
			{CarClassID: 3, LapsComplete: 11},
			{CarClassID: 2, LapsComplete: 30},
			{CarClassID: 1, LapsComplete: 9},
			{CarClassID: 2, LapsComplete: 30},
		}

		race := Race{
			results: results,
		}

		expected := map[model.CarClassID]int{
			1: 10,
			2: 30,
			3: 12,
		}
		actual := race.WinnerLapsComplete()
		assert.Equal(t, expected, actual)
	})
}

func TestFinishingPositions(t *testing.T) {
	t.Run("Should return empty for no results", func(t *testing.T) {
		race := Race{}

		actual := race.Positions()
		assert.Equal(t, make(map[model.CarClassID]map[model.CustID]Position), actual)
	})

	t.Run("Should return positions and points per class and cust", func(t *testing.T) {
		results := []model.Result{
			{CarClassID: 1, CustID: 1777, LapsComplete: 13, FinishPositionInClass: 1},
			{CarClassID: 1, CustID: 1888, LapsComplete: 13, FinishPositionInClass: 2},
			{CarClassID: 1, CustID: 1999, LapsComplete: 12, FinishPositionInClass: 3},
			{CarClassID: 2, CustID: 2111, LapsComplete: 10, FinishPositionInClass: 1},
			{CarClassID: 2, CustID: 2222, LapsComplete: 9, FinishPositionInClass: 2},
			{CarClassID: 2, CustID: 2333, LapsComplete: 9, FinishPositionInClass: 3},
			{CarClassID: 3, CustID: 3333, LapsComplete: 6, FinishPositionInClass: 1},
			{CarClassID: 3, CustID: 3444, LapsComplete: 6, FinishPositionInClass: 2},
		}

		race := Race{
			splitNum: 1,
			results:  results,
		}

		expected := map[model.CarClassID]map[model.CustID]Position{
			1: {
				1777: NewPosition(13, 1, 1),
				1888: NewPosition(13, 1, 2),
				1999: NewPosition(12, 1, 3),
			},
			2: {
				2111: NewPosition(10, 1, 1),
				2222: NewPosition(9, 1, 2),
				2333: NewPosition(9, 1, 3),
			},
			3: {
				3333: NewPosition(6, 1, 1),
				3444: NewPosition(6, 1, 2),
			},
		}
		actual := race.Positions()
		assert.Equal(t, expected, actual)
	})
}
