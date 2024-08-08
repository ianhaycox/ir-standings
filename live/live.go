// Package live work out provisional standing based on the current race positions
package live

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
)

var (
	pointsPerSplit = points.PointsPerSplit{
		//   0   1   2   3   4   5   6   7   8   9  10 11 12 13 14 15 16 17 18 19
		0: {25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		1: {14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		2: {9, 6, 4, 3, 2, 1},
	}
)

type SafeChamp struct {
	mu                sync.Mutex
	previous          *championship.Championship
	previousResults   []results.Result
	previousStandings standings.ChampionshipStandings
	ps                points.PointsStructure
}

func (c *SafeChamp) load(seriesID model.SeriesID, carClassID model.CarClassID, filename string, countBestOf int) error {
	if c.previous == nil {
		c.previous = championship.NewChampionship(seriesID, nil, c.ps, countBestOf)

		buf, err := os.ReadFile(filename) //nolint:gosec // ok
		if err != nil {
			return fmt.Errorf("can not open file %s", filename)
		}

		err = json.Unmarshal(buf, &c.previousResults)
		if err != nil {
			return fmt.Errorf("can not parse file %s", filename)
		}

		c.previous.LoadRaceData(c.previousResults)
		c.previousStandings = c.previous.Standings(carClassID)
	}

	return nil
}

func (c *SafeChamp) Live(filename string, jsonCurrentPositions string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ps = points.NewPointsStructure(pointsPerSplit)

	var currentPositions live.LiveResults

	fmt.Println("file:", filename, "JSON:", jsonCurrentPositions)

	err := json.Unmarshal([]byte(jsonCurrentPositions), &currentPositions)
	if err != nil {
		return "", fmt.Errorf("malformed request %w", err)
	}

	err = c.load(model.SeriesID(currentPositions.SeriesID), model.CarClassID(currentPositions.CarClassID), filename, currentPositions.CountBestOf)
	if err != nil {
		return "", fmt.Errorf("can load previous results, err:%w", err)
	}

	liveResults := make([]results.Result, 0, len(c.previousResults))
	liveResults = append(liveResults, c.previousResults...)

	liveResults = append(liveResults, results.Result{
		SessionID:     currentPositions.SessionID,
		SubsessionID:  currentPositions.SubsessionID,
		SeriesID:      currentPositions.SeriesID,
		SessionSplits: []results.SessionSplits{{SubsessionID: currentPositions.SubsessionID}},
		StartTime:     time.Now().UTC(),
		Track:         results.ResultTrack{TrackID: 1, TrackName: "track"},
		SessionResults: []results.SessionResults{
			{
				SimsessionName: "RACE",
				Results:        buildResults(currentPositions.CarClassID, currentPositions.Positions),
			},
		},
	})

	predicted := championship.NewChampionship(model.SeriesID(currentPositions.SeriesID), nil, c.ps, currentPositions.CountBestOf)

	predicted.LoadRaceData(liveResults)
	predicted.SetCarClasses(c.previous.CarClasses())

	predictedStandings := predicted.Standings(model.CarClassID(currentPositions.CarClassID))

	provisionalChampionship := provisionalTable(c.previousStandings, predictedStandings)

	jsonResult, err := json.MarshalIndent(provisionalChampionship[:currentPositions.TopN], "", "  ")
	if err != nil {
		return "", fmt.Errorf("can not build response, %w", err)
	}

	return string(jsonResult), nil
}

// provisionalTable calculate change between current and predicted championship tables for the Windows overlay
func provisionalTable(currentStandings, predictedStandings standings.ChampionshipStandings) []live.PredictedStanding {
	mergedStandings := make(map[model.CustID]live.PredictedStanding)

	for _, entry := range currentStandings.Table {
		mergedStandings[entry.CustID] = live.PredictedStanding{
			CurrentPosition: int(entry.Position),
			CustID:          int(entry.CustID),
			DriverName:      entry.DriverName,
			CurrentPoints:   int(entry.DroppedRoundPoints),
		}
	}

	for _, entry := range predictedStandings.Table {
		if _, ok := mergedStandings[entry.CustID]; ok {
			ls := mergedStandings[entry.CustID]

			ls.PredictedPoints = int(entry.DroppedRoundPoints)
			ls.PredictedPosition = int(entry.Position)

			mergedStandings[entry.CustID] = ls
		} else {
			mergedStandings[entry.CustID] = live.PredictedStanding{
				PredictedPosition: int(entry.Position),
				CustID:            int(entry.CustID),
				DriverName:        entry.DriverName,
				PredictedPoints:   int(entry.DroppedRoundPoints),
			}
		}
	}

	predictedResult := make([]live.PredictedStanding, 0, len(mergedStandings))

	for custID := range mergedStandings {
		ls := mergedStandings[custID]
		ls.Change = ls.CurrentPosition - ls.PredictedPosition
		predictedResult = append(predictedResult, ls)
	}

	sort.SliceStable(predictedResult, func(i, j int) bool {
		return predictedResult[i].PredictedPosition < predictedResult[j].PredictedPosition
	})

	return predictedResult
}

func buildResults(carClassID int, liveResults []live.CurrentPosition) []results.Results {
	res := make([]results.Results, 0, len(liveResults))

	for _, lr := range liveResults {
		res = append(res, results.Results{
			CustID:                lr.CustID,
			FinishPositionInClass: lr.FinishPositionInClass,
			LapsComplete:          lr.LapsComplete,
			CarID:                 lr.CarID,
			CarClassID:            carClassID,
		})
	}

	return res
}
