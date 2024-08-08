package main

// #cgo CFLAGS: -g
import "C"
import (
	"github.com/ianhaycox/ir-standings/live"
)

//export LiveStandings
func LiveStandings(filename string, jsonLivePositions string) (*C.char, int) {
	jsonLiveStandings, err := live.Live(filename, jsonLivePositions)
	if err != nil {
		return C.CString(err.Error()), 1
	}

	return C.CString(jsonLiveStandings), 0
}

func main() {}
