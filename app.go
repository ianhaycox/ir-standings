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

	"github.com/ianhaycox/ir-standings/arch"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	"github.com/ianhaycox/ir-standings/irsdk"
	"github.com/ianhaycox/ir-standings/irsdk/iryaml"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/ianhaycox/ir-standings/model/telemetry"
	"github.com/ianhaycox/ir-standings/predictor"
	"golang.org/x/exp/maps"
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
	fakeLogin     bool

	irAPI            iracing.IracingService // iRacing API
	pointsPerSplit   points.PointsPerSplit  // Points structure
	refreshSeconds   int                    // How often to read telemetry
	countBestOf      int                    // Count best of n races in season
	seriesID         int                    // iRacing series ID
	seasonYear       int                    // E.g. 2024, 2025
	seasonQuarter    int                    // E.g. 1,2,3
	triedPastResults bool                   // Guard against call API too much
	showTopN         int                    // Display top n standings
}

// NewApp creates a new App application struct
func NewApp(irAPI iracing.IracingService, pointsPerSplit points.PointsPerSplit, refreshSeconds, countBestOf, seriesID, showTopN int) *App {
	return &App{
		irAPI:            irAPI,
		refreshSeconds:   refreshSeconds,
		seriesID:         seriesID,
		telemetryData:    telemetry.TelemetryData{Status: "Disconnected"},
		pointsPerSplit:   pointsPerSplit,
		triedPastResults: false,
		countBestOf:      countBestOf,
		showTopN:         showTopN,
	}
}

var (
	cleanCRLF = strings.NewReplacer("\r", "", "\n", "")
)

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if runtime.GOOS == "windows" {
		arch.WindowOptions()
	}

	go a.irTelemetry(a.refreshSeconds)
}

func (a *App) ShowTopN() int {
	return a.showTopN
}

func (a *App) Login(email string, password string) bool {
	if email == "test" {
		a.fakeLogin = true

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
	if a.fakeLogin {
		const fakeDelay = 3

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
		seasons, err := a.irAPI.Seasons(a.ctx)
		if err != nil {
			log.Println("can not get series results:", err)
		}

		for i := range seasons {
			if seasons[i].SeriesID == a.seriesID {
				a.seasonYear = seasons[i].SeasonYear
				a.seasonQuarter = seasons[i].SeasonQuarter

				break
			}
		}

		log.Println("Getting results for :", a.seasonYear, a.seasonQuarter)

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

	a.prediction = predictor.NewPredictor(a.pointsPerSplit, a.pastResults, &a.telemetryData, a.countBestOf)

	return a.prediction.Live()
}

// Runs as a Go routine reading the Windows shared memory to get session and telemetry data
func (a *App) irTelemetry(refreshSeconds int) {
	var (
		sdk irsdk.SDK
	)

	log.Println("Starting iRacing telemetry...")

	if runtime.GOOS == "windows" {
		log.Println("Init irSDK Windows")
		sdk = irsdk.Init(nil)
	} else {
		reader, err := os.Open("/tmp/test.ibt")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Init irSDK Linux(other)")

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

			//if sdk.SessionChanged() {
			session := sdk.GetSession()
			a.updateSession(sdk, &session)
			//}

			vars, err := sdk.GetVars()
			if err != nil {
				log.Println("can not get iRacing telemetry vars", err)
				a.mtx.Unlock()

				continue
			}

			a.updateCarInfo(vars)

			log.Println("Type", a.telemetryData.SessionType)
			log.Println("Track", a.telemetryData.TrackName)
			b, _ := json.MarshalIndent(a.telemetryData.Cars[0:2], "", "  ")
			log.Println(string(b))

			a.mtx.Unlock()
		}
	}
}

// Updated often
func (a *App) updateCarInfo(vars map[string]irsdk.Variable) {
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

	carClassIDs := make(map[int]bool)

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

		if a.telemetryData.Cars[carIdx].IsSelf {
			a.telemetryData.SelfCarClassID = session.DriverInfo.Drivers[i].CarClassID
		}

		if a.telemetryData.Cars[carIdx].IsRacing() {
			carClassIDs[session.DriverInfo.Drivers[i].CarClassID] = true
		}
	}

	a.telemetryData.CarClassIDs = maps.Keys(carClassIDs)

	if a.telemetryData.SelfCarClassID == 0 && len(a.telemetryData.CarClassIDs) > 0 {
		a.telemetryData.SelfCarClassID = a.telemetryData.CarClassIDs[0]
	}
}
