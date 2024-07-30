// Package live work out provisional standing based on the current race positions
package live

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
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
	ps               = points.NewPointsStructure(pointsPerSplit)
	current          *championship.Championship
	currentResults   = make([]results.Result, 0)
	currentStandings standings.ChampionshipStandings
)

func Live(jsonLivePositions string) (string, error) {
	var livePositions live.LivePositions

	err := json.Unmarshal([]byte(jsonLivePositions), &livePositions)
	if err != nil {
		return "", fmt.Errorf("malformed request %w", err)
	}

	if current == nil {
		buf, err := os.ReadFile("/home/ian/2024-3-285-results.json")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(buf, &currentResults)
		if err != nil {
			log.Fatal(err)
		}

		current = championship.NewChampionship(model.SeriesID(livePositions.SeriesID), nil, ps, livePositions.CountBestOf)
		current.LoadRaceData(currentResults)
		currentStandings = current.Standings(model.CarClassID(livePositions.CarClassID))
	}

	liveResults := make([]results.Result, 0, len(currentResults))
	liveResults = append(liveResults, currentResults...)

	liveResults = append(liveResults, results.Result{
		SessionID:     livePositions.SessionID,
		SubsessionID:  livePositions.SubsessionID,
		SeriesID:      livePositions.SeriesID,
		SessionSplits: []results.SessionSplits{{SubsessionID: livePositions.SubsessionID}},
		StartTime:     time.Now().UTC(),
		Track:         results.ResultTrack{TrackID: 1, TrackName: "track"},
		SessionResults: []results.SessionResults{
			{
				SimsessionName: "RACE",
				Results:        buildResults(livePositions.CarClassID, livePositions.Results),
			},
		},
	})

	predicted := championship.NewChampionship(model.SeriesID(livePositions.SeriesID), nil, ps, livePositions.CountBestOf)

	predicted.LoadRaceData(liveResults)
	predicted.SetCarClasses(current.CarClasses())

	predictedStandings := predicted.Standings(model.CarClassID(livePositions.CarClassID))

	predictedResult := predictResults(currentStandings, predictedStandings)

	jsonResult, err := json.MarshalIndent(predictedResult[:livePositions.TopN], "", "  ")
	if err != nil {
		return "", fmt.Errorf("can not build response, %w", err)
	}

	return string(jsonResult), nil
}

func predictResults(currentStandings, predictedStandings standings.ChampionshipStandings) []live.LiveStandings {
	liveStandings := make(map[model.CustID]live.LiveStandings)

	for _, entry := range currentStandings.Table {
		liveStandings[entry.CustID] = live.LiveStandings{
			CurrentPosition: int(entry.Position),
			DriverName:      entry.DriverName,
			CurrentPoints:   int(entry.DroppedRoundPoints),
		}
	}

	for _, entry := range predictedStandings.Table {
		if _, ok := liveStandings[entry.CustID]; ok {
			ls := liveStandings[entry.CustID]

			ls.PredictedPoints = int(entry.DroppedRoundPoints)
			ls.PredictedPosition = int(entry.Position)

			liveStandings[entry.CustID] = ls
		} else {
			liveStandings[entry.CustID] = live.LiveStandings{
				PredictedPosition: int(entry.Position),
				DriverName:        entry.DriverName,
				PredictedPoints:   int(entry.DroppedRoundPoints),
			}
		}
	}

	predictedResult := make([]live.LiveStandings, 0, len(liveStandings))

	for custID := range liveStandings {
		ls := liveStandings[custID]
		ls.Change = ls.CurrentPosition - ls.PredictedPosition
		predictedResult = append(predictedResult, ls)
	}

	sort.SliceStable(predictedResult, func(i, j int) bool {
		return predictedResult[i].PredictedPosition < predictedResult[j].PredictedPosition
	})

	return predictedResult
}

func buildResults(carClassID int, liveResults []live.LiveResults) []results.Results {
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
