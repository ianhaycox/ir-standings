package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	"github.com/ianhaycox/ir-standings/irsdk"
	"github.com/ianhaycox/ir-standings/model"
)

const (
	defaultRefreshSeconds = 5
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

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
	TrackName   string                        `json:"track_name"`
	CountBestOf int                           `json:"count_best_of"`
	Standings   map[model.CarClassID]Standing `json:"standings"`
}

type Standing struct {
	SoFByCarClass           int              `json:"sof_by_car_class"`
	CarClassID              model.CarClassID `json:"car_class_id"`
	CarClassName            string
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

// Runs as a Go routine reading the Windows shared memory to get session and telemetry data
func irTelemetry(refreshSeconds int) {
	const (
		connectionPauseMilli = 50
		waitForDataMilli     = 100
	)

	var sdk irsdk.SDK

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
			time.Sleep(time.Duration(refreshSeconds) * time.Second)

			continue
		}

		time.Sleep(connectionPauseMilli * time.Millisecond)
		log.Println("iRacing connected")

		for {
			sdk.WaitForData(waitForDataMilli * time.Millisecond)

			if !sdk.IsConnected() {
				log.Println("iRacing disconnected")

				break
			}

			session := sdk.GetSession()

			vars, err := sdk.GetVars()
			if err != nil {
				log.Println("can not get iRacing telemetry vars", err)
				return
			}

			log.Println()
			log.Println("Session")
			log.Printf("%+v", session)

			log.Println()
			log.Println("Vars")

			for i := range vars {
				fmt.Printf("%s\n", vars[i].Name)
			}

			//		log.Printf("%+v", vars)

			varValues, err := sdk.GetVar("CarIdxClassPosition")
			if err != nil {
				log.Println("can not get iRacing telemetry var", err)
				return
			}

			log.Println()
			log.Println("Vars values")

			log.Printf("%+v", varValues.Values)
			log.Println()
		}
	}
}
