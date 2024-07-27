// Package race details for a SessionID per split
package race

import (
	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/position"
	"github.com/ianhaycox/ir-standings/model/championship/result"
)

type Race struct {
	splitNum  model.SplitNum
	sessionID model.SessionID
	results   []result.Result
}

func NewRace(splitNum model.SplitNum, sessionID model.SessionID, results []result.Result) Race {
	return Race{
		splitNum:  splitNum,
		sessionID: sessionID,
		results:   results,
	}
}

func (r *Race) SplitNum() model.SplitNum {
	return r.splitNum
}

func (r *Race) WinnerLapsComplete(carClassID model.CarClassID) model.LapsComplete {
	winnerLapsComplete := model.LapsComplete(0)

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

func (r *Race) Positions(carClassID model.CarClassID, winnerLapsComplete model.LapsComplete, awards points.PointsStructure) map[model.CustID]position.Position {
	finishingPositions := make(map[model.CustID]position.Position)

	for _, result := range r.results {
		if result.CarClassID != carClassID {
			continue
		}

		pointsAwarded := awards.Award(r.SplitNum(), result.FinishPositionInClass)

		finishingPositions[result.CustID] = position.NewPosition(result.SubsessionID, r.IsClassified(winnerLapsComplete, result.LapsComplete),
			result.LapsComplete, result.FinishPositionInClass, pointsAwarded, result.CarID)
	}

	return finishingPositions
}

func (r *Race) IsClassified(winnerLapsComplete, lapsComplete model.LapsComplete) bool {
	return lapsComplete*4 >= winnerLapsComplete*3 // 75%
}
