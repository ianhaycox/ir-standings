// Package position finishing position for a driver
package position

import (
	"cmp"
	"slices"
	"sort"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/points"
)

type Position struct {
	subsessionID model.SubsessionID
	classified   bool
	lapsComplete int
	splitNum     model.SplitNum
	position     int
	carID        model.CarID
}

func NewPosition(subsessionID model.SubsessionID, classified bool, lapsComplete int, splitNum model.SplitNum, position int, carID model.CarID) Position {
	return Position{
		subsessionID: subsessionID,
		classified:   classified,
		lapsComplete: lapsComplete,
		splitNum:     splitNum,
		position:     position,
		carID:        carID,
	}
}

func (o Position) SubsessionID() model.SubsessionID {
	return o.subsessionID
}

func (o Position) IsClassified() bool {
	return o.classified
}

func (o Position) SplitNum() model.SplitNum {
	return o.splitNum
}

func (o Position) Position() int {
	return o.position
}

func (o Position) LapsComplete() int {
	return o.lapsComplete
}

func (o Position) CarID() model.CarID {
	return o.carID
}

type Positions []Position

func (p Positions) BestPositions(countBestOf int) Positions {
	bestPositions := make(Positions, 0)

	bestPositions = append(bestPositions, p...)

	// Higher splits trump lower splits
	slices.SortFunc(bestPositions, func(a, b Position) int {
		return cmp.Or(
			cmp.Compare(a.Position(), b.Position()),
			cmp.Compare(a.SplitNum(), b.SplitNum()),
			cmp.Compare(-a.LapsComplete(), -b.LapsComplete()),
		)
	})

	if len(bestPositions) > countBestOf {
		return bestPositions[:countBestOf]
	}

	return bestPositions
}

func (p Positions) Classified(classifiedOnly bool) Positions {
	classifiedPositions := make(Positions, 0)

	for i := range p {
		if classifiedOnly && !p[i].classified {
			continue
		}

		classifiedPositions = append(classifiedPositions, p[i])
	}

	return classifiedPositions
}

func (p Positions) Total(ps points.PointsStructure, classifiedOnly bool, countBestOf int) int {
	total := 0

	filteredPositions := p.Classified(classifiedOnly).BestPositions(countBestOf)

	for i := range filteredPositions {
		total += ps.Award(filteredPositions[i].SplitNum(), filteredPositions[i].Position())
	}

	return total
}

func (p Positions) Positions(classifiedOnly bool, countBestOf int) []model.TieBreaker {
	filteredPositions := p.Classified(classifiedOnly).BestPositions(countBestOf)

	justPositions := make([]model.TieBreaker, 0, len(filteredPositions))

	for i := range filteredPositions {
		justPositions = append(justPositions, model.NewTieBreaker(filteredPositions[i].SubsessionID(), filteredPositions[i].Position()))
	}

	return justPositions
}

func (p Positions) Laps(classifiedOnly bool, countBestOf int) int {
	laps := 0

	filteredPositions := p.Classified(classifiedOnly).BestPositions(countBestOf)

	for i := range filteredPositions {
		laps += filteredPositions[i].LapsComplete()
	}

	return laps
}

func (p Positions) Counted(classifiedOnly bool, countBestOf int) int {
	filteredPositions := p.Classified(classifiedOnly).BestPositions(countBestOf)

	return len(filteredPositions)
}

func (p Positions) CarsDriven(classifiedOnly bool, countBestOf int) []model.CarID {
	carIDs := make(map[model.CarID]bool)

	filteredPositions := p.Classified(classifiedOnly).BestPositions(countBestOf)

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
