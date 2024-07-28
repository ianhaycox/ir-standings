// Package championship to calculate positions
package championship

import (
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
	"github.com/ianhaycox/ir-standings/model/championship/result"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/ianhaycox/ir-standings/model/data/results"
)

type Championship struct {
	seriesID       int
	events         map[model.SessionID]event.Event // Events have multiple races in different splits
	carClasses     car.CarClasses                  // Competing car classes and names
	excludeTrackID map[int]bool                    // Exclude these candidates from the results
	awards         points.PointsStructure          // Points awarded by split
	drivers        map[model.CustID]driver.Driver
	countBestOf    int
}

func NewChampionship(seriesID int, excludeTrackID map[int]bool, awards points.PointsStructure, countBestOf int) *Championship {
	return &Championship{
		seriesID:       seriesID,
		events:         make(map[model.SessionID]event.Event),
		excludeTrackID: excludeTrackID,
		awards:         awards,
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

	for _, irResult := range data {
		if c.isExcluded(irResult.Track.TrackID) {
			continue
		}

		sessionID := model.SessionID(irResult.SessionID)
		subsessionID := model.SubsessionID(irResult.SubsessionID)

		var (
			sessionEvent event.Event
			ok           bool
		)

		if sessionEvent, ok = c.events[sessionID]; !ok {
			sessionEvent = event.NewEvent(sessionID, irResult.StartTime, irResult.Track)
		}

		sessionResults := make([]result.Result, 0)

		for i := range irResult.SessionResults {
			if irResult.SessionResults[i].SimsessionName == "RACE" {
				for _, sessionResult := range irResult.SessionResults[i].Results {
					sessionResult := result.Result{
						SessionID:               model.SessionID(irResult.SessionID),
						SubsessionID:            model.SubsessionID(irResult.SubsessionID),
						CustID:                  model.CustID(sessionResult.CustID),
						DisplayName:             sessionResult.DisplayName,
						FinishPositionInClass:   model.FinishPositionInClass(sessionResult.FinishPositionInClass),
						LapsLead:                sessionResult.LapsLead,
						LapsComplete:            model.LapsComplete(sessionResult.LapsComplete),
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

			race := race.NewRace(c.splitNum(subsessionID, irResult.SessionSplits), sessionID, sessionResults)

			sessionEvent.AddRace(subsessionID, race)

			c.events[sessionID] = sessionEvent
		}
	}
}

func (c *Championship) Standings(carClassID model.CarClassID) standings.ChampionshipStandings {
	cs := standings.ChampionshipStandings{
		BestOf:       c.countBestOf,
		CarClassName: "GTP",
		Table:        make([]standings.ChampionshipTable, 0),
	}

	events := c.Events()

	custFinishingPositions := make(map[model.CustID]position.Positions)

	for _, event := range events {
		subsessionIDs := event.SubSessions()

		for i := range subsessionIDs {
			race, err := event.Race(subsessionIDs[i])
			if err != nil {
				log.Fatal(err)
			}

			winnerLapsComplete := race.WinnerLapsComplete(carClassID)

			racePositions := race.Positions(carClassID, winnerLapsComplete, c.awards)
			for custID, position := range racePositions {
				custFinishingPositions[custID] = append(custFinishingPositions[custID], position)
			}
		}
	}

	for custID, positions := range custFinishingPositions {
		driver := c.drivers[custID]

		cs.Table = append(cs.Table, standings.ChampionshipTable{
			DroppedRoundPoints:      positions.Total(true, c.countBestOf),
			AllRoundsPoints:         positions.Total(true, len(events)),
			TieBreakFinishPositions: positions.TieBreakerPositions(true, c.countBestOf),
			CarNames:                strings.Join(c.carClasses.Names(positions.CarsDriven(true, c.countBestOf)), ","),
			DriverName:              driver.DisplayName(),
			Counted:                 positions.Counted(false, c.countBestOf),
			TotalLaps:               positions.Laps(false, c.countBestOf),
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
