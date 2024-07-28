// Package position finishing position for a driver
package position

import (
	"cmp"
	"slices"
	"sort"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
)

type Position struct {
	subsessionID model.SubsessionID
	classified   bool
	lapsComplete model.LapsComplete
	position     model.FinishPositionInClass
	points       model.Point
	carID        model.CarID
}

func NewPosition(subsessionID model.SubsessionID, classified bool, lapsComplete model.LapsComplete, position model.FinishPositionInClass,
	points model.Point, carID model.CarID) Position {
	return Position{
		subsessionID: subsessionID,
		classified:   classified,
		lapsComplete: lapsComplete,
		position:     position,
		points:       points,
		carID:        carID,
	}
}

func (o Position) SubsessionID() model.SubsessionID {
	return o.subsessionID
}

func (o Position) IsClassified() bool {
	return o.classified
}

func (o Position) Position() model.FinishPositionInClass {
	return o.position
}

func (o Position) LapsComplete() model.LapsComplete {
	return o.lapsComplete
}

func (o Position) Points() model.Point {
	return o.points
}

func (o Position) CarID() model.CarID {
	return o.carID
}

type Positions []Position

func (p Positions) BestResults(countBestOf int) Positions {
	bestResults := make(Positions, 0)

	bestResults = append(bestResults, p...)

	slices.SortFunc(bestResults, func(a, b Position) int {
		return cmp.Or(
			cmp.Compare(-a.Points(), -b.Points()),
			cmp.Compare(-a.LapsComplete(), -b.LapsComplete()),
			cmp.Compare(a.Position(), b.Position()),
		)
	})

	if len(bestResults) > countBestOf {
		return bestResults[:countBestOf]
	}

	return bestResults
}

func (p Positions) Classified(classifiedOnly bool) Positions {
	classifiedResults := make(Positions, 0)

	for i := range p {
		if classifiedOnly && !p[i].classified {
			continue
		}

		classifiedResults = append(classifiedResults, p[i])
	}

	return classifiedResults
}

// Total of all candidate finishing positions
func (p Positions) Total(classifiedOnly bool, countBestOf int) model.Point {
	total := model.Point(0)

	filteredPositions := p.Classified(classifiedOnly).BestResults(countBestOf)

	for i := range filteredPositions {
		if filteredPositions[i].Points() == model.NotCounted {
			continue
		}

		total += filteredPositions[i].Points()
	}

	return total
}

func (p Positions) TieBreakerPositions(classifiedOnly bool, countBestOf int) []standings.TieBreaker {
	filteredPositions := p.Classified(classifiedOnly).BestResults(countBestOf)

	justPositions := make([]standings.TieBreaker, 0, len(filteredPositions))

	for i := range filteredPositions {
		justPositions = append(justPositions, standings.NewTieBreaker(filteredPositions[i].SubsessionID(), filteredPositions[i].Position()))
	}

	return justPositions
}

func (p Positions) Laps(classifiedOnly bool, countBestOf int) model.LapsComplete {
	laps := model.LapsComplete(0)

	filteredPositions := p.Classified(classifiedOnly).BestResults(countBestOf)

	for i := range filteredPositions {
		laps += filteredPositions[i].LapsComplete()
	}

	return laps
}

func (p Positions) Counted(classifiedOnly bool, countBestOf int) int {
	filteredPositions := p.Classified(classifiedOnly).BestResults(countBestOf)

	return len(filteredPositions)
}

func (p Positions) CarsDriven(classifiedOnly bool, countBestOf int) []model.CarID {
	carIDs := make(map[model.CarID]bool)

	filteredPositions := p.Classified(classifiedOnly).BestResults(countBestOf)

	for i := range filteredPositions {
		carIDs[filteredPositions[i].CarID()] = true
	}

	uniqueCarIDs := []model.CarID{}

	for carID := range carIDs {
		uniqueCarIDs = append(uniqueCarIDs, carID)
	}

	sort.SliceStable(uniqueCarIDs, func(i, j int) bool { return uniqueCarIDs[i] > uniqueCarIDs[j] })

	return uniqueCarIDs
}
