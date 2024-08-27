package telemetry

import (
	"log"
	"strings"
	"time"

	"github.com/ianhaycox/ir-standings/irsdk"
)

const (
	connectionRetrySecs = 5
	waitForDataMilli    = 100
)

type Data struct {
	sdk  irsdk.SDK
	data *TelemetryData
}

func NewData(sdk irsdk.SDK) Data {
	return Data{
		sdk: sdk,
	}
}

func (d *Data) Telemetry() *TelemetryData {
	d.data = &TelemetryData{
		Status: Waiting,
	}

	if !d.sdk.IsConnected() {
		log.Println("Waiting for iRacing connection...")

		return d.data
	}

	log.Println("iRacing connected")

	d.data.Status = Connected

	d.sdk.WaitForData(waitForDataMilli * time.Millisecond)

	if !d.sdk.IsConnected() {
		log.Println("iRacing disconnected")

		d.data.Status = Disconnected

		return d.data
	}

	d.updateSession()

	vars, err := d.sdk.GetVars()
	if err != nil {
		log.Println("can not get iRacing telemetry vars", err)

		d.data.Status = Problem

		return d.data
	}

	d.updateCarInfo(vars)

	log.Println("Type", d.data.SessionType)

	return d.data
}

// Updated often
func (d *Data) updateCarInfo(vars map[string]irsdk.Variable) {
	var (
		carIdxPositionsByClass []int
		carIdxLaps             []int
	)

	carIdxPositionByClass := vars["CarIdxClassPosition"].Values
	if carIdxPositionByClass != nil {
		carIdxPositionsByClass = carIdxPositionByClass.([]int)
	}

	carIdxLap := vars["CarIdxLap"].Values
	if carIdxLap != nil {
		carIdxLaps = carIdxLap.([]int)
	}

	for i := range carIdxPositionsByClass {
		d.data.Cars[i].RacePositionInClass = carIdxPositionsByClass[i]
		d.data.Cars[i].LapsComplete = carIdxLaps[i]
	}
}

var (
	cleanCRLF = strings.NewReplacer("\r", "", "\n", "")
)

// Updated rarely
func (d *Data) updateSession() {
	session := d.sdk.GetSession()

	d.data.SeriesID = session.WeekendInfo.SeriesID
	d.data.SessionID = session.WeekendInfo.SessionID
	d.data.SubsessionID = session.WeekendInfo.SubSessionID

	sessionNum, err := d.sdk.GetVar("SessionNum")
	if err == nil {
		d.data.SessionType = session.SessionInfo.Sessions[sessionNum.Value.(int)].SessionName
	}

	d.data.TrackName = session.WeekendInfo.TrackDisplayName + " " + session.WeekendInfo.TrackConfigName
	d.data.TrackID = session.WeekendInfo.TrackID
	d.data.DriverCarIdx = session.DriverInfo.DriverCarIdx

	for i := range session.DriverInfo.Drivers {
		carIdx := session.DriverInfo.Drivers[i].CarIdx

		d.data.Cars[carIdx].CarClassID = session.DriverInfo.Drivers[i].CarClassID
		d.data.Cars[carIdx].CarID = session.DriverInfo.Drivers[i].CarID
		d.data.Cars[carIdx].CarNumber = session.DriverInfo.Drivers[i].CarNumber
		d.data.Cars[carIdx].CustID = session.DriverInfo.Drivers[i].UserID
		d.data.Cars[carIdx].DriverName = cleanCRLF.Replace(session.DriverInfo.Drivers[i].UserName)
		d.data.Cars[carIdx].IRating = session.DriverInfo.Drivers[i].IRating
		d.data.Cars[carIdx].IsPaceCar = session.DriverInfo.Drivers[i].CarIsPaceCar == 1
		d.data.Cars[carIdx].IsSelf = carIdx == d.data.DriverCarIdx
		d.data.Cars[carIdx].IsSpectator = session.DriverInfo.Drivers[i].IsSpectator == 1

		if d.data.Cars[carIdx].IsSelf {
			d.data.SelfCarClassID = session.DriverInfo.Drivers[i].CarClassID
		}
	}
}
