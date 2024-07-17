package cdn

import (
	"testing"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/stretchr/testify/assert"
)

func TestCDN(t *testing.T) {
	t.Run("Should retrun an instance of CDNService", func(t *testing.T) {
		c := NewCDNService(api.NewHTTPClient(api.NewConfiguration(nil, "")))
		assert.NotNil(t, c)
	})
}
