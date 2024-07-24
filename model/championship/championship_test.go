package championship

import (
	"testing"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/driver"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/position"
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
		c := NewChampionship(1, map[int]bool{1: true}, 2, points.NewPointsStructure(pointsPerSplit), 9)

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
							{CustID: 9001, DisplayName: "Driver-9001", FinishPositionInClass: 1, LapsComplete: 30, CarClassID: 84, CarID: 77},
							{CustID: 9002, DisplayName: "Driver-9002", FinishPositionInClass: 2, LapsComplete: 30, CarClassID: 84, CarID: 77},
							{CustID: 9003, DisplayName: "Driver-9003", FinishPositionInClass: 1, LapsComplete: 25, CarClassID: 83, CarID: 76},
							{CustID: 9004, DisplayName: "Driver-9004", FinishPositionInClass: 2, LapsComplete: 24, CarClassID: 83, CarID: 76},
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
							{CustID: 8001, DisplayName: "Driver-8001", FinishPositionInClass: 1, LapsComplete: 29, CarClassID: 84, CarID: 77},
							{CustID: 8002, DisplayName: "Driver-8002", FinishPositionInClass: 2, LapsComplete: 28, CarClassID: 84, CarID: 77},
							{CustID: 8003, DisplayName: "Driver-8003", FinishPositionInClass: 1, LapsComplete: 24, CarClassID: 83, CarID: 76},
							{CustID: 8004, DisplayName: "Driver-8004", FinishPositionInClass: 2, LapsComplete: 24, CarClassID: 83, CarID: 76},
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
		assert.Equal(t, 30, race1.WinnerLapsComplete(84))
		assert.Equal(t, 25, race1.WinnerLapsComplete(83))

		race2, err := events[0].Race(subsessionIDs[1])
		assert.NoError(t, err)
		assert.Equal(t, 29, race2.WinnerLapsComplete(84))
		assert.Equal(t, 24, race2.WinnerLapsComplete(83))

		expectedPositions84 := map[model.CustID]position.Position{
			9001: position.NewPosition(30, 0, 1, 77),
			9002: position.NewPosition(30, 0, 2, 77),
		}
		positions84 := race1.Positions(84)
		assert.Equal(t, expectedPositions84, positions84)

		expectedPositions83 := map[model.CustID]position.Position{
			9003: position.NewPosition(25, 0, 1, 76),
			9004: position.NewPosition(24, 0, 2, 76),
		}
		positions83 := race1.Positions(83)
		assert.Equal(t, expectedPositions83, positions83)

		expectedPositions84 = map[model.CustID]position.Position{
			8001: position.NewPosition(29, 1, 1, 77),
			8002: position.NewPosition(28, 1, 2, 77),
		}
		positions84 = race2.Positions(84)
		assert.Equal(t, expectedPositions84, positions84)

		expectedPositions83 = map[model.CustID]position.Position{
			8003: position.NewPosition(24, 1, 1, 76),
			8004: position.NewPosition(24, 1, 2, 76),
		}
		positions83 = race2.Positions(83)
		assert.Equal(t, expectedPositions83, positions83)
	})

	t.Run("Verify an excluded track is ignored", func(t *testing.T) {
		c := NewChampionship(1, map[int]bool{219: true}, 2, points.NewPointsStructure(pointsPerSplit), 1)

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
