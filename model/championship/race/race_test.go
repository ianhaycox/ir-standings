package race

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/position"
	"github.com/stretchr/testify/assert"
)

func TestRaceWinnerLapComplete(t *testing.T) {
	t.Run("Should return zero for no results", func(t *testing.T) {
		race := NewRace(1, 2, []model.Result{})

		actual := race.WinnerLapsComplete(333)
		assert.Equal(t, 0, actual)
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

		actual := race.WinnerLapsComplete(1)
		assert.Equal(t, 10, actual)

		actual = race.WinnerLapsComplete(2)
		assert.Equal(t, 30, actual)

		actual = race.WinnerLapsComplete(3)
		assert.Equal(t, 12, actual)
	})
}

func TestIsClassified(t *testing.T) {
	t.Run("IsClassified if completed 75 percent or more laps", func(t *testing.T) {
		res := Race{}

		assert.True(t, res.IsClassified(10, 10))
		assert.True(t, res.IsClassified(11, 10))
		assert.True(t, res.IsClassified(12, 10))
		assert.True(t, res.IsClassified(13, 10))
		assert.False(t, res.IsClassified(14, 10))
		assert.False(t, res.IsClassified(15, 10))
		assert.False(t, res.IsClassified(16, 10))
	})
}

func TestFinishingPositions(t *testing.T) {
	t.Run("Should return empty for no results", func(t *testing.T) {
		race := Race{}

		actual := race.Positions(1, 1)
		assert.Equal(t, make(map[model.CustID]position.Position), actual)
	})

	t.Run("Should return positions and points per class and cust", func(t *testing.T) {
		results := []model.Result{
			{SubsessionID: 444, CarClassID: 1, CustID: 1777, LapsComplete: 13, FinishPositionInClass: 1, CarID: 97},
			{SubsessionID: 444, CarClassID: 3, CustID: 3333, LapsComplete: 6, FinishPositionInClass: 1, CarID: 99},
			{SubsessionID: 444, CarClassID: 1, CustID: 1999, LapsComplete: 12, FinishPositionInClass: 3, CarID: 97},
			{SubsessionID: 444, CarClassID: 2, CustID: 2111, LapsComplete: 10, FinishPositionInClass: 1, CarID: 98},
			{SubsessionID: 444, CarClassID: 1, CustID: 1888, LapsComplete: 13, FinishPositionInClass: 2, CarID: 97},
			{SubsessionID: 444, CarClassID: 2, CustID: 2222, LapsComplete: 9, FinishPositionInClass: 2, CarID: 98},
			{SubsessionID: 444, CarClassID: 2, CustID: 2333, LapsComplete: 9, FinishPositionInClass: 3, CarID: 98},
			{SubsessionID: 444, CarClassID: 3, CustID: 3444, LapsComplete: 8, FinishPositionInClass: 2, CarID: 99},
		}

		race := Race{
			splitNum: 1,
			results:  results,
		}

		actual := race.Positions(1, 10)
		assert.Equal(t,
			map[model.CustID]position.Position{
				1777: position.NewPosition(444, true, 13, 1, 1, 97),
				1888: position.NewPosition(444, true, 13, 1, 2, 97),
				1999: position.NewPosition(444, true, 12, 1, 3, 97),
			}, actual)

		actual = race.Positions(2, 20)
		assert.Equal(t,
			map[model.CustID]position.Position{
				2111: position.NewPosition(444, false, 10, 1, 1, 98),
				2222: position.NewPosition(444, false, 9, 1, 2, 98),
				2333: position.NewPosition(444, false, 9, 1, 3, 98),
			}, actual)

		actual = race.Positions(3, 10)
		assert.Equal(t,
			map[model.CustID]position.Position{
				3333: position.NewPosition(444, false, 6, 1, 1, 99),
				3444: position.NewPosition(444, true, 8, 1, 2, 99),
			}, actual)
	})
}

func TestSplitNum(t *testing.T) {
	t.Run("Should return split number", func(t *testing.T) {
		race := NewRace(2, 1, []model.Result{})

		assert.Equal(t, model.SplitNum(2), race.SplitNum())
	})
}
