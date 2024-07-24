package position

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	t.Run("NewPosition returns an instance", func(t *testing.T) {
		p := NewPosition(10, 11, 12, 13)
		assert.Equal(t, Position{lapsComplete: 10, splitNum: 11, position: 12, carID: 13}, p)
		assert.Equal(t, 10, p.LapsComplete())
		assert.Equal(t, model.SplitNum(11), p.SplitNum())
		assert.Equal(t, 12, p.Position())
		assert.Equal(t, model.CarID(13), p.CarID())
	})
}

func TestBestPosition(t *testing.T) {
	t.Run("BestPosition should return an empty array for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Len(t, positions.BestPositions(10), 0)
	})

	t.Run("BestPosition should return an empty array for zero bestof", func(t *testing.T) {
		positions := Positions{
			NewPosition(1, 1, 2, 99),
			NewPosition(2, 2, 2, 99),
		}
		assert.Len(t, positions.BestPositions(0), 0)
	})

	t.Run("BestPosition should return four best results", func(t *testing.T) {
		positions := Positions{
			NewPosition(1, 0, 1, 99),
			NewPosition(1, 0, 5, 99),
			NewPosition(1, 0, 0, 99),
			NewPosition(1, 0, 19, 99),
			NewPosition(1, 0, 4, 99),
			NewPosition(1, 0, 0, 99),
		}

		expected := Positions{
			NewPosition(1, 0, 0, 99),
			NewPosition(1, 0, 0, 99),
			NewPosition(1, 0, 1, 99),
			NewPosition(1, 0, 4, 99),
		}

		bestPositions := positions.BestPositions(4)
		assert.Len(t, bestPositions, 4)
		assert.Equal(t, expected, bestPositions)
	})

	t.Run("BestPosition should return all best results because less than countOf", func(t *testing.T) {
		positions := Positions{
			NewPosition(1, 0, 1, 99),
			NewPosition(1, 0, 5, 99),
			NewPosition(1, 0, 0, 99),
		}

		expected := Positions{
			NewPosition(1, 0, 0, 99),
			NewPosition(1, 0, 1, 99),
			NewPosition(1, 0, 5, 99),
		}

		bestPositions := positions.BestPositions(4)
		assert.Len(t, bestPositions, 3)
		assert.Equal(t, expected, bestPositions)
	})
}

func TestTotal(t *testing.T) {
	t.Run("Total should return zero for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Equal(t, 0, positions.Total(points.NewPointsStructure(make(points.PointsPerSplit))))
	})

	t.Run("Total should return sum of points according to the points structure per split", func(t *testing.T) {
		awards := points.PointsPerSplit{
			0: {25, 22, 20, 18},
			1: {14, 12, 10},
			2: {9, 6},
		}

		ps := points.NewPointsStructure(awards)

		positions := Positions{
			{splitNum: 0, position: 0},
			{splitNum: 1, position: 2},
			{splitNum: 2, position: 1},
			{splitNum: 3, position: 1}, // Ignored
		}

		assert.Equal(t, 41, positions.Total(ps))
	})
}

func TestSum(t *testing.T) {
	t.Run("Sum should return zero for empty input", func(t *testing.T) {
		positions := Positions{}
		assert.Equal(t, 0, positions.Sum())
	})

	t.Run("Sum should return summation of finishing positions", func(t *testing.T) {
		positions := Positions{{position: 10}, {position: 34}}
		assert.Equal(t, 44, positions.Sum())
	})
}
