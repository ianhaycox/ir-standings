package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResults(t *testing.T) {
	t.Run("IsClassified if completed 75 percent or more laps", func(t *testing.T) {
		res := Result{LapsComplete: 10}

		assert.True(t, res.IsClassified(10))
		assert.True(t, res.IsClassified(11))
		assert.True(t, res.IsClassified(12))
		assert.True(t, res.IsClassified(13))
		assert.False(t, res.IsClassified(14))
		assert.False(t, res.IsClassified(15))
		assert.False(t, res.IsClassified(16))
	})
}
