// Package event is subsession race results at a track
package event

import (
	"fmt"
	"sort"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/race"
	"github.com/ianhaycox/ir-standings/model/data/results"
)

type Event struct {
	startTime time.Time
	sessionID model.SessionID
	track     results.ResultTrack
	race      map[model.SubsessionID]race.Race
}

func NewEvent(sessionID model.SessionID, startTime time.Time, track results.ResultTrack) Event {
	return Event{
		sessionID: sessionID,
		startTime: startTime,
		track:     track,
		race:      make(map[model.SubsessionID]race.Race),
	}
}

func (e *Event) AddRace(subsessionID model.SubsessionID, race race.Race) {
	e.race[subsessionID] = race
}

// SubSessions in order of splitNum
func (e *Event) SubSessions() []model.SubsessionID {
	type sorter struct {
		splitNum     model.SplitNum
		subsessionID model.SubsessionID
	}

	var bySplitNum []sorter

	for subsessionID, race := range e.race {
		bySplitNum = append(bySplitNum, sorter{subsessionID: subsessionID, splitNum: race.SplitNum()})
	}

	sort.SliceStable(bySplitNum, func(i, j int) bool { return bySplitNum[i].splitNum < bySplitNum[j].splitNum })

	subsessionIDs := make([]model.SubsessionID, 0)

	for i := range bySplitNum {
		subsessionIDs = append(subsessionIDs, bySplitNum[i].subsessionID)
	}

	return subsessionIDs
}

func (e *Event) Race(subsessionID model.SubsessionID) (*race.Race, error) {
	if race, ok := e.race[subsessionID]; ok {
		return &race, nil
	}

	return nil, fmt.Errorf("race %d not found", subsessionID)
}

func (e *Event) StartTime() time.Time {
	return e.startTime
}

func (e *Event) SessionID() model.SessionID {
	return e.sessionID
}

func (e *Event) TrackName() string {
	return e.track.TrackName
}
