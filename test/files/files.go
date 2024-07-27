// Package files test utils to do stuff with files
package files

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ianhaycox/ir-standings/model/data/results"
)

func ReadFile(t *testing.T, filename string) []byte {
	t.Helper()

	buf, err := os.ReadFile(filename) //nolint:gosec // testing
	if err != nil {
		t.Fatal(err)
	}

	return buf
}

func ReadResultsFixture(t *testing.T, filename string) []results.Result {
	t.Helper()

	buf := ReadFile(t, filename)

	res := make([]results.Result, 0)

	err := json.Unmarshal(buf, &res)
	if err != nil {
		t.Fatal(err)
	}

	return res
}
