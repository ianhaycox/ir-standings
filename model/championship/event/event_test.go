package event

import (
	"testing"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship/race"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	t.Run("NewEvent gets object", func(t *testing.T) {
		now := time.Now().UTC()
		e := NewEvent(model.SessionID(1), now, results.ResultTrack{TrackName: "Silverstone"})

		assert.Equal(t, model.SessionID(1), e.SessionID())
		assert.Equal(t, now, e.StartTime())
		assert.Equal(t, "Silverstone", e.TrackName())
	})

	t.Run("Event race management", func(t *testing.T) {
		now := time.Now().UTC()
		e := NewEvent(model.SessionID(1), now, results.ResultTrack{TrackName: "Silverstone"})

		r2 := race.NewRace(model.SplitNum(1), model.SessionID(1), []model.Result{{SessionID: 6, SubsessionID: 777}})
		e.AddRace(model.SubsessionID(777), r2)

		r1 := race.NewRace(model.SplitNum(0), model.SessionID(1), []model.Result{})
		e.AddRace(model.SubsessionID(888), r1)

		r, err := e.Race(model.SubsessionID(888))
		assert.NoError(t, err)
		assert.Equal(t, &r1, r)

		r, err = e.Race(model.SubsessionID(777))
		assert.NoError(t, err)
		assert.Equal(t, &r2, r)

		r, err = e.Race(model.SubsessionID(666))
		assert.ErrorContains(t, err, "not found")
		assert.Nil(t, r)

		bySplitNum := e.SubSessions()
		assert.Len(t, bySplitNum, 2)
		assert.Equal(t, model.SubsessionID(888), bySplitNum[0])
		assert.Equal(t, model.SubsessionID(777), bySplitNum[1])
	})
}
