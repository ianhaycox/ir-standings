// Package points to be awarded per split
package points

import (
	"github.com/ianhaycox/ir-standings/model"
)

type PointsPerSplit map[model.SplitNum][]int

type PointsStructure struct {
	structure PointsPerSplit
}

func NewPointsStructure(structure PointsPerSplit) PointsStructure {
	return PointsStructure{
		structure: structure,
	}
}

func (ps *PointsStructure) Award(splitNum model.SplitNum, finishingPosition int) int {
	if len(ps.structure[splitNum]) > finishingPosition {
		return ps.structure[splitNum][finishingPosition]
	}

	return 0
}
