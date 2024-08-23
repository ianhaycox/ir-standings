package main

import (
	"embed"
	"net/http"
	"os"
	"strconv"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const (
	defaultRefreshSeconds = 5
	defaultWidth          = 800
	defaultHeight         = 650
	countBestOf           = 10
	seasonYear            = 2024
	seasonQuarter         = 3
	showTopN              = 10
)

var (
	selectedCarClassIDs = []int{84, 83}
)

//go:embed all:frontend/dist
var assets embed.FS

var pointsPerSplit = points.PointsPerSplit{
	//   0   1   2   3   4   5   6   7   8   9  10 11 12 13 14 15 16 17 18 19
	0: {25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	1: {14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	2: {9, 6, 4, 3, 2, 1},
}

func main() {
	refresh := os.Getenv("IR_STANDINGS_REFRESH_SECONDS")

	refreshSeconds, err := strconv.Atoi(refresh)
	if err != nil {
		refreshSeconds = defaultRefreshSeconds
	}

	httpClient := http.DefaultClient
	cookieStore := cookiejar.NewStore(iracing.CookiesFile)
	httpClient.Jar = cookiejar.NewCookieJar(cookieStore)

	cfg := api.NewConfiguration(httpClient, api.UserAgent)
	cfg.AddDefaultHeader("Accept", "application/json")
	cfg.AddDefaultHeader("Content-Type", "application/json")
	client := api.NewHTTPClient(cfg)

	// See auth.go that authenticates separately and saves encrypted credentials in a cookie jar
	ir := iracing.NewIracingService(
		client,
		iracing.NewIracingDataService(
			client, cdn.NewCDNService(api.NewHTTPClient(api.NewConfiguration(http.DefaultClient, ""))),
		),
		api.NewAuthenticationService(),
	)

	// Create an instance of the app structure
	app := NewApp(ir, pointsPerSplit, refreshSeconds, countBestOf, int(iracing.KamelSeriesID), seasonYear, seasonQuarter, showTopN, selectedCarClassIDs)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "iRacing Championship Standings",
		Width:  defaultWidth,
		Height: defaultHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1}, //nolint:mnd // ok
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
