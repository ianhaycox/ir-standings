package searchseries

import (
	"encoding/json"
	"testing"

	"github.com/ianhaycox/ir-standings/test/files"
	"github.com/stretchr/testify/assert"
)

func TestIsBroadcast(t *testing.T) {
	var ssr []SearchSeriesResult

	b := files.ReadFile(t, "../../../example/results-search_series-results.json")

	err := json.Unmarshal(b, &ssr)
	assert.NoError(t, err)

	broadcasts := make(map[int]bool)

	for i := range ssr {
		if ssr[i].IsBroadcast() {
			broadcasts[ssr[i].SessionID] = true
		}
	}

	assert.NotEqual(t, 12, len(ssr))
	assert.Len(t, broadcasts, 12)
}
