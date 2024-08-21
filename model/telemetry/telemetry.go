// Package telemetry read from Windows shared memory
package telemetry

const (
	IrMaxCars = 64
)

type CarsInfo [IrMaxCars]CarInfo

type CarInfo struct {
	CarClassID          int
	CarID               int
	CarName             string
	CarNumber           string
	CustID              int
	DriverName          string
	IRating             int
	IsPaceCar           bool
	IsSelf              bool
	IsSpectator         bool
	LapsComplete        int
	RacePositionInClass int
}

type TelemetryData struct {
	SeriesID     int
	SessionID    int
	SubsessionID int
	SessionType  string // PRACTICE, QUALIFY, RACE
	Status       string // Connected, Driving
	TrackName    string
	TrackID      int
	DriverCarIdx int
	Cars         CarsInfo
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
