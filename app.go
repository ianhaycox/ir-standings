package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ianhaycox/ir-standings/connectors/iracing"
	"github.com/ianhaycox/ir-standings/irsdk"
	"github.com/ianhaycox/ir-standings/irsdk/iryaml"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/ianhaycox/ir-standings/model/telemetry"
	"github.com/ianhaycox/ir-standings/predictor"
)

const (
	connectionRetrySecs = 5
	waitForDataMilli    = 100
)

// App struct
type App struct {
	ctx context.Context
	sync.Mutex

	irAPI            iracing.IracingService  // iRacing API
	pointsPerSplit   points.PointsPerSplit   // Points structure
	refreshSeconds   int                     // How often to read telemetry
	countBestOf      int                     // Count best of n races in season
	seriesID         int                     // iRacing series ID
	seasonYear       int                     // E.g. 2024, 2025
	seasonQuarter    int                     // E.g. 1,2,3
	triedPastResults bool                    // Guard against call API too much
	pastResults      []results.Result        // Previous weeks results for this season from the iRacing API
	telemetryData    telemetry.TelemetryData // Shared memory updated every `refreshSeconds`
}

// NewApp creates a new App application struct
func NewApp(irAPI iracing.IracingService, pointsPerSplit points.PointsPerSplit, refreshSeconds, countBestOf, seriesID, seasonYear, seasonQuarter int) *App {
	return &App{
		irAPI:            irAPI,
		refreshSeconds:   refreshSeconds,
		seriesID:         seriesID,
		seasonYear:       seasonYear,
		seasonQuarter:    seasonQuarter,
		telemetryData:    telemetry.TelemetryData{Status: "Disconnected"},
		pointsPerSplit:   pointsPerSplit,
		triedPastResults: true,
		countBestOf:      countBestOf,
	}
}

var (
	cleanCRLF = strings.NewReplacer("\r", "", "\n", "")
)

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	go a.irTelemetry(a.refreshSeconds)
}

func (a *App) Login(email string, password string) bool {
	if email == "test" {
		return true
	}

	err := a.irAPI.Authenticate(a.ctx, email, password)
	if err != nil {
		log.Println("Can not login: ", err)

		return false
	}

	return true
}

func (a *App) PastResults() bool {
	log.Println("PastResults")

	if !a.triedPastResults {
		searchSeriesResults, err := a.irAPI.SearchSeriesResults(a.ctx, a.seasonYear, a.seasonQuarter, a.seriesID)
		if err != nil {
			log.Println("can not get series results:", err)
		}

		a.pastResults, err = a.irAPI.SeasonBroadcastResults(a.ctx, searchSeriesResults)
		if err != nil {
			log.Println("can not get series results:", err)
		}

		log.Println("Done")

		return true
	}

	log.Println("Skipped")

	return false
}

func (a *App) LatestStandings() live.PredictedStandings {
	log.Println("LatestStandings")

	ps := predictor.NewPredictor(a.pointsPerSplit, a.pastResults, a.countBestOf)
	s := ps.Live(&a.telemetryData)

	log.Println("Returning", s)

	return s
}

// Runs as a Go routine reading the Windows shared memory to get session and telemetry data
func (a *App) irTelemetry(refreshSeconds int) {
	a.Lock()
	defer a.Unlock()

	var (
		sdk irsdk.SDK
	)

	log.Println("Starting iRacing telemetry...")

	if runtime.GOOS == "windows" {
		sdk = irsdk.Init(nil)
	} else {
		reader, err := os.Open("/tmp/test.ibt")
		if err != nil {
			log.Fatal(err) //nolint:gocritic // for dev only
		}

		sdk = irsdk.Init(reader)
	}

	defer sdk.Close()

	for {
		if !sdk.IsConnected() {
			log.Println("Waiting for iRacing connection...")
			time.Sleep(time.Duration(connectionRetrySecs) * time.Second)

			continue
		}

		log.Println("iRacing connected")

		a.telemetryData.Status = "Connected"

		for {
			time.Sleep(time.Duration(refreshSeconds) * time.Second)

			sdk.WaitForData(waitForDataMilli * time.Millisecond)

			if !sdk.IsConnected() {
				log.Println("iRacing disconnected")

				break
			}

			if sdk.SessionChanged() {
				session := sdk.GetSession()

				a.updateSession(sdk, &session)
			}

			vars, err := sdk.GetVars()
			if err != nil {
				log.Println("can not get iRacing telemetry vars", err)
				continue
			}

			a.updateTelemetry(vars)

			log.Println(a.telemetryData.TrackName)
		}
	}
}

// Updated often
func (a *App) updateTelemetry(vars map[string]irsdk.Variable) {
	var (
		carIdxPositionsByClass [telemetry.IrMaxCars]int
		carIdxLaps             [telemetry.IrMaxCars]int
	)

	carIdxPositionByClass := vars["CarIdxClassPosition"].Values
	if carIdxPositionByClass != nil {
		carIdxPositionsByClass = carIdxPositionByClass.([telemetry.IrMaxCars]int)
	}

	carIdxLap := vars["CarIdxLap"].Values
	if carIdxLap != nil {
		carIdxLaps = carIdxLap.([telemetry.IrMaxCars]int)
	}

	for i := 0; i < telemetry.IrMaxCars; i++ {
		a.telemetryData.Cars[i].RacePositionInClass = carIdxPositionsByClass[i]
		a.telemetryData.Cars[i].LapsComplete = carIdxLaps[i]
	}
}

// Updated rarely
func (a *App) updateSession(sdk irsdk.SDK, session *iryaml.IRSession) {
	a.telemetryData.SeriesID = session.WeekendInfo.SeriesID
	a.telemetryData.SessionID = session.WeekendInfo.SessionID
	a.telemetryData.SubsessionID = session.WeekendInfo.SubSessionID

	sessionNum, err := sdk.GetVar("SessionNum")
	if err == nil {
		a.telemetryData.SessionType = session.SessionInfo.Sessions[sessionNum.Value.(int)].SessionName
	}

	a.telemetryData.TrackName = session.WeekendInfo.TrackDisplayName + " " + session.WeekendInfo.TrackConfigName
	a.telemetryData.TrackID = session.WeekendInfo.TrackID
	a.telemetryData.DriverCarIdx = session.DriverInfo.DriverCarIdx

	for i := range session.DriverInfo.Drivers {
		carIdx := session.DriverInfo.Drivers[i].CarIdx

		a.telemetryData.Cars[carIdx].CarClassID = session.DriverInfo.Drivers[i].CarClassID
		a.telemetryData.Cars[carIdx].CarID = session.DriverInfo.Drivers[i].CarID
		a.telemetryData.Cars[carIdx].CarName = session.DriverInfo.Drivers[i].CarScreenName
		a.telemetryData.Cars[carIdx].CarNumber = session.DriverInfo.Drivers[i].CarNumber
		a.telemetryData.Cars[carIdx].CustID = session.DriverInfo.Drivers[i].UserID
		a.telemetryData.Cars[carIdx].DriverName = cleanCRLF.Replace(session.DriverInfo.Drivers[i].UserName)
		a.telemetryData.Cars[carIdx].IRating = session.DriverInfo.Drivers[i].IRating
		a.telemetryData.Cars[carIdx].IsPaceCar = session.DriverInfo.Drivers[i].CarIsPaceCar == 1
		a.telemetryData.Cars[carIdx].IsSelf = carIdx == a.telemetryData.DriverCarIdx
		a.telemetryData.Cars[carIdx].IsSpectator = session.DriverInfo.Drivers[i].IsSpectator == 1
	}
}
