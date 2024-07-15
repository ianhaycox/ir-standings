package iracing

import (
	"testing"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/stretchr/testify/assert"
)

func TestIRacing(t *testing.T) {
	t.Run("should return an IRacingService instance", func(t *testing.T) {
		i := NewIracingService(
			NewIracingDataService(api.NewAPIClient(api.NewConfiguration(nil, ""))),
			api.NewAuthenticationService("", ""),
		)
		assert.NotNil(t, i)
	})
}
