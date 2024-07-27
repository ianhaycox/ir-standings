package race

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/position"
	"github.com/ianhaycox/ir-standings/model/championship/result"
	"github.com/stretchr/testify/assert"
)

func TestRaceWinnerLapComplete(t *testing.T) {
	t.Run("Should return zero for no results", func(t *testing.T) {
		race := NewRace(1, 2, []result.Result{})

		actual := race.WinnerLapsComplete(333)
		assert.Equal(t, model.LapsComplete(0), actual)
	})

	t.Run("Should return max laps complete per class", func(t *testing.T) {
		results := []result.Result{
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
		assert.Equal(t, model.LapsComplete(10), actual)

		actual = race.WinnerLapsComplete(2)
		assert.Equal(t, model.LapsComplete(30), actual)

		actual = race.WinnerLapsComplete(3)
		assert.Equal(t, model.LapsComplete(12), actual)
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
	awards := points.PointsPerSplit{
		0: {25, 22, 20, 18},
		1: {14, 12, 10},
		2: {9, 6},
	}

	ps := points.NewPointsStructure(awards)

	results := []result.Result{
		{SubsessionID: 444, CarClassID: 1, CustID: 1777, LapsComplete: 13, FinishPositionInClass: 0, CarID: 97},
		{SubsessionID: 444, CarClassID: 1, CustID: 1888, LapsComplete: 13, FinishPositionInClass: 2, CarID: 97},
		{SubsessionID: 444, CarClassID: 1, CustID: 1999, LapsComplete: 12, FinishPositionInClass: 3, CarID: 97},
		{SubsessionID: 444, CarClassID: 2, CustID: 2111, LapsComplete: 20, FinishPositionInClass: 0, CarID: 98},
		{SubsessionID: 444, CarClassID: 2, CustID: 2222, LapsComplete: 9, FinishPositionInClass: 1, CarID: 98},
		{SubsessionID: 444, CarClassID: 2, CustID: 2333, LapsComplete: 9, FinishPositionInClass: 2, CarID: 98},
		{SubsessionID: 444, CarClassID: 3, CustID: 3333, LapsComplete: 6, FinishPositionInClass: 1, CarID: 99},
		{SubsessionID: 444, CarClassID: 3, CustID: 3444, LapsComplete: 8, FinishPositionInClass: 2, CarID: 99},
	}

	t.Run("Should return empty for no results", func(t *testing.T) {
		race := Race{}

		actual := race.Positions(1, 1, ps)
		assert.Equal(t, make(map[model.CustID]position.Position), actual)
	})

	t.Run("Should return positions and points per class and cust for split 0 (top split)", func(t *testing.T) {
		race := Race{
			splitNum: 0,
			results:  results,
		}

		actual := race.Positions(model.CarClassID(1), model.LapsComplete(10), ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				1777: position.NewPosition(444, true, 13, 0, 25, 97),
				1888: position.NewPosition(444, true, 13, 2, 20, 97),
				1999: position.NewPosition(444, true, 12, 3, 18, 97),
			}, actual)

		actual = race.Positions(model.CarClassID(2), model.LapsComplete(20), ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				2111: position.NewPosition(444, true, 20, 0, 25, 98),
				2222: position.NewPosition(444, false, 9, 1, 22, 98),
				2333: position.NewPosition(444, false, 9, 2, 20, 98),
			}, actual)

		actual = race.Positions(model.CarClassID(3), model.LapsComplete(30), ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				3333: position.NewPosition(444, false, 6, 1, 22, 99),
				3444: position.NewPosition(444, false, 8, 2, 20, 99),
			}, actual)
	})

	t.Run("Should return positions and points per class and cust for split 1 (second split)", func(t *testing.T) {
		race := Race{
			splitNum: 1,
			results:  results,
		}

		actual := race.Positions(model.CarClassID(1), model.LapsComplete(10), ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				1777: position.NewPosition(444, true, 13, 0, 14, 97),
				1888: position.NewPosition(444, true, 13, 2, 10, 97),
				1999: position.NewPosition(444, true, 12, 3, 0, 97),
			}, actual)

		actual = race.Positions(model.CarClassID(2), model.LapsComplete(20), ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				2111: position.NewPosition(444, true, 20, 0, 14, 98),
				2222: position.NewPosition(444, false, 9, 1, 12, 98),
				2333: position.NewPosition(444, false, 9, 2, 10, 98),
			}, actual)

		actual = race.Positions(model.CarClassID(3), model.LapsComplete(30), ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				3333: position.NewPosition(444, false, 6, 1, 12, 99),
				3444: position.NewPosition(444, false, 8, 2, 10, 99),
			}, actual)
	})

	t.Run("Should return positions and points per class and cust for split 2 (third split)", func(t *testing.T) {
		race := Race{
			splitNum: 2,
			results:  results,
		}

		actual := race.Positions(1, 10, ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				1777: position.NewPosition(444, true, 13, 0, 9, 97),
				1888: position.NewPosition(444, true, 13, 2, 0, 97),
				1999: position.NewPosition(444, true, 12, 3, 0, 97),
			}, actual)

		actual = race.Positions(2, 20, ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				2111: position.NewPosition(444, true, 20, 0, 9, 98),
				2222: position.NewPosition(444, false, 9, 1, 6, 98),
				2333: position.NewPosition(444, false, 9, 2, 0, 98),
			}, actual)

		actual = race.Positions(3, 30, ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				3333: position.NewPosition(444, false, 6, 1, 6, 99),
				3444: position.NewPosition(444, false, 8, 2, 0, 99),
			}, actual)
	})

	t.Run("Should return positions and NOT_COUNTED points per class and cust for split 3 (fourth split)", func(t *testing.T) {
		race := Race{
			splitNum: 3,
			results:  results,
		}

		actual := race.Positions(1, 10, ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				1777: position.NewPosition(444, true, 13, 0, model.NotCounted, 97),
				1888: position.NewPosition(444, true, 13, 2, model.NotCounted, 97),
				1999: position.NewPosition(444, true, 12, 3, model.NotCounted, 97),
			}, actual)

		actual = race.Positions(2, 20, ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				2111: position.NewPosition(444, true, 20, 0, model.NotCounted, 98),
				2222: position.NewPosition(444, false, 9, 1, model.NotCounted, 98),
				2333: position.NewPosition(444, false, 9, 2, model.NotCounted, 98),
			}, actual)

		actual = race.Positions(3, 30, ps)
		assert.Equal(t,
			map[model.CustID]position.Position{
				3333: position.NewPosition(444, false, 6, 1, model.NotCounted, 99),
				3444: position.NewPosition(444, false, 8, 2, model.NotCounted, 99),
			}, actual)
	})
}

func TestSplitNum(t *testing.T) {
	t.Run("Should return split number", func(t *testing.T) {
		race := NewRace(2, 1, []result.Result{})

		assert.Equal(t, model.SplitNum(2), race.SplitNum())
	})
}
