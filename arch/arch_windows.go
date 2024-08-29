// Package arch specific options
package arch

import (
	"syscall"

	"github.com/lxn/win"
)

func WindowOptions() {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr("iRacing Championship Standings"))
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)
}
