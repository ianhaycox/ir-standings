// Package race details for a SessionID per split
package race

import (
	"github.com/ianhaycox/ir-standings/model"
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

type Position struct {
	lapsComplete int
	splitNum     model.SplitNum
	position     int
}

func (r *Race) WinnerLapsComplete() map[model.CarClassID]int {
	winnerLapsComplete := make(map[model.CarClassID]int)

	for i := range r.results {
		if r.results[i].LapsComplete > winnerLapsComplete[r.results[i].CarClassID] {
			winnerLapsComplete[r.results[i].CarClassID] = r.results[i].LapsComplete
		}
	}

	return winnerLapsComplete
}

func (r *Race) Positions() map[model.CarClassID]map[model.CustID]Position {
	finishingPositions := make(map[model.CarClassID]map[model.CustID]Position)

	for _, result := range r.results {
		if len(finishingPositions[result.CarClassID]) == 0 {
			finishingPositions[result.CarClassID] = make(map[model.CustID]Position)
		}

		positionForCust := finishingPositions[result.CarClassID][result.CustID]

		positionForCust.splitNum = r.splitNum
		positionForCust.lapsComplete = result.LapsComplete
		positionForCust.position = result.FinishPositionInClass

		finishingPositions[result.CarClassID][result.CustID] = positionForCust
	}

	return finishingPositions
}
