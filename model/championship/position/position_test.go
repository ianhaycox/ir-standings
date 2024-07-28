package position

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	t.Run("NewPosition returns an instance", func(t *testing.T) {
		p := NewPosition(444, true, 10, 12, 25, 13)
		assert.Equal(t, Position{subsessionID: 444, classified: true, lapsComplete: 10, position: 12, points: 25, carID: 13}, p)
		assert.Equal(t, model.LapsComplete(10), p.LapsComplete())
		assert.Equal(t, model.FinishPositionInClass(12), p.Position())
		assert.Equal(t, model.CarID(13), p.CarID())
		assert.True(t, p.IsClassified())
		assert.Equal(t, model.SubsessionID(444), p.subsessionID)
		assert.Equal(t, model.Point(25), p.Points())
	})
}

func TestBestResults(t *testing.T) {
	t.Run("BestResults should return an empty slice for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Len(t, positions.BestResults(10), 0)
	})

	t.Run("BestResults should return an empty slice for zero bestof", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, false, 1, 2, 25, 99),
			NewPosition(444, false, 2, 2, 25, 99),
		}
		assert.Len(t, positions.BestResults(0), 0)
	})

	t.Run("BestResults should return four best results", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, false, 1, 1, 25, 99),
			NewPosition(444, false, 1, 5, 25, 99),
			NewPosition(444, false, 1, 0, 25, 99),
			NewPosition(444, false, 1, 19, 25, 99),
			NewPosition(444, false, 1, 4, 25, 99),
			NewPosition(444, false, 1, 0, 25, 99),
		}

		expected := Positions{
			NewPosition(444, false, 1, 0, 25, 99),
			NewPosition(444, false, 1, 0, 25, 99),
			NewPosition(444, false, 1, 1, 25, 99),
			NewPosition(444, false, 1, 4, 25, 99),
		}

		bestPositions := positions.BestResults(4)
		assert.Len(t, bestPositions, 4)
		assert.Equal(t, expected, bestPositions)
	})

	t.Run("BestResults should return all best results because less than countOf", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, false, 1, 1, 25, 99),
			NewPosition(444, false, 1, 5, 25, 99),
			NewPosition(444, false, 1, 0, 25, 99),
		}

		expected := Positions{
			NewPosition(444, false, 1, 0, 25, 99),
			NewPosition(444, false, 1, 1, 25, 99),
			NewPosition(444, false, 1, 5, 25, 99),
		}

		bestPositions := positions.BestResults(4)
		assert.Len(t, bestPositions, 3)
		assert.Equal(t, expected, bestPositions)
	})

	t.Run("BestResults should return best results including one not counted", func(t *testing.T) {
		positions := Positions{
			NewPosition(440, false, 10, 1, 25, 99),
			NewPosition(441, false, 10, 1, 25, 99),
			NewPosition(442, false, 10, 1, model.NotCounted, 99),
		}

		expected := Positions{
			NewPosition(440, false, 10, 1, 25, 99),
			NewPosition(441, false, 10, 1, 25, 99),
			NewPosition(442, false, 10, 1, model.NotCounted, 99),
		}

		bestPositions := positions.BestResults(3)
		assert.Len(t, bestPositions, 3)
		assert.Equal(t, expected, bestPositions)
	})
}

func TestTotal(t *testing.T) {
	t.Run("Total should return zero for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Equal(t, model.Point(0), positions.Total(true, 10))
	})

	t.Run("Total should return sum of points ignoring non-counted races", func(t *testing.T) {
		positions := Positions{
			{classified: true, points: 0},
			{classified: true, points: 2},
			{classified: true, points: 1},
			{classified: true, points: 1},
			{classified: true, points: model.NotCounted},
		}

		assert.Equal(t, model.Point(0+2+1+1), positions.Total(true, 10))
	})

	t.Run("Total should return sum of points only for classified races", func(t *testing.T) {
		positions := Positions{
			{classified: true, points: 0},
			{classified: false, points: 25},
			{classified: true, points: 1},
			{classified: true, points: 1},
		}

		assert.Equal(t, model.Point(0+1+1), positions.Total(true, 10))
	})

	t.Run("Total should return sum of points for 5 best results", func(t *testing.T) {
		positions := Positions{
			{classified: true, points: 12},
			{classified: false, points: 6},
			{classified: true, points: 2},
			{classified: false, points: 13},
			{classified: true, points: model.NotCounted},
			{classified: true, points: 8},
			{classified: false, points: 7},
		}

		assert.Equal(t, model.Point(12+13+8+7+6), positions.Total(false, 5))
	})
}

func TestPositions(t *testing.T) {
	t.Run("Positions should return empty slice for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Equal(t, []standings.TieBreaker{}, positions.TieBreakerPositions(false, 10))
	})

	t.Run("Positions should return positions of all finishing positions", func(t *testing.T) {
		positions := Positions{{subsessionID: 1, position: 34}, {subsessionID: 1, position: 10}}
		assert.Equal(t, []standings.TieBreaker{standings.NewTieBreaker(1, 10), standings.NewTieBreaker(1, 34)}, positions.TieBreakerPositions(false, 10))
	})

	t.Run("Positions should return positions of classified finishing positions", func(t *testing.T) {
		positions := Positions{{subsessionID: 1, classified: true, position: 10}, {subsessionID: 1, classified: false, position: 34}}
		assert.Equal(t, []standings.TieBreaker{standings.NewTieBreaker(1, 10)}, positions.TieBreakerPositions(true, 10))
	})

	t.Run("Positions should return positions of best 2 finishing positions", func(t *testing.T) {
		positions := Positions{{subsessionID: 1, position: 10}, {position: 34}, {subsessionID: 1, position: 25}}
		assert.Equal(t, []standings.TieBreaker{standings.NewTieBreaker(1, 10), standings.NewTieBreaker(1, 25)}, positions.TieBreakerPositions(false, 2))
	})
}

func TestIsClassified(t *testing.T) {
	t.Run("IsClassified should return an empty slice for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Len(t, positions.BestResults(10), 0)
	})

	t.Run("IsClassified should return two results as classified only", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, true, 10, 1, 25, 99),
			NewPosition(444, false, 10, 5, 25, 99),
			NewPosition(444, true, 10, 3, 25, 99),
		}

		expected := Positions{
			NewPosition(444, true, 10, 1, 25, 99),
			NewPosition(444, true, 10, 3, 25, 99),
		}

		isClassified := positions.Classified(true)
		assert.Len(t, isClassified, 2)
		assert.Equal(t, expected, isClassified)
	})

	t.Run("IsClassified should return all results as classified only == false", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, true, 10, 1, 25, 99),
			NewPosition(444, false, 10, 5, 25, 99),
			NewPosition(444, true, 10, 3, 25, 99),
		}

		expected := Positions{
			NewPosition(444, true, 10, 1, 25, 99),
			NewPosition(444, false, 10, 5, 25, 99),
			NewPosition(444, true, 10, 3, 25, 99),
		}

		isClassified := positions.Classified(false)
		assert.Len(t, isClassified, 3)
		assert.Equal(t, expected, isClassified)
	})
}

func TestLaps(t *testing.T) {
	positions := Positions{
		NewPosition(444, true, 10, 1, 25, 99),
		NewPosition(444, false, 10, 5, 25, 99),
		NewPosition(444, true, 10, 3, 25, 99),
	}

	t.Run("Laps should return zero for empty input", func(t *testing.T) {
		empty := Positions{}
		assert.Equal(t, model.LapsComplete(0), empty.Laps(true, 10))
	})

	t.Run("Laps should total laps for all finishes", func(t *testing.T) {
		assert.Equal(t, model.LapsComplete(30), positions.Laps(false, 10))
	})

	t.Run("Laps should total laps for classified finishes", func(t *testing.T) {
		assert.Equal(t, model.LapsComplete(20), positions.Laps(true, 10))
	})

	t.Run("Laps should total laps for best 1 classified finishes", func(t *testing.T) {
		assert.Equal(t, model.LapsComplete(10), positions.Laps(true, 1))
	})
}

func TestCounted(t *testing.T) {
	positions := Positions{
		NewPosition(444, true, 10, 1, 25, 99),
		NewPosition(444, false, 10, 5, 25, 99),
		NewPosition(444, true, 10, 3, 25, 99),
	}

	t.Run("Counted should zero for empty input", func(t *testing.T) {
		empty := Positions{}
		assert.Equal(t, 0, empty.Counted(true, 10))
	})

	t.Run("Counted should total races for all finishes", func(t *testing.T) {
		assert.Equal(t, 3, positions.Counted(false, 10))
	})

	t.Run("Counted should total races for classified finishes", func(t *testing.T) {
		assert.Equal(t, 2, positions.Counted(true, 10))
	})

	t.Run("Counted should total races for best 1 classified finishes", func(t *testing.T) {
		assert.Equal(t, 1, positions.Counted(true, 1))
	})
}

func TestCarsDriven(t *testing.T) {
	positions := Positions{
		NewPosition(444, true, 10, 1, 25, 99),
		NewPosition(444, true, 10, 2, 25, 99),
		NewPosition(444, false, 10, 5, 25, 98),
		NewPosition(444, true, 10, 3, 25, 97),
	}

	t.Run("CarsDriven should return an empty slice for empty input", func(t *testing.T) {
		empty := Positions{}
		assert.Len(t, empty.CarsDriven(true, 10), 0)
	})

	t.Run("CarsDriven should return all cars for all finishes", func(t *testing.T) {
		assert.Equal(t, []model.CarID{99, 98, 97}, positions.CarsDriven(false, 10))
	})

	t.Run("CarsDriven should return cars for classified finishes", func(t *testing.T) {
		assert.Equal(t, []model.CarID{99, 97}, positions.CarsDriven(true, 10))
	})

	t.Run("CarsDriven should return cars for best 1 classified finishes", func(t *testing.T) {
		assert.Equal(t, []model.CarID{99}, positions.CarsDriven(true, 1))
	})
}
