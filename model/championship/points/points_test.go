package points

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/stretchr/testify/assert"
)

func TestPoints(t *testing.T) {
	t.Run("Empty structure returns not counted", func(t *testing.T) {
		points := PointsPerSplit{}

		ps := NewPointsStructure(points)
		assert.Equal(t, model.NotCounted, ps.Award(0, 10, 1))
	})

	t.Run("Structure returns corresponding points", func(t *testing.T) {
		points := PointsPerSplit{
			0: {25, 22, 20, 18},
			1: {14, 12, 10},
			2: {9, 6},
		}

		ps := NewPointsStructure(points)
		assert.Equal(t, model.Point(25), ps.Award(0, 0, 1))
		assert.Equal(t, model.Point(18), ps.Award(0, 3, 1))
		assert.Equal(t, model.Point(0), ps.Award(0, 4, 1))

		assert.Equal(t, model.Point(14), ps.Award(1, 0, 1))
		assert.Equal(t, model.Point(9), ps.Award(2, 0, 1))
		assert.Equal(t, model.Point(0), ps.Award(1, 10, 1))

		assert.Equal(t, model.NotCounted, ps.Award(3, 10, 1))
		assert.Equal(t, model.NotCounted, ps.Award(8, 10, 1))

		// Race not started yet
		assert.Equal(t, model.Point(0), ps.Award(0, 0, 0))
	})
}
