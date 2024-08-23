package irsdk

import "testing"

func TestInit(t *testing.T) {
	t.Skip()

	sdk := Init(nil)
	sdk.Close()
}
