// Package race details for a SessionID per split
package race

import (
	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/position"
)

type Race struct {
	splitNum  model.SplitNum
	sessionID model.SessionID
	results   []model.Result
}

func NewRace(splitNum model.SplitNum, sessionID model.SessionID, results []model.Result) Race {
	return Race{
		splitNum:  splitNum,
		sessionID: sessionID,
		results:   results,
	}
}

func (r *Race) SplitNum() model.SplitNum {
	return r.splitNum
}

func (r *Race) IsClassified(winnerLapsComplete int, lapsComplete int) bool {
	return lapsComplete*4 >= winnerLapsComplete*3 // 75%
}

func (r *Race) WinnerLapsComplete(carClassID model.CarClassID) int {
	winnerLapsComplete := 0

	for i := range r.results {
		if r.results[i].CarClassID != carClassID {
			continue
		}

		if r.results[i].LapsComplete > winnerLapsComplete {
			winnerLapsComplete = r.results[i].LapsComplete
		}
	}

	return winnerLapsComplete
}

func (r *Race) Positions(carClassID model.CarClassID) map[model.CustID]position.Position {
	finishingPositions := make(map[model.CustID]position.Position)

	for _, result := range r.results {
		if result.CarClassID != carClassID {
			continue
		}

		finishingPositions[result.CustID] = position.NewPosition(result.LapsComplete, r.splitNum, result.FinishPositionInClass, result.CarID)
	}

	return finishingPositions
}
