package telemetry

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/ianhaycox/ir-standings/irsdk"
)

const (
	connectionRetrySecs = 5
	waitForDataMilli    = 100
)

type Data struct {
	sdk  *irsdk.IRSDK
	data *TelemetryData
}

func NewData(sdk *irsdk.IRSDK) Data {
	return Data{
		sdk: sdk,
	}
}

func (d *Data) Telemetry() *TelemetryData {
	d.data = &TelemetryData{
		Status: Waiting,
	}

	d.sdk.WaitForData(waitForDataMilli * time.Millisecond)

	if !d.sdk.IsConnected() {
		log.Println("iRacing disconnected")

		return d.data
	}

	d.data.Status = Connected

	d.updateSession()

	sessionState, err := d.sdk.GetVarValue("SessionState")
	if err != nil {
		log.Println("Error getting SessionState, %w", err)
	} else {
		d.data.SessionState = sessionState.(int)
	}

	d.updateCarInfo()

	return d.data
}

// Updated often
func (d *Data) updateCarInfo() {
	cicp, err := d.sdk.GetVarValues("CarIdxClassPosition")
	if err != nil {
		log.Println("Error getting CarIdxClassPosition, %w", err)
	} else {
		c := cicp.([]int)
		for i := range c {
			if d.data.Cars[i].IsRacing() {
				if c[i] > 1000 {
					fmt.Println("Pos:", i, c[i])
				}
				d.data.Cars[i].RacePositionInClass = c[i]
			} else {
				d.data.Cars[i].RacePositionInClass = 0
			}
		}
	}

	cil, err := d.sdk.GetVarValues("CarIdxLap")
	if err != nil {
		log.Println("Error getting CarIdxLap, %w", err)
	} else {
		c := cil.([]int)
		for i := range c {
			if d.data.Cars[i].IsRacing() {
				// I think byte4ToInt() should handle signed ints - bodge it here.
				if c[i] > math.MaxInt32 {
					d.data.Cars[i].LapsComplete = -1
				} else {
					d.data.Cars[i].LapsComplete = c[i]
				}
			} else {
				d.data.Cars[i].LapsComplete = 0
			}
		}
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
