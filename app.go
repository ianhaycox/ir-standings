package main

import (
	"context"
	"encoding/json"
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
	"github.com/ianhaycox/ir-standings/test/devmode"
)

const (
	connectionRetrySecs = 5
	waitForDataMilli    = 100
)

// App struct
type App struct {
	ctx context.Context
	mtx sync.Mutex

	pastResults   []results.Result        // Previous weeks results for this season from the iRacing API
	telemetryData telemetry.TelemetryData // Shared memory updated every `refreshSeconds`
	prediction    *predictor.Predictor

	irAPI               iracing.IracingService // iRacing API
	pointsPerSplit      points.PointsPerSplit  // Points structure
	refreshSeconds      int                    // How often to read telemetry
	countBestOf         int                    // Count best of n races in season
	seriesID            int                    // iRacing series ID
	seasonYear          int                    // E.g. 2024, 2025
	seasonQuarter       int                    // E.g. 1,2,3
	triedPastResults    bool                   // Guard against call API too much
	showTopN            int                    // Display top n standings
	selectedCarClassIDs []int                  // Display car class IDs in table in order
}

// NewApp creates a new App application struct
func NewApp(irAPI iracing.IracingService, pointsPerSplit points.PointsPerSplit, refreshSeconds, countBestOf, seriesID,
	seasonYear, seasonQuarter, showTopN int, selectedCarClassIDs []int) *App {
	return &App{
		irAPI:               irAPI,
		refreshSeconds:      refreshSeconds,
		seriesID:            seriesID,
		seasonYear:          seasonYear,
		seasonQuarter:       seasonQuarter,
		telemetryData:       telemetry.TelemetryData{Status: "Disconnected"},
		pointsPerSplit:      pointsPerSplit,
		triedPastResults:    false,
		countBestOf:         countBestOf,
		showTopN:            showTopN,
		selectedCarClassIDs: selectedCarClassIDs,
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

func (a *App) ShowTopN() int {
	return a.showTopN
}

func (a *App) SelectedCarClassIDs() []int {
	return a.selectedCarClassIDs
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
	if devmode.IsDevMode() {
		const fakeDelay = 5

		filename := "./model/fixtures/2024-2-285-results-redacted.json"
		log.Println("Using results fixture in dev mode:", filename)

		a.triedPastResults = true

		buf, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		a.pastResults = make([]results.Result, 0)

		err = json.Unmarshal(buf, &a.pastResults)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(fakeDelay * time.Second)
	}

	if !a.triedPastResults {
		searchSeriesResults, err := a.irAPI.SearchSeriesResults(a.ctx, a.seasonYear, a.seasonQuarter, a.seriesID)
		if err != nil {
			log.Println("can not get series results:", err)
		}

		a.pastResults, err = a.irAPI.SeasonBroadcastResults(a.ctx, searchSeriesResults)
		if err != nil {
			log.Println("can not get series results:", err)
		}

		log.Println("Got results from iRacing API")

		return true
	}

	log.Println("Skipped iRacing results API")

	return false
}

func (a *App) LatestStandings() live.PredictedStandings {
	log.Println("LatestStandings")

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.prediction == nil {
		a.prediction = predictor.NewPredictor(a.pointsPerSplit, a.pastResults, &a.telemetryData, a.countBestOf)
	}

	return a.prediction.Live()
}

// Runs as a Go routine reading the Windows shared memory to get session and telemetry data
func (a *App) irTelemetry(refreshSeconds int) {
	var (
		sdk irsdk.SDK
	)

	log.Println("Starting iRacing telemetry...")

	if runtime.GOOS == "windows" {
		sdk = irsdk.Init(nil)
	} else {
		reader, err := os.Open("/tmp/test.ibt")
		if err != nil {
			log.Fatal(err)
		}

		sdk = irsdk.Init(reader)
	}

	defer sdk.Close()
	defer a.mtx.Unlock()

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

			a.mtx.Lock()

			if sdk.SessionChanged() {
				session := sdk.GetSession()

				sessionNum, err := sdk.GetVar("SessionNum")
				if err != nil {
					log.Println("can not get iRacing SessionNum var", err)
					a.mtx.Unlock()

					continue
				}

				if session.SessionInfo.Sessions[sessionNum.Value.(int)].SessionName != "RACE" {
					a.mtx.Unlock()
					continue
				}

				a.updateSession(sdk, &session)
			}

			vars, err := sdk.GetVars()
			if err != nil {
				log.Println("can not get iRacing telemetry vars", err)
				a.mtx.Unlock()

				continue
			}

			a.updateCarInfo(vars)

			a.mtx.Unlock()
		}
	}
}

// Updated often
func (a *App) updateCarInfo(vars map[string]irsdk.Variable) {
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
		a.telemetryData.Cars[carIdx].CarClassName = session.DriverInfo.Drivers[i].CarClassShortName
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
