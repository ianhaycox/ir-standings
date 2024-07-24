package model

import (
	"cmp"
	"slices"
)

type ChampionshipStandings struct {
	BestOf       int
	CarClassName string
	Table        []ChampionshipTable
}

type ChampionshipTable struct {
	Position                int
	DriverName              string
	CarNames                string
	DroppedRoundPoints      int
	AllRoundsPoints         int // Tie-breaker: points without drops
	TieBreakFinishPositions int // then: higher number of better positions.
	Counted                 int
	TotalLaps               int
}

func (cs ChampionshipStandings) Sort() ChampionshipStandings {
	slices.SortFunc(cs.Table, func(a, b ChampionshipTable) int {
		return cmp.Or(
			cmp.Compare(-a.DroppedRoundPoints, -b.DroppedRoundPoints),         // High to low
			cmp.Compare(-a.AllRoundsPoints, -b.AllRoundsPoints),               // High to low
			cmp.Compare(a.TieBreakFinishPositions, b.TieBreakFinishPositions), // Low to high
		)
	})

	for i := range cs.Table {
		cs.Table[i].Position = i + 1
	}

	return cs
}
