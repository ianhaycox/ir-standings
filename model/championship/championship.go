// Package championship to calculate positions
package championship

import (
	"log"
	"sort"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/driver"
	"github.com/ianhaycox/ir-standings/model/championship/event"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/race"
	"github.com/ianhaycox/ir-standings/model/data/results"
)

type Championship struct {
	seriesID       int
	events         map[model.SessionID]event.Event // Events have multiple races in different splits
	carClasses     []results.CarClasses            // Competing cars
	excludeTrackID map[int]bool                    // Exclude these candidates from the results
	maxSplits      int                             // Ignore more splits than this when calculating results
	points         points.PointsStructure          // Points awarded by split
	drivers        map[model.CustID]driver.Driver
}

func NewChampionship(seriesID int, excludeTrackID map[int]bool, maxSplits int, points points.PointsStructure) *Championship {
	return &Championship{
		seriesID:       seriesID,
		events:         make(map[model.SessionID]event.Event),
		excludeTrackID: excludeTrackID,
		maxSplits:      maxSplits,
		points:         points,
		drivers:        make(map[model.CustID]driver.Driver),
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

func (c *Championship) AddDriver(custID model.CustID, driver driver.Driver) {
	c.drivers[custID] = driver
}

func (c *Championship) LoadRaceData(data []results.Result) {
	if len(data) > 0 {
		c.carClasses = data[0].CarClasses
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
						FinishPosition:          sessionResult.FinishPosition,
						FinishPositionInClass:   sessionResult.FinishPositionInClass,
						LapsLead:                sessionResult.LapsLead,
						LapsComplete:            sessionResult.LapsComplete,
						Position:                sessionResult.Position,
						QualLapTime:             sessionResult.QualLapTime,
						StartingPosition:        sessionResult.StartingPosition,
						StartingPositionInClass: sessionResult.StartingPositionInClass,
						CarClassID:              model.CarClassID(sessionResult.CarClassID),
						CarClassName:            sessionResult.CarClassName,
						CarClassShortName:       sessionResult.CarClassShortName,
						ClubID:                  sessionResult.ClubID,
						ClubName:                sessionResult.ClubName,
						ClubShortname:           sessionResult.ClubShortname,
						Division:                sessionResult.Division,
						DivisionName:            sessionResult.DivisionName,
						Incidents:               sessionResult.Incidents,
						CarID:                   sessionResult.CarID,
						CarName:                 sessionResult.CarName,
					}

					sessionResults = append(sessionResults, sessionResult)

					c.AddDriver(sessionResult.CustID, driver.NewDriver(sessionResult.CustID, sessionResult.DisplayName))
				}
			}

			race := race.NewRace(c.SplitNum(subsessionID, result.SessionSplits), sessionID, sessionResults)

			sessionEvent.AddRace(subsessionID, race)

			c.events[sessionID] = sessionEvent
		}
	}
}

func (c *Championship) SplitNum(subsessionID model.SubsessionID, sessionSplits []results.SessionSplits) model.SplitNum {
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
