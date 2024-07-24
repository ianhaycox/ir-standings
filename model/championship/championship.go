// Package championship to calculate positions
package championship

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/car"
	"github.com/ianhaycox/ir-standings/model/championship/driver"
	"github.com/ianhaycox/ir-standings/model/championship/event"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/position"
	"github.com/ianhaycox/ir-standings/model/championship/race"
	"github.com/ianhaycox/ir-standings/model/data/results"
)

type Championship struct {
	seriesID       int
	events         map[model.SessionID]event.Event // Events have multiple races in different splits
	carClasses     car.CarClasses                  // Competing car classes and names
	excludeTrackID map[int]bool                    // Exclude these candidates from the results
	maxSplits      int                             // Ignore more splits than this when calculating results
	points         points.PointsStructure          // Points awarded by split
	drivers        map[model.CustID]driver.Driver
	countBestOf    int
}

func NewChampionship(seriesID int, excludeTrackID map[int]bool, maxSplits int, points points.PointsStructure, countBestOf int) *Championship {
	return &Championship{
		seriesID:       seriesID,
		events:         make(map[model.SessionID]event.Event),
		excludeTrackID: excludeTrackID,
		maxSplits:      maxSplits,
		points:         points,
		drivers:        make(map[model.CustID]driver.Driver),
		countBestOf:    countBestOf,
	}
}

// Events list of events sorted by start time
func (c *Championship) Events() []event.Event {
	sortedEvents := make([]event.Event, 0, len(c.events))

	for _, event := range c.events {
		sortedEvents = append(sortedEvents, event)
	}

	sort.SliceStable(sortedEvents, func(i, j int) bool { return sortedEvents[i].StartTime().Before(sortedEvents[j].StartTime()) })

	return sortedEvents
}

func (c *Championship) LoadRaceData(data []results.Result) {
	if len(data) > 0 {
		c.carClasses = car.NewCarClasses(data[0].CarClasses)
	}

	for _, result := range data {
		if c.isExcluded(result.Track.TrackID) {
			continue
		}

		sessionID := model.SessionID(result.SessionID)
		subsessionID := model.SubsessionID(result.SubsessionID)

		var (
			sessionEvent event.Event
			ok           bool
		)

		if sessionEvent, ok = c.events[sessionID]; !ok {
			sessionEvent = event.NewEvent(sessionID, result.StartTime, result.Track)
		}

		sessionResults := make([]model.Result, 0)

		for i := range result.SessionResults {
			if result.SessionResults[i].SimsessionName == "RACE" {
				for _, sessionResult := range result.SessionResults[i].Results {
					sessionResult := model.Result{
						SessionID:               model.SessionID(result.SessionID),
						SubsessionID:            model.SubsessionID(result.SubsessionID),
						CustID:                  model.CustID(sessionResult.CustID),
						DisplayName:             sessionResult.DisplayName,
						FinishPositionInClass:   sessionResult.FinishPositionInClass,
						LapsLead:                sessionResult.LapsLead,
						LapsComplete:            sessionResult.LapsComplete,
						Position:                sessionResult.Position,
						QualLapTime:             sessionResult.QualLapTime,
						StartingPosition:        sessionResult.StartingPosition,
						StartingPositionInClass: sessionResult.StartingPositionInClass,
						CarClassID:              model.CarClassID(sessionResult.CarClassID),
						ClubID:                  sessionResult.ClubID,
						ClubName:                sessionResult.ClubName,
						ClubShortname:           sessionResult.ClubShortname,
						Division:                sessionResult.Division,
						DivisionName:            sessionResult.DivisionName,
						Incidents:               sessionResult.Incidents,
						CarID:                   model.CarID(sessionResult.CarID),
						CarName:                 sessionResult.CarName,
					}

					sessionResults = append(sessionResults, sessionResult)

					c.addDriver(sessionResult.CustID, driver.NewDriver(sessionResult.CustID, sessionResult.DisplayName))

					c.carClasses.AddCarName(sessionResult.CarID, sessionResult.CarName)
				}
			}

			race := race.NewRace(c.splitNum(subsessionID, result.SessionSplits), sessionID, sessionResults)

			sessionEvent.AddRace(subsessionID, race)

			c.events[sessionID] = sessionEvent
		}
	}
}

func (c *Championship) Standings(carClassID model.CarClassID) model.ChampionshipStandings {
	cs := model.ChampionshipStandings{
		BestOf:       c.countBestOf,
		CarClassName: "GTP",
		Table:        make([]model.ChampionshipTable, 0),
	}

	events := c.Events()

	classifiedPositions := make(map[model.CustID]position.Positions)
	droppedRoundPositions := make(map[model.CustID]position.Positions)
	allPositions := make(map[model.CustID]position.Positions)
	custIDs := make(map[model.CustID]bool)
	carsDriven := make(map[model.CustID][]model.CarID)

	for _, event := range events {
		subsessionIDs := event.SubSessions()

		for i := range subsessionIDs {
			race, err := event.Race(subsessionIDs[i])
			if err != nil {
				log.Fatal(err)
			}

			winnerLapsComplete := race.WinnerLapsComplete(carClassID)

			racePositions := race.Positions(carClassID)
			for custID, position := range racePositions {
				custIDs[custID] = true

				if custID == 455246 {
					fmt.Println(position)
				}

				carsDriven[custID] = append(carsDriven[custID], position.CarID())
				allPositions[custID] = append(allPositions[custID], position)

				if !race.IsClassified(winnerLapsComplete, position.LapsComplete()) {
					continue
				}

				classifiedPositions[custID] = append(classifiedPositions[custID], position)
				droppedRoundPositions[custID] = classifiedPositions[custID].BestPositions(c.countBestOf)
			}
		}
	}

	for custID := range custIDs {
		driver := c.drivers[custID]

		cs.Table = append(cs.Table, model.ChampionshipTable{
			DroppedRoundPoints:      droppedRoundPositions[custID].Total(c.points),
			AllRoundsPoints:         classifiedPositions[custID].Total(c.points),
			TieBreakFinishPositions: allPositions[custID].Sum(),
			CarNames:                strings.Join(c.carClasses.Names(carsDriven[custID]), ","),
			DriverName:              driver.DisplayName(),
			Counted:                 droppedRoundPositions[custID].Counted(),
			TotalLaps:               droppedRoundPositions[custID].Laps(),
		})
	}

	return cs.Sort()
}

func (c *Championship) splitNum(subsessionID model.SubsessionID, sessionSplits []results.SessionSplits) model.SplitNum {
	for i, sessionSplit := range sessionSplits {
		if sessionSplit.SubsessionID == int(subsessionID) {
			return model.SplitNum(i)
		}
	}

	log.Fatal("can not determine split number for SubsessionID", subsessionID)

	return 0
}

func (c *Championship) isExcluded(trackID int) bool {
	_, ok := c.excludeTrackID[trackID]

	return ok
}

func (c *Championship) addDriver(custID model.CustID, driver driver.Driver) {
	c.drivers[custID] = driver
}

/*
func (s *Season) CalculateChampionshipPoints(points [][]int) {
	for splitNum, splitResult := range s.splits {
		winnerLapsComplete := make(map[SubsessionID]map[CarClassID]int)

		for subsessionID, result := range splitResult.results {
			for _, sessionResult := range result {
				carClassID := sessionResult.CarClassID

				if len(winnerLapsComplete[subsessionID]) == 0 {
					winnerLapsComplete[subsessionID] = make(map[CarClassID]int)
				}

				if sessionResult.LapsComplete > winnerLapsComplete[subsessionID][carClassID] {
					winnerLapsComplete[subsessionID][carClassID] = sessionResult.LapsComplete
				}
			}
		}

		for carClassID, car := range splitResult.pointsByCarClass {
			for custID := range car {
				for subsessionID := range car[custID].finishingPositions {
					if splitResult.results[subsessionID][custID].IsClassified(winnerLapsComplete[subsessionID][carClassID]) {
						if len(points[splitNum]) > splitResult.pointsByCarClass[carClassID][custID].finishingPositions[subsessionID] {
							splitResult.pointsByCarClass[carClassID][custID].championshipsPoints[subsessionID] =
								points[splitNum][splitResult.pointsByCarClass[carClassID][custID].finishingPositions[subsessionID]]
						}
					}
				}

				all := make([]int, 0)
				for _, champPoints := range splitResult.pointsByCarClass[carClassID][custID].championshipsPoints {
					all = append(all, champPoints)
				}

				sort.SliceStable(all, func(i, j int) bool { return all[i] > all[j] })

				bestOf := 0

				for i := 0; i < 9 && i < len(all); i++ {
					bestOf += all[i]
				}

				x := splitResult.pointsByCarClass[carClassID][custID]
				x.bestOf = bestOf
				splitResult.pointsByCarClass[carClassID][custID] = x
			}
		}
	}
}
*/
