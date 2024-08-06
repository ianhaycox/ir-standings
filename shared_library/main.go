package main

// #cgo CFLAGS: -g
import "C"
import (
	"os"

	"github.com/ianhaycox/ir-standings/live"
)

//export LiveStandings
func LiveStandings(filename string, jsonLivePositions string) (*C.char, int) {
	jsonLiveStandings, err := live.Live(filename, jsonLivePositions)
	if err != nil {
		return C.CString(err.Error()), 1
	}

	os.WriteFile("gls.log", []byte(jsonLivePositions), 0700)
	os.WriteFile("gls2.log", []byte(jsonLiveStandings), 0700)

	return C.CString(jsonLiveStandings), 0
}

func main() {}
