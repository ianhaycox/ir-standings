package championship

import (
	"time"

	"github.com/ianhaycox/ir-standings/model/data/results"
)

type Event struct {
	startTime time.Time
	sessionID SessionID
	track     results.ResultTrack // In order of event
}
