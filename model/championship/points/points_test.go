package points

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoints(t *testing.T) {
	t.Run("Empty structure returns zero", func(t *testing.T) {
		points := PointsPerSplit{}

		ps := NewPointsStructure(points)
		assert.Equal(t, 0, ps.Award(0, 10))
	})

	t.Run("Structure returns corresponding points", func(t *testing.T) {
		points := PointsPerSplit{
			0: {25, 22, 20, 18},
			1: {14, 12, 10},
			2: {9, 6},
		}

		ps := NewPointsStructure(points)
		assert.Equal(t, 25, ps.Award(0, 0))
		assert.Equal(t, 18, ps.Award(0, 3))
		assert.Equal(t, 0, ps.Award(0, 4))

		assert.Equal(t, 14, ps.Award(1, 0))
		assert.Equal(t, 9, ps.Award(2, 0))
		assert.Equal(t, 0, ps.Award(1, 10))
	})
}
