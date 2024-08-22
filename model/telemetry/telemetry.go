// Package telemetry read from Windows shared memory
package telemetry

import (
	"github.com/ianhaycox/ir-standings/model/data/results"
)

const (
	IrMaxCars = 64
)

type CarsInfo [IrMaxCars]CarInfo

type CarInfo struct {
	CarClassID          int    `json:"car_class_id"`
	CarClassName        string `json:"car_class_name"`
	CarID               int    `json:"car_id"`
	CarName             string `json:"car_name"`
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
	SeriesID     int      `json:"series_id"`
	SessionID    int      `json:"session_id"`
	SubsessionID int      `json:"subsession_id"`
	SessionType  string   `json:"session_type"` // PRACTICE, QUALIFY, RACE
	Status       string   `json:"status"`       // Connected, Driving
	TrackName    string   `json:"track_name"`
	TrackID      int      `json:"track_id"`
	DriverCarIdx int      `json:"driver_car_idx"`
	Cars         CarsInfo `json:"cars"`
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

// CarClasses build cars in session - TODO for a Race should only do once.
func (td *TelemetryData) CarClasses() []results.CarClasses {
	carClasses := make([]results.CarClasses, 0)

	cc := make(map[int]results.CarClasses) // By CarClassID
	cic := make(map[int]int)               // By CarID

	for i := range td.Cars {
		if !td.Cars[i].IsRacing() {
			continue
		}

		if _, ok := cc[td.Cars[i].CarClassID]; !ok {
			cc[td.Cars[i].CarClassID] = results.CarClasses{
				CarClassID:  td.Cars[i].CarClassID,
				ShortName:   td.Cars[i].CarClassName,
				Name:        td.Cars[i].CarClassName,
				CarsInClass: make([]results.CarsInClass, 0),
			}
		}

		if _, ok := cic[td.Cars[i].CarID]; !ok {
			cic[td.Cars[i].CarID] = td.Cars[i].CarClassID
		}
	}

	for carClassID, carClass := range cc {
		for carID, memberCarClassID := range cic {
			if memberCarClassID == carClassID {
				carClass.CarsInClass = append(carClass.CarsInClass, results.CarsInClass{CarID: carID})
			}
		}

		carClasses = append(carClasses, carClass)
	}

	// b, _ := json.MarshalIndent(carClasses, "", "  ")
	// fmt.Println(string(b))

	return carClasses
}
