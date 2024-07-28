package main

import (
	"C"
)

//export Results
func Results(livePositions string) (*C.char, int) {
	return C.CString("Hello world" + livePositions), 0
}

func main() {}
