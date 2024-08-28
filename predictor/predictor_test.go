package predictor

import (
	"encoding/json"
	"testing"

	"github.com/ianhaycox/ir-standings/irsdk/telemetry"
	"github.com/ianhaycox/ir-standings/model/championship/car"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/test/files"
	"github.com/stretchr/testify/assert"
)

const (
	telemetryFile = "./fixtures/telemetry.json"
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

	json.Unmarshal(files.ReadFile(t, telemetryFile), &telem)

	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live([]results.Result{}, &telem)

	assert.Equal(t, 5580, ps.Standings[83].SoFByCarClass)
	assert.Equal(t, "Audi GTO", ps.Standings[83].CarClassName)
	assert.Len(t, ps.Standings[83].Items, 3)

	assert.Equal(t, 5766, ps.Standings[84].SoFByCarClass)
	assert.Equal(t, "Nissan GTP", ps.Standings[84].CarClassName)
	assert.Len(t, ps.Standings[84].Items, 3)
}

func TestPredictorFirstRaceNotConnected(t *testing.T) {
	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live([]results.Result{}, &telemetry.TelemetryData{})

	assert.Equal(t, 0, ps.Standings[83].SoFByCarClass)
	assert.Equal(t, "Audi GTO", ps.Standings[83].CarClassName)
	assert.Len(t, ps.Standings[83].Items, 0)

	assert.Equal(t, 0, ps.Standings[84].SoFByCarClass)
	assert.Equal(t, "Nissan GTP", ps.Standings[84].CarClassName)
	assert.Len(t, ps.Standings[84].Items, 0)
}

func TestPredictorInSeasonConnected(t *testing.T) {
	var telem telemetry.TelemetryData

	json.Unmarshal(files.ReadFile(t, "./fixtures/telemetry.json"), &telem)

	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live(files.ReadResultsFixture(t, "../model/fixtures/2024-2-285-results-redacted.json"), &telem)

	assert.Equal(t, 5580, ps.Standings[83].SoFByCarClass)
	assert.Equal(t, "Audi GTO", ps.Standings[83].CarClassName)
	assert.Len(t, ps.Standings[83].Items, 127)

	assert.Equal(t, 5766, ps.Standings[84].SoFByCarClass)
	assert.Equal(t, "Nissan GTP", ps.Standings[84].CarClassName)
	assert.Len(t, ps.Standings[84].Items, 160)
}

func TestPredictorInSeasonNotConnected(t *testing.T) {
	p := NewPredictor(pointsPerSplit, 10, carClasses)

	ps := p.Live(files.ReadResultsFixture(t, "../model/fixtures/2024-2-285-results-redacted.json"), &telemetry.TelemetryData{})

	assert.Equal(t, 0, ps.Standings[83].SoFByCarClass)
	assert.Equal(t, "Audi GTO", ps.Standings[83].CarClassName)
	assert.Len(t, ps.Standings[83].Items, 126)

	assert.Equal(t, 0, ps.Standings[84].SoFByCarClass)
	assert.Equal(t, "Nissan GTP", ps.Standings[84].CarClassName)
	assert.Len(t, ps.Standings[84].Items, 157)
}
