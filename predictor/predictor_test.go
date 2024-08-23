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

func TestPredictorFirstRace(t *testing.T) {
	var pointsPerSplit = points.PointsPerSplit{
		//   0   1   2   3   4   5   6   7   8   9  10 11 12 13 14 15 16 17 18 19
		0: {25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		1: {14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		2: {9, 6, 4, 3, 2, 1},
	}

	var telem telemetry.TelemetryData

	json.Unmarshal(files.ReadFile(t, "/tmp/telemetry.json"), &telem)

	p := NewPredictor(pointsPerSplit, []results.Result{}, &telem, 10)

	ps := p.Live()

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}

func TestPredictorInSeason(t *testing.T) {
	var pointsPerSplit = points.PointsPerSplit{
		//   0   1   2   3   4   5   6   7   8   9  10 11 12 13 14 15 16 17 18 19
		0: {25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		1: {14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		2: {9, 6, 4, 3, 2, 1},
	}

	var telem telemetry.TelemetryData

	json.Unmarshal(files.ReadFile(t, "/tmp/telemetry.json"), &telem)

	p := NewPredictor(pointsPerSplit, files.ReadResultsFixture(t, "../model/fixtures/2024-2-285-results-redacted.json"), &telem, 10)

	ps := p.Live()

	b, _ := json.MarshalIndent(ps, "", "  ")

	fmt.Println(string(b))
}
