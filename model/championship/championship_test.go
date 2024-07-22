package championship

import (
	"testing"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/driver"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/race"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChampionship(t *testing.T) {
	pointsPerSplit := points.PointsPerSplit{
		0: []int{5, 3, 1},
		1: []int{3, 1},
	}

	t.Run("Loads a simplified set of results into Events and Races", func(t *testing.T) {
		c := NewChampionship(1, map[int]bool{1: true}, 2, points.NewPointsStructure(pointsPerSplit))

		start1, err := time.Parse(time.RFC3339, "2024-03-16T17:00:00Z")
		require.NoError(t, err)

		start2, err := time.Parse(time.RFC3339, "2024-03-23T17:00:00Z")
		require.NoError(t, err)

		fixture := []results.Result{
			{
				SessionID:    1,
				SubsessionID: 1001,
				SessionSplits: []results.SessionSplits{
					{SubsessionID: 1001},
					{SubsessionID: 1002},
				},
				CarClasses: []results.CarClasses{
					{CarClassID: 84, ShortName: "Nissan GTP", Name: "Nissan GTP", CarsInClass: []results.CarsInClass{{CarID: 77}}},
					{CarClassID: 83, ShortName: "Audi 90 GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},
				},
				SeriesID:   285,
				SeriesName: "IMSA Vintage Series",
				SessionResults: []results.SessionResults{
					{
						SimsessionName: "RACE",
						Results: []results.Results{
							{CustID: 9001, DisplayName: "Driver-9001", FinishPositionInClass: 1, LapsComplete: 30, CarClassID: 84},
							{CustID: 9002, DisplayName: "Driver-9002", FinishPositionInClass: 2, LapsComplete: 30, CarClassID: 84},
							{CustID: 9003, DisplayName: "Driver-9003", FinishPositionInClass: 1, LapsComplete: 25, CarClassID: 83},
							{CustID: 9004, DisplayName: "Driver-9004", FinishPositionInClass: 2, LapsComplete: 24, CarClassID: 83},
						},
					},
				},
				StartTime: start1,
				Track:     results.ResultTrack{TrackID: 219, TrackName: "Mount Panorama Circuit"},
			},
			{
				SessionID:    1,
				SubsessionID: 1002,
				SessionSplits: []results.SessionSplits{
					{SubsessionID: 1001},
					{SubsessionID: 1002},
				},
				CarClasses: []results.CarClasses{
					{CarClassID: 84, ShortName: "Nissan GTP", Name: "Nissan GTP", CarsInClass: []results.CarsInClass{{CarID: 77}}},
					{CarClassID: 83, ShortName: "Audi 90 GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},
				},
				SeriesID:   285,
				SeriesName: "IMSA Vintage Series",
				SessionResults: []results.SessionResults{
					{
						SimsessionName: "RACE",
						Results: []results.Results{
							{CustID: 8001, DisplayName: "Driver-8001", FinishPositionInClass: 1, LapsComplete: 29, CarClassID: 84},
							{CustID: 8002, DisplayName: "Driver-8002", FinishPositionInClass: 2, LapsComplete: 28, CarClassID: 84},
							{CustID: 8003, DisplayName: "Driver-8003", FinishPositionInClass: 1, LapsComplete: 24, CarClassID: 83},
							{CustID: 8004, DisplayName: "Driver-8004", FinishPositionInClass: 2, LapsComplete: 24, CarClassID: 83},
						},
					},
				},
				StartTime: start2,
				Track:     results.ResultTrack{TrackID: 219, TrackName: "Mount Panorama Circuit"},
			},
		}

		c.LoadRaceData(fixture)

		// One event with two splits
		events := c.Events()
		assert.Len(t, events, 1)

		expectedDrivers := make(map[model.CustID]driver.Driver)
		expectedDrivers[8001] = driver.NewDriver(8001, "Driver-8001")
		expectedDrivers[8002] = driver.NewDriver(8002, "Driver-8002")
		expectedDrivers[8003] = driver.NewDriver(8003, "Driver-8003")
		expectedDrivers[8004] = driver.NewDriver(8004, "Driver-8004")
		expectedDrivers[9001] = driver.NewDriver(9001, "Driver-9001")
		expectedDrivers[9002] = driver.NewDriver(9002, "Driver-9002")
		expectedDrivers[9003] = driver.NewDriver(9003, "Driver-9003")
		expectedDrivers[9004] = driver.NewDriver(9004, "Driver-9004")

		assert.Len(t, c.drivers, len(expectedDrivers))
		assert.Equal(t, expectedDrivers, c.drivers)

		subsessionIDs := events[0].SubSessions()
		assert.Len(t, subsessionIDs, 2)

		race1, err := events[0].Race(subsessionIDs[0])
		assert.NoError(t, err)
		assert.Equal(t, map[model.CarClassID]int{84: 30, 83: 25}, race1.WinnerLapsComplete())

		race2, err := events[0].Race(subsessionIDs[1])
		assert.NoError(t, err)
		assert.Equal(t, map[model.CarClassID]int{84: 29, 83: 24}, race2.WinnerLapsComplete())

		expectedPositions1 := map[model.CarClassID]map[model.CustID]race.Position{
			84: {
				9001: race.NewPosition(30, 0, 1),
				9002: race.NewPosition(30, 0, 2),
			},
			83: {
				9003: race.NewPosition(25, 0, 1),
				9004: race.NewPosition(24, 0, 2),
			},
		}
		positions1 := race1.Positions()
		assert.Equal(t, expectedPositions1, positions1)

		expectedPositions2 := map[model.CarClassID]map[model.CustID]race.Position{
			84: {
				8001: race.NewPosition(29, 1, 1),
				8002: race.NewPosition(28, 1, 2),
			},
			83: {
				8003: race.NewPosition(24, 1, 1),
				8004: race.NewPosition(24, 1, 2),
			},
		}
		positions2 := race2.Positions()
		assert.Equal(t, expectedPositions2, positions2)
	})

	t.Run("Verify an excluded track is ignored", func(t *testing.T) {
		c := NewChampionship(1, map[int]bool{219: true}, 2, points.NewPointsStructure(pointsPerSplit))

		fixture := []results.Result{
			{
				SessionID:    1,
				SubsessionID: 1001,
				Track:        results.ResultTrack{TrackID: 219, TrackName: "Mount Panorama Circuit"},
			},
		}

		c.LoadRaceData(fixture)

		// No events as 219 excluded
		assert.Len(t, c.Events(), 0)
	})
}
