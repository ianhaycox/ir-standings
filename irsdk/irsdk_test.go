package irsdk

import "testing"

func TestInit(t *testing.T) {
	sdk := Init(nil)
	sdk.Close()
}
