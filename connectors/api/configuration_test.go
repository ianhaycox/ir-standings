package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDefaultHeader(t *testing.T) {
	t.Parallel()

	t.Run("add header", func(t *testing.T) {
		c := NewConfiguration(nil, "agent")
		c.AddDefaultHeader("key", "value")
		c.AddDefaultHeader("Accept", "application/json")
		c.AddDefaultHeader("Content-Type", "application/json")
		assert.Equal(t, map[string]string{"Accept": "application/json", "Content-Type": "application/json", "key": "value"}, c.DefaultHeader)
		assert.Equal(t, "agent", c.UserAgent)
	})
}
