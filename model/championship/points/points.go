// Package points to be awarded per split
package points

import (
	"github.com/ianhaycox/ir-standings/model"
)

type PointsPerSplit map[model.SplitNum][]model.Point

type PointsStructure struct {
	structure PointsPerSplit
}

func NewPointsStructure(structure PointsPerSplit) PointsStructure {
	return PointsStructure{
		structure: structure,
	}
}

func (ps *PointsStructure) Award(splitNum model.SplitNum, finishingPosition model.FinishPositionInClass) model.Point {
	if _, ok := ps.structure[splitNum]; !ok {
		return model.NotCounted
	}

	if len(ps.structure[splitNum]) > int(finishingPosition) {
		return ps.structure[splitNum][finishingPosition]
	}

	return model.Point(0)
}
