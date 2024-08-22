package predictor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/telemetry"
	"github.com/ianhaycox/ir-standings/test/files"
)

func TestPredictor(t *testing.T) {
	var pointsPerSplit = points.PointsPerSplit{
		//   0   1   2   3   4   5   6   7   8   9  10 11 12 13 14 15 16 17 18 19
		0: {25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		1: {14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		2: {9, 6, 4, 3, 2, 1},
	}

	files.ReadFile(t, "../telemetry.json")

	var telem telemetry.TelemetryData

	json.Unmarshal(files.ReadFile(t, "/tmp/telemetry.json"), &telem)

	p := NewPredictor(pointsPerSplit, []results.Result{}, &telem, 10)

	ps := p.Live()

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}

/*

	i83 := []live.PredictedStanding{
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

	i84 := []live.PredictedStanding{
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

	s83 := live.Standing{
		SoFByCarClass:           2087,
		CarClassID:              83,
		CarClassName:            "GTO",
		ClassLeaderLapsComplete: 10,
		Items:                   i83,
	}

	s84 := live.Standing{
		SoFByCarClass:           3025,
		CarClassID:              84,
		CarClassName:            "GTP",
		ClassLeaderLapsComplete: 12,
		Items:                   i84,
	}

	ps := live.PredictedStandings{
		TrackName:   "Motegi Resort",
		CountBestOf: 10,
		Standings:   map[model.CarClassID]live.Standing{84: s84, 83: s83},
	}
*/
