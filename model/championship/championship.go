// Package championship to calculate positions
package championship

import (
	"sort"

	"github.com/ianhaycox/ir-standings/model/data/results"
)

type SessionID int
type SubsessionID int
type CarClassID int
type CustID int

type Championship struct {
	seriesID       int
	races          []results.Result // All race,qualifying,practice results from iRacing for this series and season
	excludeTrackID map[int]bool     // Exclude these candidates from the results
	maxSplits      int              // Ignore more splits than this when calculating results
}

func NewChampionship(seriesID int, excludeTrackID map[int]bool, maxSplits int, results []results.Result) *Championship {
	return &Championship{
		seriesID:       seriesID,
		races:          results,
		excludeTrackID: excludeTrackID,
		maxSplits:      maxSplits,
	}
}

func (c *Championship) BroadcastRaces() Season {
	s := Season{
		events:     c.events(),
		carClasses: c.races[0].CarClasses,
	}

	s.splits = make([]SplitResult, 0)

	// Include only races and filter out unnecessary tracks and splits
	for splitNum := 0; splitNum < c.maxSplits; splitNum++ {
		splitResult := SplitResult{
			results:          make(map[SubsessionID]map[CustID]*Result),
			pointsByCarClass: make(map[CarClassID]map[CustID]SeasonPoints),
		}

		for _, result := range c.races {
			if len(result.SessionSplits) <= splitNum {
				continue
			}

			if c.isExcluded(result.Track.TrackID) {
				continue
			}

			subsessionID := result.SessionSplits[splitNum].SubsessionID

			if result.SubsessionID == subsessionID {
				for i := range result.SessionResults {
					if result.SessionResults[i].SimsessionName == "RACE" {
						custResult := make(map[CustID]*Result)

						for _, sessionResult := range result.SessionResults[i].Results {
							custID := CustID(sessionResult.CustID)

							custResult[custID] = &Result{
								SessionID:               SessionID(result.SessionID),
								SubsessionID:            SubsessionID(result.SubsessionID),
								CustID:                  CustID(sessionResult.CustID),
								DisplayName:             sessionResult.DisplayName,
								FinishPosition:          sessionResult.FinishPosition,
								FinishPositionInClass:   sessionResult.FinishPositionInClass,
								LapsLead:                sessionResult.LapsLead,
								LapsComplete:            sessionResult.LapsComplete,
								Position:                sessionResult.Position,
								QualLapTime:             sessionResult.QualLapTime,
								StartingPosition:        sessionResult.StartingPosition,
								StartingPositionInClass: sessionResult.StartingPositionInClass,
								CarClassID:              CarClassID(sessionResult.CarClassID),
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
						}

						splitResult.results[SubsessionID(subsessionID)] = custResult
					}
				}
			}
		}

		if len(splitResult.results) > 0 {
			s.splits = append(s.splits, splitResult)
		}
	}

	return s
}

func (c *Championship) events() []Event {
	tracks := make(map[int]Event)

	for i := range c.races {
		if c.isExcluded(c.races[i].Track.TrackID) {
			continue
		}

		tracks[c.races[i].Track.TrackID] = Event{
			sessionID: SessionID(c.races[i].SessionID),
			startTime: c.races[i].StartTime,
			track:     c.races[i].Track,
		}
	}

	events := make([]Event, 0, len(tracks))

	for _, track := range tracks {
		events = append(events, track)
	}

	sort.SliceStable(events, func(i, j int) bool { return events[i].startTime.Before(events[j].startTime) })

	return events
}

func (c *Championship) isExcluded(trackID int) bool {
	_, ok := c.excludeTrackID[trackID]

	return ok
}
