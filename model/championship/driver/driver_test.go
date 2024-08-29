package driver

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	d := NewDriver(23, "test", 2354)
	assert.Equal(t, model.CustID(23), d.CustID())
	assert.Equal(t, "test", d.DisplayName())
	assert.Equal(t, 2354, d.IRating())
}
