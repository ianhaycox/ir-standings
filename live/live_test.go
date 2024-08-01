package live

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/stretchr/testify/assert"
)

func TestLive(t *testing.T) {
	t.Run("New standing each lap", func(t *testing.T) {
		livePositions := live.LiveResults{
			SeriesID:     285,
			SessionID:    999,
			SubsessionID: 123,
			Track:        "new",
			CountBestOf:  10,
			CarClassID:   84,
			TopN:         4,
			Positions: []live.CurrentPosition{
				{CustID: 341977, FinishPositionInClass: 3, LapsComplete: 10, CarID: 77},
				{CustID: 86672, FinishPositionInClass: 2, LapsComplete: 10, CarID: 77},
				{CustID: 667693, FinishPositionInClass: 1, LapsComplete: 10, CarID: 77},
				{CustID: 301749, FinishPositionInClass: 0, LapsComplete: 10, CarID: 77},
			},
		}

		b, err := json.MarshalIndent(livePositions, "", "  ")
		assert.NoError(t, err)

		s, err := Live(string(b))
		assert.NoError(t, err)

		fmt.Println(s)

		var prediction1 []live.LiveResults

		err = json.Unmarshal([]byte(s), &prediction1)
		assert.NoError(t, err)

		livePositions2 := live.LiveResults{
			SeriesID:     285,
			SessionID:    888,
			SubsessionID: 321,
			Track:        "new",
			CountBestOf:  10,
			CarClassID:   84,
			TopN:         5,
			Positions: []live.CurrentPosition{
				{CustID: 341977, FinishPositionInClass: 0, LapsComplete: 10, CarID: 77},
				{CustID: 86672, FinishPositionInClass: 1, LapsComplete: 10, CarID: 77},
				{CustID: 667693, FinishPositionInClass: 2, LapsComplete: 10, CarID: 77},
				{CustID: 301749, FinishPositionInClass: 3, LapsComplete: 10, CarID: 77},
			},
		}

		b, err = json.MarshalIndent(livePositions2, "", "  ")
		assert.NoError(t, err)

		s, err = Live(string(b))
		assert.NoError(t, err)

		fmt.Println(s)

		var prediction2 []live.LiveResults

		err = json.Unmarshal([]byte(s), &prediction2)
		assert.NoError(t, err)
	})
}
