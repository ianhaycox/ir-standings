package live

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePositionsTestData(t *testing.T) {
	p := LiveResults{
		SeriesID:     285,
		SessionID:    1,
		SubsessionID: 2,
		Track:        "Lime Rock",
		CountBestOf:  10,
		CarClassID:   84,
		TopN:         5,
		Positions: []CurrentPosition{
			{
				CustID:                123,
				FinishPositionInClass: 0,
				LapsComplete:          10,
				CarID:                 76,
			},
		},
	}

	b, err := json.Marshal(p)
	assert.NoError(t, err)
	fmt.Println(string(b))
}
