package live

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/stretchr/testify/assert"
)

var ex = `{
  "car_class_id": 83,
  "count_best_of": 12,
  "results": [{
      "car_id": 76,
      "cust_id": 141372,
      "finish_position_in_class": 6,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 321008,
      "finish_position_in_class": 2,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 467784,
      "finish_position_in_class": 1,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 637832,
      "finish_position_in_class": 4,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 851509,
      "finish_position_in_class": 3,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 173387,
      "finish_position_in_class": 5,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 390264,
      "finish_position_in_class": 17,
      "laps_complete": -1
    },{
      "car_id": 76,
      "cust_id": 203536,
      "finish_position_in_class": 7,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 374674,
      "finish_position_in_class": 9,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 175720,
      "finish_position_in_class": 14,
      "laps_complete": 11
    },{
      "car_id": 76,
      "cust_id": 411930,
      "finish_position_in_class": 8,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 424252,
      "finish_position_in_class": 18,
      "laps_complete": -1
    },{
      "car_id": 76,
      "cust_id": 36256,
      "finish_position_in_class": 11,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 899576,
      "finish_position_in_class": 19,
      "laps_complete": -1
    },{
      "car_id": 76,
      "cust_id": 231698,
      "finish_position_in_class": 12,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 756538,
      "finish_position_in_class": 10,
      "laps_complete": 12
    },{
      "car_id": 76,
      "cust_id": 731520,
      "finish_position_in_class": 15,
      "laps_complete": 9
    },{
      "car_id": 76,
      "cust_id": 1021885,
      "finish_position_in_class": 13,
      "laps_complete": 11
    },{
      "car_id": 76,
      "cust_id": 268321,
      "finish_position_in_class": 16,
      "laps_complete": 4
    }],
  "season_id": 244591934,
  "series_id": 285,
  "subsession_id": 70230671,
  "top_n": 20,
  "track": "Long Beach "
}
`

func TestLive(t *testing.T) {
	t.Run("New standing each lap", func(t *testing.T) {
		filename := "../windows/test-results.json"

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

		s, err := Live(filename, string(ex))
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

		s, err = Live(filename, string(b))
		assert.NoError(t, err)

		fmt.Println(s)

		var prediction2 []live.LiveResults

		err = json.Unmarshal([]byte(s), &prediction2)
		assert.NoError(t, err)
	})
}
