package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandingSort(t *testing.T) {
	t.Run("No entries should produce an empty table", func(t *testing.T) {
		cs := ChampionshipStandings{
			Table: []ChampionshipTable{},
		}

		cs.Sort()

		expected := []ChampionshipTable{}
		assert.Equal(t, expected, cs.Table)
	})

	t.Run("Unique entries should order from high to low for dropped rounds", func(t *testing.T) {
		cs := ChampionshipStandings{
			Table: []ChampionshipTable{
				{DroppedRoundPoints: 23},
				{DroppedRoundPoints: 20},
				{DroppedRoundPoints: 33},
				{DroppedRoundPoints: 22},
			},
		}

		cs.Sort()

		expected := []ChampionshipTable{
			{DroppedRoundPoints: 33, Position: 1},
			{DroppedRoundPoints: 23, Position: 2},
			{DroppedRoundPoints: 22, Position: 3},
			{DroppedRoundPoints: 20, Position: 4},
		}
		assert.Equal(t, expected, cs.Table)
	})

	t.Run("Tie break should invoke all rounds points", func(t *testing.T) {
		cs := ChampionshipStandings{
			Table: []ChampionshipTable{
				{DriverName: "second", DroppedRoundPoints: 23, AllRoundsPoints: 25, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 12}}},
				{DriverName: "second tied better no dropped rounds", DroppedRoundPoints: 23, AllRoundsPoints: 26, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 10}}},
				{DriverName: "first", DroppedRoundPoints: 30, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 8}}},
				{DriverName: "third", DroppedRoundPoints: 22, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 31}}},
			},
		}

		cs.Sort()

		expected := []ChampionshipTable{
			{Position: 1, DriverName: "first", DroppedRoundPoints: 30, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 8}}},
			{Position: 2, DriverName: "second tied better no dropped rounds", DroppedRoundPoints: 23, AllRoundsPoints: 26, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 10}}},
			{Position: 3, DriverName: "second", DroppedRoundPoints: 23, AllRoundsPoints: 25, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 12}}},
			{Position: 4, DriverName: "third", DroppedRoundPoints: 22, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 31}}},
		}
		assert.Equal(t, expected, cs.Table)
	})

	t.Run("Tie break should invoke best finishing positions", func(t *testing.T) {
		cs := ChampionshipStandings{
			Table: []ChampionshipTable{
				{DriverName: "second", DroppedRoundPoints: 23, AllRoundsPoints: 25, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 12}}},
				{DriverName: "second tied better finishing", DroppedRoundPoints: 23, AllRoundsPoints: 25, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 10}}},
				{DriverName: "first", DroppedRoundPoints: 30, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 8}}},
				{DriverName: "third", DroppedRoundPoints: 22, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 31}}},
			},
		}

		cs.Sort()

		expected := []ChampionshipTable{
			{Position: 1, DriverName: "first", DroppedRoundPoints: 30, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 8}}},
			{Position: 2, DriverName: "second tied better finishing", DroppedRoundPoints: 23, AllRoundsPoints: 25, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 10}}},
			{Position: 3, DriverName: "second", DroppedRoundPoints: 23, AllRoundsPoints: 25, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 12}}},
			{Position: 4, DriverName: "third", DroppedRoundPoints: 22, TieBreakFinishPositions: []TieBreaker{{subsessionID: 1, position: 31}}},
		}
		assert.Equal(t, expected, cs.Table)
	})
}
