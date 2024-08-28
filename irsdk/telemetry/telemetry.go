// Package telemetry read from Windows shared memory
package telemetry

const (
	IrMaxCars    = 64
	Connected    = "Connected"
	Disconnected = "Disconnected"
	Waiting      = "Waiting for iRacing connection..."
	Problem      = "Telemetry unavailable"
)

type CarsInfo [IrMaxCars]CarInfo

type CarInfo struct {
	CarClassID          int    `json:"car_class_id"`
	CarID               int    `json:"car_id"`
	CarNumber           string `json:"car_number"`
	CustID              int    `json:"cust_id"`
	DriverName          string `json:"driver_name"`
	IRating             int    `json:"irating"`
	IsPaceCar           bool   `json:"is_pace_car"`
	IsSelf              bool   `json:"is_self"`
	IsSpectator         bool   `json:"is_spectator"`
	LapsComplete        int    `json:"laps_complete"`
	RacePositionInClass int    `json:"race_position_in_class"`
}

type TelemetryData struct {
	SeriesID       int      `json:"series_id"`
	SessionID      int      `json:"session_id"`
	SubsessionID   int      `json:"subsession_id"`
	SessionType    string   `json:"session_type"`  // PRACTICE, QUALIFY, RACE
	SessionState   int      `json:"session_state"` // Warmup, Racing, Cooldown etc.
	Status         string   `json:"status"`        // Connected, Driving
	TrackName      string   `json:"track_name"`
	TrackID        int      `json:"track_id"`
	DriverCarIdx   int      `json:"driver_car_idx"`
	SelfCarClassID int      `json:"self_car_class_id"`
	Cars           CarsInfo `json:"cars,omitempty"`
}

func (td *TelemetryData) SofByCarClass() map[int]int {
	total := make(map[int]int)
	count := make(map[int]int)

	for i := range td.Cars {
		if !td.Cars[i].IsRacing() {
			continue
		}

		count[td.Cars[i].CarClassID]++
		total[td.Cars[i].CarClassID] += td.Cars[i].IRating
	}

	sof := make(map[int]int)

	for carClassID := range total {
		sof[carClassID] = total[carClassID] / count[carClassID]
	}

	return sof
}

func (td *TelemetryData) LeaderLapsComplete(carClassID int) int {
	leaderLapsComplete := 0

	for i := range td.Cars {
		if td.Cars[i].CarClassID != carClassID {
			continue
		}

		if !td.Cars[i].IsRacing() {
			continue
		}

		if td.Cars[i].LapsComplete > leaderLapsComplete {
			leaderLapsComplete = td.Cars[i].LapsComplete
		}
	}

	return leaderLapsComplete
}

func (ci *CarInfo) IsRacing() bool {
	return !(ci.IsPaceCar || ci.IsSpectator || ci.DriverName == "")
}
