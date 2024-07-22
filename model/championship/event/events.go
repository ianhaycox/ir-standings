// Package event is subsession race results at a track
package event

import (
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

func (e *Event) StartTime() time.Time {
	return e.startTime
}

func (e *Event) SessionID() model.SessionID {
	return e.sessionID
}

func (e *Event) TrackName() string {
	return e.track.TrackName
}
