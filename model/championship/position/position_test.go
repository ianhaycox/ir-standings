package position

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	t.Run("NewPosition returns an instance", func(t *testing.T) {
		p := NewPosition(444, true, 10, 11, 12, 13)
		assert.Equal(t, Position{subsessionID: 444, classified: true, lapsComplete: 10, splitNum: 11, position: 12, carID: 13}, p)
		assert.Equal(t, 10, p.LapsComplete())
		assert.Equal(t, model.SplitNum(11), p.SplitNum())
		assert.Equal(t, 12, p.Position())
		assert.Equal(t, model.CarID(13), p.CarID())
		assert.True(t, p.IsClassified())
		assert.Equal(t, model.SubsessionID(444), p.subsessionID)
	})
}

func TestBestPosition(t *testing.T) {
	t.Run("BestPosition should return an empty slice for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Len(t, positions.BestPositions(10), 0)
	})

	t.Run("BestPosition should return an empty slice for zero bestof", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, false, 1, 1, 2, 99),
			NewPosition(444, false, 2, 2, 2, 99),
		}
		assert.Len(t, positions.BestPositions(0), 0)
	})

	t.Run("BestPosition should return four best results - single split", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, false, 1, 0, 1, 99),
			NewPosition(444, false, 1, 0, 5, 99),
			NewPosition(444, false, 1, 0, 0, 99),
			NewPosition(444, false, 1, 0, 19, 99),
			NewPosition(444, false, 1, 0, 4, 99),
			NewPosition(444, false, 1, 0, 0, 99),
		}

		expected := Positions{
			NewPosition(444, false, 1, 0, 0, 99),
			NewPosition(444, false, 1, 0, 0, 99),
			NewPosition(444, false, 1, 0, 1, 99),
			NewPosition(444, false, 1, 0, 4, 99),
		}

		bestPositions := positions.BestPositions(4)
		assert.Len(t, bestPositions, 4)
		assert.Equal(t, expected, bestPositions)
	})

	t.Run("BestPosition should return all best results because less than countOf - single split", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, false, 1, 0, 1, 99),
			NewPosition(444, false, 1, 0, 5, 99),
			NewPosition(444, false, 1, 0, 0, 99),
		}

		expected := Positions{
			NewPosition(444, false, 1, 0, 0, 99),
			NewPosition(444, false, 1, 0, 1, 99),
			NewPosition(444, false, 1, 0, 5, 99),
		}

		bestPositions := positions.BestPositions(4)
		assert.Len(t, bestPositions, 3)
		assert.Equal(t, expected, bestPositions)
	})

	t.Run("BestPosition should return best results - multiple splits - sorted position then split", func(t *testing.T) {
		positions := Positions{
			NewPosition(440, false, 10, 0, 1, 99),
			NewPosition(441, false, 10, 1, 1, 99),
			NewPosition(442, false, 10, 2, 1, 99),
			NewPosition(441, false, 10, 1, 3, 99),
			NewPosition(440, false, 10, 0, 1, 99),
		}

		expected := Positions{
			NewPosition(440, false, 10, 0, 1, 99),
			NewPosition(440, false, 10, 0, 1, 99),
			NewPosition(441, false, 10, 1, 1, 99),
		}

		bestPositions := positions.BestPositions(3)
		assert.Len(t, bestPositions, 3)
		assert.Equal(t, expected, bestPositions)
	})
}

func TestTotal(t *testing.T) {
	t.Run("Total should return zero for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Equal(t, 0, positions.Total(points.NewPointsStructure(make(points.PointsPerSplit)), true, 10))
	})

	t.Run("Total should return sum of points according to the points structure per split", func(t *testing.T) {
		awards := points.PointsPerSplit{
			0: {25, 22, 20, 18},
			1: {14, 12, 10},
			2: {9, 6},
		}

		ps := points.NewPointsStructure(awards)

		positions := Positions{
			{classified: true, splitNum: 0, position: 0},
			{classified: true, splitNum: 1, position: 2},
			{classified: true, splitNum: 2, position: 1},
			{classified: true, splitNum: 3, position: 1}, // Ignored due to split 3
		}

		assert.Equal(t, 25+10+6, positions.Total(ps, true, 10))
	})

	t.Run("Total should return sum of points only for classified races", func(t *testing.T) {
		awards := points.PointsPerSplit{
			0: {25, 22, 20, 18},
			1: {14, 12, 10},
			2: {9, 6},
		}

		ps := points.NewPointsStructure(awards)

		positions := Positions{
			{classified: true, splitNum: 0, position: 0},
			{classified: false, splitNum: 1, position: 25},
			{classified: true, splitNum: 2, position: 1},
			{classified: true, splitNum: 0, position: 1},
		}

		assert.Equal(t, 25+22+6, positions.Total(ps, true, 10))
	})

	t.Run("Total should return sum of points for all races", func(t *testing.T) {
		awards := points.PointsPerSplit{
			0: {25, 22, 20, 18},
			1: {14, 12, 10},
			2: {9, 6},
		}

		ps := points.NewPointsStructure(awards)

		positions := Positions{
			{classified: true, splitNum: 0, position: 0},
			{classified: false, splitNum: 1, position: 1},
			{classified: true, splitNum: 2, position: 0},
			{classified: false, splitNum: 0, position: 3},
		}

		assert.Equal(t, 25+12+9+18, positions.Total(ps, false, 10))
	})
}

func TestPositions(t *testing.T) {
	t.Run("Positions should return empty slice for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Equal(t, []standings.TieBreaker{}, positions.Positions(false, 10))
	})

	t.Run("Positions should return positions of all finishing positions", func(t *testing.T) {
		positions := Positions{{subsessionID: 1, position: 34}, {subsessionID: 1, position: 10}}
		assert.Equal(t, []standings.TieBreaker{standings.NewTieBreaker(1, 10), standings.NewTieBreaker(1, 34)}, positions.Positions(false, 10))
	})

	t.Run("Positions should return positions of classified finishing positions", func(t *testing.T) {
		positions := Positions{{subsessionID: 1, classified: true, position: 10}, {subsessionID: 1, classified: false, position: 34}}
		assert.Equal(t, []standings.TieBreaker{standings.NewTieBreaker(1, 10)}, positions.Positions(true, 10))
	})

	t.Run("Positions should return positions of best 2 finishing positions", func(t *testing.T) {
		positions := Positions{{subsessionID: 1, position: 10}, {position: 34}, {subsessionID: 1, position: 25}}
		assert.Equal(t, []standings.TieBreaker{standings.NewTieBreaker(1, 10), standings.NewTieBreaker(1, 25)}, positions.Positions(false, 2))
	})
}

func TestIsClassified(t *testing.T) {
	t.Run("IsClassified should return an empty slice for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Len(t, positions.BestPositions(10), 0)
	})

	t.Run("IsClassified should return two results as classified only", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, true, 10, 0, 1, 99),
			NewPosition(444, false, 10, 0, 5, 99),
			NewPosition(444, true, 10, 0, 3, 99),
		}

		expected := Positions{
			NewPosition(444, true, 10, 0, 1, 99),
			NewPosition(444, true, 10, 0, 3, 99),
		}

		isClassified := positions.Classified(true)
		assert.Len(t, isClassified, 2)
		assert.Equal(t, expected, isClassified)
	})

	t.Run("IsClassified should return all results as classified only == false", func(t *testing.T) {
		positions := Positions{
			NewPosition(444, true, 10, 0, 1, 99),
			NewPosition(444, false, 10, 0, 5, 99),
			NewPosition(444, true, 10, 0, 3, 99),
		}

		expected := Positions{
			NewPosition(444, true, 10, 0, 1, 99),
			NewPosition(444, false, 10, 0, 5, 99),
			NewPosition(444, true, 10, 0, 3, 99),
		}

		isClassified := positions.Classified(false)
		assert.Len(t, isClassified, 3)
		assert.Equal(t, expected, isClassified)
	})
}

func TestLaps(t *testing.T) {
	positions := Positions{
		NewPosition(444, true, 10, 0, 1, 99),
		NewPosition(444, false, 10, 0, 5, 99),
		NewPosition(444, true, 10, 0, 3, 99),
	}

	t.Run("Laps should return zero for empty input", func(t *testing.T) {
		empty := Positions{}
		assert.Equal(t, 0, empty.Laps(true, 10))
	})

	t.Run("Laps should total laps for all finishes", func(t *testing.T) {
		assert.Equal(t, 30, positions.Laps(false, 10))
	})

	t.Run("Laps should total laps for classified finishes", func(t *testing.T) {
		assert.Equal(t, 20, positions.Laps(true, 10))
	})

	t.Run("Laps should total laps for best 1 classified finishes", func(t *testing.T) {
		assert.Equal(t, 10, positions.Laps(true, 1))
	})
}

func TestCounted(t *testing.T) {
	positions := Positions{
		NewPosition(444, true, 10, 0, 1, 99),
		NewPosition(444, false, 10, 0, 5, 99),
		NewPosition(444, true, 10, 0, 3, 99),
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
		NewPosition(444, true, 10, 0, 1, 99),
		NewPosition(444, true, 10, 0, 2, 99),
		NewPosition(444, false, 10, 0, 5, 98),
		NewPosition(444, true, 10, 0, 3, 97),
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
