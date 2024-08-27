package predictor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ianhaycox/ir-standings/irsdk/telemetry"
	"github.com/ianhaycox/ir-standings/model/championship/car"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/test/files"
)

var (
	pointsPerSplit = points.PointsPerSplit{
		//   0   1   2   3   4   5   6   7   8   9  10 11 12 13 14 15 16 17 18 19
		0: {25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		1: {14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		2: {9, 6, 4, 3, 2, 1},
	}
	carClasses = car.NewCarClasses([]int{83, 84},
		[]cars.Car{{CarID: 76, CarName: "Audi"}, {CarID: 77, CarName: "Nissan"}},
		[]cars.CarClass{
			{CarClassID: 83, Name: "Audi GTO", ShortName: "GTO", CarsInClass: []cars.CarsInClass{{CarID: 76}}},
			{CarClassID: 84, Name: "Nissan GTP", ShortName: "GTP", CarsInClass: []cars.CarsInClass{{CarID: 77}}},
		})
)

func TestPredictorFirstRaceConnected(t *testing.T) {
	var telem telemetry.TelemetryData

	json.Unmarshal(files.ReadFile(t, "/tmp/telemetry.json"), &telem)

	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live([]results.Result{}, &telem)

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}

func TestPredictorFirstRaceNotConnected(t *testing.T) {
	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live([]results.Result{}, &telemetry.TelemetryData{})

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}

func TestPredictorInSeasonConnected(t *testing.T) {
	var telem telemetry.TelemetryData

	json.Unmarshal(files.ReadFile(t, "/tmp/telemetry.json"), &telem)

	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live(files.ReadResultsFixture(t, "../model/fixtures/2024-2-285-results-redacted.json"), &telem)

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}

func TestPredictorInSeasonNotConnected(t *testing.T) {
	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live(files.ReadResultsFixture(t, "../model/fixtures/2024-2-285-results-redacted.json"), &telemetry.TelemetryData{})

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}
