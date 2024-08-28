package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/ianhaycox/ir-standings/arch"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	"github.com/ianhaycox/ir-standings/irsdk"
	"github.com/ianhaycox/ir-standings/irsdk/telemetry"
	"github.com/ianhaycox/ir-standings/model/championship/car"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/ianhaycox/ir-standings/predictor"
)

// App struct
type App struct {
	ctx context.Context
	mtx sync.Mutex

	carclasses     car.CarClasses         // Car classes and car membership for this series
	pastResults    []results.Result       // Previous weeks results for this season from the iRacing API
	telemetryData  telemetry.Data         // Shared memory updated every `refreshSeconds`
	prediction     *predictor.Predictor   // Predict standing using past results and current telemetry
	irAPI          iracing.IracingService // iRacing API
	pointsPerSplit points.PointsPerSplit  // Points structure
	refreshSeconds int                    // How often to read telemetry
	countBestOf    int                    // Count best of n races in season
	seriesID       int                    // iRacing series ID
	seasonYear     int                    // E.g. 2024, 2025
	seasonQuarter  int                    // E.g. 1,2,3
	showTopN       int                    // Display top n standings
}

// NewApp creates a new App application struct
func NewApp(sdk *irsdk.IRSDK, irAPI iracing.IracingService, pointsPerSplit points.PointsPerSplit, refreshSeconds, countBestOf, seriesID, showTopN int) *App {
	return &App{
		irAPI:          irAPI,
		refreshSeconds: refreshSeconds,
		seriesID:       seriesID,
		telemetryData:  telemetry.NewData(sdk),
		pointsPerSplit: pointsPerSplit,
		countBestOf:    countBestOf,
		showTopN:       showTopN,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if runtime.GOOS == "windows" {
		arch.WindowOptions()
	}
}

func (a *App) ShowTopN() int {
	return a.showTopN
}

func (a *App) Login(email string, password string) bool { // TODO return error
	if email == "test" {
		const fakeDelay = 5

		filename := "./model/fixtures/2024-2-285-results-redacted.json"
		log.Println("Using results fixture in dev mode:", filename)

		buf, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		a.pastResults = make([]results.Result, 0)

		err = json.Unmarshal(buf, &a.pastResults)
		if err != nil {
			log.Fatal(err)
		}

		a.seasonQuarter = 2
		a.seasonYear = 2024
		a.carclasses = car.NewCarClasses([]int{83, 84},
			[]cars.Car{{CarID: 76, CarName: "Audi"}, {CarID: 77, CarName: "Nissan"}},
			[]cars.CarClass{
				{CarClassID: 83, Name: "Audi GTO", ShortName: "GTO", CarsInClass: []cars.CarsInClass{{CarID: 76}}},
				{CarClassID: 84, Name: "Nissan GTP", ShortName: "GTP", CarsInClass: []cars.CarsInClass{{CarID: 77}}},
			})

		time.Sleep(fakeDelay * time.Second)

		return true
	}

	err := a.irAPI.Authenticate(a.ctx, email, password)
	if err != nil {
		log.Println("Can not login: ", err)

		return false
	}

	cars, err := a.irAPI.Cars(a.ctx)
	if err != nil {
		log.Println("can not get cars:", err)
	}

	carclasses, err := a.irAPI.CarClasses(a.ctx)
	if err != nil {
		log.Println("can not get car classes:", err)
	}

	seasons, err := a.irAPI.Seasons(a.ctx)
	if err != nil {
		log.Println("can not get series results:", err)
	}

	var carClassIDs []int

	for i := range seasons {
		if seasons[i].SeriesID == a.seriesID {
			a.seasonYear = seasons[i].SeasonYear
			a.seasonQuarter = seasons[i].SeasonQuarter
			carClassIDs = seasons[i].CarClassIds

			break
		}
	}

	a.carclasses = car.NewCarClasses(carClassIDs, cars, carclasses)

	log.Println("Getting results for :", a.seasonYear, a.seasonQuarter)

	searchSeriesResults, err := a.irAPI.SearchSeriesResults(a.ctx, a.seasonYear, a.seasonQuarter, a.seriesID)
	if err != nil {
		log.Println("can not get series results:", err)

		return false
	}

	a.pastResults, err = a.irAPI.SeasonBroadcastResults(a.ctx, searchSeriesResults)
	if err != nil {
		log.Println("can not get series results:", err)

		return false
	}

	log.Println("Got results from iRacing API")

	return true
}

func (a *App) LatestStandings() live.PredictedStandings {
	log.Println("LatestStandings")

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.prediction == nil {
		a.prediction = predictor.NewPredictor(a.pointsPerSplit, a.countBestOf, a.carclasses)
	}

	data := a.telemetryData.Telemetry()

	log.Println(data.Status, " ", data.SessionType)

	return a.prediction.Live(a.pastResults, data)
}
