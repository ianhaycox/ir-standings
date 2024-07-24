// Package position finishing position for a driver
package position

import (
	"sort"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/points"
)

type Position struct {
	lapsComplete int
	splitNum     model.SplitNum
	position     int
	carID        model.CarID
}

func NewPosition(lapsComplete int, splitNum model.SplitNum, position int, carID model.CarID) Position {
	return Position{
		lapsComplete: lapsComplete,
		splitNum:     splitNum,
		position:     position,
		carID:        carID,
	}
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

	sort.SliceStable(bestPositions, func(i, j int) bool { return bestPositions[i].position < bestPositions[j].position })

	if len(bestPositions) > countBestOf {
		return bestPositions[:countBestOf]
	}

	return bestPositions
}

func (p Positions) Total(ps points.PointsStructure) int {
	total := 0

	for i := range p {
		total += ps.Award(p[i].SplitNum(), p[i].Position())
	}

	return total
}

func (p Positions) Sum() int {
	sum := 0

	for i := range p {
		sum += p[i].Position()
	}

	return sum
}

func (p Positions) Laps() int {
	laps := 0

	for i := range p {
		laps += p[i].LapsComplete()
	}

	return laps
}

func (p Positions) Counted() int {
	return len(p)
}
