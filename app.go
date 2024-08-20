package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	"github.com/ianhaycox/ir-standings/irsdk"
	"github.com/ianhaycox/ir-standings/irsdk/iryaml"
	"github.com/ianhaycox/ir-standings/model"
)

const (
	defaultRefreshSeconds = 5
)

// App struct
type App struct {
	ctx context.Context
	sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

var (
	cleanCRLF = strings.NewReplacer("\r", "", "\n", "")
)

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	refresh := os.Getenv("IR_STANDINGS_REFRESH_SECONDS")

	refreshSeconds, err := strconv.Atoi(refresh)
	if err != nil {
		refreshSeconds = defaultRefreshSeconds
	}

	go irTelemetry(refreshSeconds)
}

func (a *App) Login(username string, password string) bool {
	if username == "test" {
		return true
	}

	ctx := context.Background()
	httpClient := http.DefaultClient
	cookieStore := cookiejar.NewStore(iracing.CookiesFile)
	jar := cookiejar.NewCookieJar(cookieStore)

	if !jar.HasExpired("members-ng.iracing.com", "authtoken_members") {
		log.Println("Using cached login token")

		return true
	}

	httpClient.Jar = jar
	cfg := api.NewConfiguration(httpClient, api.UserAgent)
	cfg.AddDefaultHeader("Accept", "application/json")
	cfg.AddDefaultHeader("Content-Type", "application/json")

	auth := api.NewAuthenticationService(username, password)
	client := api.NewHTTPClient(cfg)

	ir := iracing.NewIracingService(client, nil, auth)

	err := ir.Authenticate(ctx)
	if err != nil {
		log.Println("Can not login: ", err)

		return false
	}

	return true
}

// {CarClassID: 84, ShortName: "GTP", Name: "Nissan GTP ZX-T", CarsInClass: []results.CarsInClass{{CarID: 77}}},
// {CarClassID: 83, ShortName: "GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},

type PredictedStandings struct {
	Status      string                        `json:"status"` // iRacing connection status/session, Race, Qualifying,...
	TrackName   string                        `json:"track_name"`
	CountBestOf int                           `json:"count_best_of"`
	Standings   map[model.CarClassID]Standing `json:"standings"`
}

type Standing struct {
	SoFByCarClass           int                 `json:"sof_by_car_class"`
	CarClassID              model.CarClassID    `json:"car_class_id"`
	CarClassName            string              `json:"car_class_name"`
	ClassLeaderLapsComplete model.LapsComplete  `json:"class_leader_laps_complete"`
	Items                   []PredictedStanding `json:"items"`
}

type PredictedStanding struct {
	CustID            model.CustID                `json:"cust_id"`                    // Key for React
	DriverName        string                      `json:"driver_name"`                // Driver
	CarNumber         string                      `json:"car_number,omitempty"`       // May be blank if not in the session
	CurrentPosition   model.FinishPositionInClass `json:"current_position,omitempty"` // May be first race in the current session
	PredictedPosition model.FinishPositionInClass `json:"predicted_position"`         // Championship position as is
	CurrentPoints     model.Point                 `json:"current_points"`             // Championship position before race
	PredictedPoints   model.Point                 `json:"predicted_points"`           // Championship points as is
	Change            int                         `json:"change"`                     // +/- change from current position
	CarNames          []string                    `json:"car_names"`                  // Cars driven in this class
}

//nolint:mnd,lll,dupl // ok
func (a *App) PredictedStandings() PredictedStandings {
	i83 := []PredictedStanding{
		{CustID: 1, DriverName: "Audi 1", CarNumber: "1", CurrentPosition: 1, PredictedPosition: 2, CurrentPoints: 100, PredictedPoints: 100, Change: -1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 2, DriverName: "Audi 2", CarNumber: "2", CurrentPosition: 2, PredictedPosition: 1, CurrentPoints: 95, PredictedPoints: 110, Change: 1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 3, DriverName: "Audi 3", CarNumber: "3", CurrentPosition: 3, PredictedPosition: 3, CurrentPoints: 90, PredictedPoints: 90, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 4, DriverName: "Audi 4", CarNumber: "4", CurrentPosition: 4, PredictedPosition: 4, CurrentPoints: 85, PredictedPoints: 85, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 5, DriverName: "Audi 5", CarNumber: "5", CurrentPosition: 5, PredictedPosition: 5, CurrentPoints: 80, PredictedPoints: 80, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 6, DriverName: "Audi 6", CarNumber: "6", CurrentPosition: 6, PredictedPosition: 7, CurrentPoints: 75, PredictedPoints: 75, Change: -1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 7, DriverName: "Audi 7", CarNumber: "7", CurrentPosition: 7, PredictedPosition: 6, CurrentPoints: 70, PredictedPoints: 78, Change: 1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 8, DriverName: "Audi 8", CarNumber: "8", CurrentPosition: 8, PredictedPosition: 8, CurrentPoints: 65, PredictedPoints: 65, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 9, DriverName: "Audi 9", CarNumber: "9", CurrentPosition: 9, PredictedPosition: 9, CurrentPoints: 60, PredictedPoints: 60, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 10, DriverName: "Audi 10", CarNumber: "10", CurrentPosition: 10, PredictedPosition: 10, CurrentPoints: 55, PredictedPoints: 55, Change: 0, CarNames: []string{"Audi 90 GTO"}},
	}

	i84 := []PredictedStanding{
		{CustID: 11, DriverName: "Nissan 1", CarNumber: "1", CurrentPosition: 1, PredictedPosition: 2, CurrentPoints: 100, PredictedPoints: 100, Change: -1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 12, DriverName: "Nissan 2", CarNumber: "2", CurrentPosition: 2, PredictedPosition: 1, CurrentPoints: 95, PredictedPoints: 110, Change: 1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 13, DriverName: "Nissan 3", CarNumber: "3", CurrentPosition: 3, PredictedPosition: 3, CurrentPoints: 90, PredictedPoints: 90, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 14, DriverName: "Nissan 4", CarNumber: "4", CurrentPosition: 4, PredictedPosition: 4, CurrentPoints: 85, PredictedPoints: 85, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 15, DriverName: "Nissan 5", CarNumber: "5", CurrentPosition: 5, PredictedPosition: 5, CurrentPoints: 80, PredictedPoints: 80, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 16, DriverName: "Nissan 6", CarNumber: "6", CurrentPosition: 6, PredictedPosition: 7, CurrentPoints: 75, PredictedPoints: 75, Change: -1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 17, DriverName: "Nissan 7", CarNumber: "7", CurrentPosition: 7, PredictedPosition: 6, CurrentPoints: 70, PredictedPoints: 78, Change: 1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 18, DriverName: "Nissan 8", CarNumber: "8", CurrentPosition: 8, PredictedPosition: 8, CurrentPoints: 65, PredictedPoints: 65, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 19, DriverName: "Nissan 9", CarNumber: "9", CurrentPosition: 9, PredictedPosition: 9, CurrentPoints: 60, PredictedPoints: 60, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 20, DriverName: "Nissan 10", CarNumber: "10", CurrentPosition: 10, PredictedPosition: 10, CurrentPoints: 55, PredictedPoints: 55, Change: 0, CarNames: []string{"Nissan ZX-T"}},
	}

	s83 := Standing{
		SoFByCarClass:           2087,
		CarClassID:              83,
		CarClassName:            "GTO",
		ClassLeaderLapsComplete: 10,
		Items:                   i83,
	}

	s84 := Standing{
		SoFByCarClass:           3025,
		CarClassID:              84,
		CarClassName:            "GTP",
		ClassLeaderLapsComplete: 12,
		Items:                   i84,
	}

	ps := PredictedStandings{
		TrackName:   "Motegi Resort",
		CountBestOf: 10,
		Standings:   map[model.CarClassID]Standing{84: s84, 83: s83},
	}

	return ps
}

func (a *App) SendTelemetry(telemetryJSON string) {
	log.Println("Telemetry:", telemetryJSON)
}

func (a *App) SendSessionInfo(sessionJSON string) {
	log.Println("SessionInfo:", sessionJSON)
}

const irMaxCars = 64

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
	LapComplete         int
	RacePositionInClass int
}

type TelemetryData struct {
	SeriesID     int
	SessionID    int
	SubsessionID int
	SessionType  string // PRACTICE, QUALIFY, RACE
	Status       string // Connected, Driving
	TrackName    string
	DriverCarIdx int
	Cars         [irMaxCars]CarInfo
}

// Runs as a Go routine reading the Windows shared memory to get session and telemetry data
func irTelemetry(refreshSeconds int) {
	const (
		connectionRetrySecs = 5
		waitForDataMilli    = 100
	)

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

	for {
		if !sdk.IsConnected() {
			log.Println("Waiting for iRacing connection...")
			time.Sleep(time.Duration(connectionRetrySecs) * time.Second)

			continue
		}

		log.Println("iRacing connected")

		telemetryData := TelemetryData{Status: "Connected"}

		for {
			time.Sleep(time.Duration(refreshSeconds) * time.Second)

			sdk.WaitForData(waitForDataMilli * time.Millisecond)

			if !sdk.IsConnected() {
				log.Println("iRacing disconnected")

				break
			}

			if sdk.SessionChanged() {
				session := sdk.GetSession()

				updateSession(sdk, &session, &telemetryData)
			}

			vars, err := sdk.GetVars()
			if err != nil {
				log.Println("can not get iRacing telemetry vars", err)
				continue
			}

			updateTelemetry(vars, &telemetryData)

			log.Println(telemetryData)
		}
	}
}

// Updated often
func updateTelemetry(vars map[string]irsdk.Variable, telemetryData *TelemetryData) {
	var (
		carIdxPositionsByClass [irMaxCars]int
		carIdxLaps             [irMaxCars]int
	)

	carIdxPositionByClass := vars["CarIdxClassPosition"].Values
	if carIdxPositionByClass != nil {
		carIdxPositionsByClass = carIdxPositionByClass.([irMaxCars]int)
	}

	carIdxLap := vars["CarIdxLap"].Values
	if carIdxLap != nil {
		carIdxLaps = carIdxLap.([irMaxCars]int)
	}

	for i := 0; i < irMaxCars; i++ {
		telemetryData.Cars[i].RacePositionInClass = carIdxPositionsByClass[i]
		telemetryData.Cars[i].LapComplete = carIdxLaps[i]
	}
}

// Updated rarely
func updateSession(sdk irsdk.SDK, session *iryaml.IRSession, telemetryData *TelemetryData) {
	telemetryData.SeriesID = session.WeekendInfo.SeriesID
	telemetryData.SessionID = session.WeekendInfo.SessionID
	telemetryData.SubsessionID = session.WeekendInfo.SubSessionID

	sessionNum, err := sdk.GetVar("SessionNum")
	if err == nil {
		telemetryData.SessionType = session.SessionInfo.Sessions[sessionNum.Value.(int)].SessionName
	}

	telemetryData.TrackName = session.WeekendInfo.TrackDisplayName + " " + session.WeekendInfo.TrackConfigName
	telemetryData.DriverCarIdx = session.DriverInfo.DriverCarIdx

	for i := range session.DriverInfo.Drivers {
		carIdx := session.DriverInfo.Drivers[i].CarIdx

		telemetryData.Cars[carIdx].CarClassID = session.DriverInfo.Drivers[i].CarClassID
		telemetryData.Cars[carIdx].CarID = session.DriverInfo.Drivers[i].CarID
		telemetryData.Cars[carIdx].CarName = session.DriverInfo.Drivers[i].CarScreenName
		telemetryData.Cars[carIdx].CarNumber = session.DriverInfo.Drivers[i].CarNumber
		telemetryData.Cars[carIdx].CustID = session.DriverInfo.Drivers[i].UserID
		telemetryData.Cars[carIdx].DriverName = cleanCRLF.Replace(session.DriverInfo.Drivers[i].UserName)
		telemetryData.Cars[carIdx].IRating = session.DriverInfo.Drivers[i].IRating
		telemetryData.Cars[carIdx].IsPaceCar = session.DriverInfo.Drivers[i].CarIsPaceCar == 1
		telemetryData.Cars[carIdx].IsSelf = carIdx == telemetryData.DriverCarIdx
		telemetryData.Cars[carIdx].IsSpectator = session.DriverInfo.Drivers[i].IsSpectator == 1
	}
}

func sofByCarClass(td *TelemetryData) map[int]int {
	total := make(map[int]int)
	count := make(map[int]int)

	for i := range td.Cars {
		if td.Cars[i].IsPaceCar || td.Cars[i].IsSpectator || td.Cars[i].DriverName == "" {
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
