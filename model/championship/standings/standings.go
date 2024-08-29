// Package standings to work out the final scores
package standings

import (
	"sort"

	"github.com/ianhaycox/ir-standings/model"
)

type ChampionshipStandings struct {
	BestOf       int
	CarClassName string
	Table        []ChampionshipTable
}

type TieBreaker struct {
	SubsessionID model.SubsessionID
	Position     model.FinishPositionInClass
}

func NewTieBreaker(subsessionID model.SubsessionID, position model.FinishPositionInClass) TieBreaker {
	return TieBreaker{
		SubsessionID: subsessionID,
		Position:     position,
	}
}

type ChampionshipTable struct {
	Position                model.FinishPositionInClass
	CustID                  model.CustID
	IRating                 int
	DriverName              string
	CarNames                []string
	DroppedRoundPoints      model.Point
	AllRoundsPoints         model.Point  // Tie-breaker: points without drops
	TieBreakFinishPositions []TieBreaker // then: higher number of better positions., i.e. promote driver with more 1st, then 2nd, etc.
	Counted                 int
	TotalLaps               model.LapsComplete
}

func (cs ChampionshipStandings) orderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (cs ChampionshipStandings) Sort() ChampionshipStandings {
	points := func(c1, c2 *ChampionshipTable) bool {
		return c1.DroppedRoundPoints > c2.DroppedRoundPoints
	}

	tieBreakOne := func(c1, c2 *ChampionshipTable) bool {
		return c1.AllRoundsPoints > c2.AllRoundsPoints
	}

	tieBreakTwo := func(c1, c2 *ChampionshipTable) bool {
		var highestPosition model.FinishPositionInClass

		c1HigherPostions := make(map[model.FinishPositionInClass]int)
		c2HigherPostions := make(map[model.FinishPositionInClass]int)

		for _, c1Pos := range c1.TieBreakFinishPositions {
			c1HigherPostions[c1Pos.Position]++

			if c1Pos.Position > highestPosition {
				highestPosition = c1Pos.Position
			}
		}

		for _, c2Pos := range c2.TieBreakFinishPositions {
			c2HigherPostions[c2Pos.Position]++

			if c2Pos.Position > highestPosition {
				highestPosition = c2Pos.Position
			}
		}

		for i := model.FinishPositionInClass(0); i <= highestPosition; i++ {
			c1Total := c1HigherPostions[i]
			c2Total := c2HigherPostions[i]

			if c1Total == c2Total {
				continue
			}

			return c1Total > c2Total
		}

		return false
	}

	tieBreakFirstRaceOfSeason := func(c1, c2 *ChampionshipTable) bool {
		return c1.IRating > c2.IRating
	}

	cs.orderedBy(points, tieBreakOne, tieBreakTwo, tieBreakFirstRaceOfSeason).Sort(cs.Table)

	for i := range cs.Table {
		cs.Table[i].Position = model.FinishPositionInClass(i + 1)
	}

	return cs
}

type lessFunc func(p1, p2 *ChampionshipTable) bool

type multiSorter struct {
	standings []ChampionshipTable
	less      []lessFunc
}

func (ms *multiSorter) Sort(standings []ChampionshipTable) {
	ms.standings = standings
	sort.Sort(ms)
}

func (ms *multiSorter) Len() int {
	return len(ms.standings)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.standings[i], ms.standings[j] = ms.standings[j], ms.standings[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call.
//
//nolint:wsl // comment
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.standings[i], &ms.standings[j]

	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]

		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}

		// p == q; try the next comparison.
	}

	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}
