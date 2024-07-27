package model

import (
	"sort"
)

type ChampionshipStandings struct {
	BestOf       int
	CarClassName string
	Table        []ChampionshipTable
}

type TieBreaker struct {
	subsessionID SubsessionID
	position     int
}

func NewTieBreaker(subsessionID SubsessionID, position int) TieBreaker {
	return TieBreaker{
		subsessionID: subsessionID,
		position:     position,
	}
}

type ChampionshipTable struct {
	Position                int
	DriverName              string
	CarNames                string
	DroppedRoundPoints      int
	AllRoundsPoints         int          // Tie-breaker: points without drops
	TieBreakFinishPositions []TieBreaker // then: higher number of better positions.
	Counted                 int
	TotalLaps               int
}

func (cs ChampionshipStandings) orderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

type bySession map[SubsessionID]int

func betterFinish(c1Sessions, c2Sessions bySession) int {
	c1Better := 0

	for subsessionID, pos := range c1Sessions {
		if other, ok := c2Sessions[subsessionID]; ok {
			if pos < other {
				c1Better++
			}
		} else {
			// No comparative session
			c1Better++
		}
	}

	return c1Better
}

func (cs ChampionshipStandings) Sort() ChampionshipStandings {
	points := func(c1, c2 *ChampionshipTable) bool {
		return c1.DroppedRoundPoints > c2.DroppedRoundPoints
	}

	tieBreakOne := func(c1, c2 *ChampionshipTable) bool {
		return c1.AllRoundsPoints > c2.AllRoundsPoints
	}

	tieBreakTwo := func(c1, c2 *ChampionshipTable) bool {
		return len(c1.TieBreakFinishPositions) > len(c2.TieBreakFinishPositions)
	}

	tieBreakThree := func(c1, c2 *ChampionshipTable) bool {
		c1Sessions := make(bySession)
		c2Sessions := make(bySession)

		for _, pos := range c1.TieBreakFinishPositions {
			c1Sessions[pos.subsessionID] = pos.position
		}

		for _, pos := range c2.TieBreakFinishPositions {
			c2Sessions[pos.subsessionID] = pos.position
		}

		c1Better := betterFinish(c1Sessions, c2Sessions)
		c2Better := betterFinish(c2Sessions, c1Sessions)

		return c1Better < c2Better // Not less()
	}

	tieBreakFour := func(c1, c2 *ChampionshipTable) bool {
		return c1.Counted > c2.Counted
	}

	// According to #vcr_rules the second tie break, higher number of better positions, should be enough, but it
	// doesn't seem to be enough. By trial and error tieBreakTwo seems extra and also need tieBreakFour
	// to make low numbers of scored rounds deterministic.
	cs.orderedBy(points, tieBreakOne, tieBreakTwo, tieBreakThree, tieBreakFour).Sort(cs.Table)

	for i := range cs.Table {
		cs.Table[i].Position = i + 1
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
