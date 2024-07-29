package main

import (
	"C"
)
import "github.com/ianhaycox/ir-standings/live"

//export LiveStandings
func LiveStandings(jsonLivePositions string) (*C.char, int) {
	jsonLiveStandings, err := live.Live(jsonLivePositions)
	if err != nil {
		return C.CString(err.Error()), 1
	}

	return C.CString(jsonLiveStandings), 0
}

func main() {}
