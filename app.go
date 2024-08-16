package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
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

type PredictedStanding struct {
	CustID            int    `json:"cust_id"`
	DriverName        string `json:"driver_name"`
	CarNumber         string `json:"car_number,omitempty"`       // May be blank if not in the session
	CurrentPosition   int    `json:"current_position,omitempty"` // May be first race in the current session
	PredictedPosition int    `json:"predicted_position"`
	CurrentPoints     int    `json:"current_points"`
	PredictedPoints   int    `json:"predicted_points"`
	Change            int    `json:"change"` // +/- change from current position
}

type PredictedStandings struct {
	Items []PredictedStanding `json:"items"`
}

func (a *App) PredictedStandings() PredictedStandings {
	return PredictedStandings{Items: make([]PredictedStanding, 0)}
}
