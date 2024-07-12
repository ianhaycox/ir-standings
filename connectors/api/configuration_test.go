package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDefaultHeader(t *testing.T) {
	t.Parallel()

	t.Run("add header", func(t *testing.T) {
		c := NewConfiguration("https://example.com", nil)
		c.AddDefaultHeader("key", "value")
		assert.Equal(t, map[string]string{"Accept": "application/json", "Content-Type": "application/json", "key": "value"}, c.DefaultHeader)
		assert.Equal(t, "https://example.com", c.BasePath)
	})
}
