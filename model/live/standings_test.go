package live

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTestData(t *testing.T) {
	s := []PredictedStanding{
		{
			DriverName:        "D1",
			CurrentPosition:   1,
			PredictedPosition: 2,
			CurrentPoints:     123,
			PredictedPoints:   120,
			Change:            -1,
		},
		{
			DriverName:        "D2",
			CurrentPosition:   2,
			PredictedPosition: 1,
			CurrentPoints:     120,
			PredictedPoints:   123,
			Change:            1,
		},
	}

	b, err := json.Marshal(s)
	assert.NoError(t, err)
	fmt.Println(string(b))
}
