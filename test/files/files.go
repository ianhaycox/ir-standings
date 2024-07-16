// Package files test utils to do stuff with files
package files

import (
	"os"
	"testing"
)

func ReadFile(t *testing.T, filename string) []byte {
	t.Helper()

	buf, err := os.ReadFile(filename) //nolint:gosec // testing
	if err != nil {
		t.Fatal(err)
	}

	return buf
}
