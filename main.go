package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	"github.com/ianhaycox/ir-standings/irsdk"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

const (
	defaultRefreshSeconds = 5
	defaultWidth          = 800
	defaultHeight         = 500
	countBestOf           = 10
	showTopN              = 10
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

	var sdk *irsdk.IRSDK

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

	// Create an instance of the app structure
	app := NewApp(sdk, ir, pointsPerSplit, refreshSeconds, countBestOf, int(iracing.KamelSeriesID), showTopN)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "iRacing Championship Standings",
		Width:  defaultWidth,
		Height: defaultHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		AlwaysOnTop:      true,
		Windows:          &windows.Options{WebviewIsTransparent: true, WindowIsTranslucent: false},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
