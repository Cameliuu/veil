package window

import (
	"fmt"
	"syscall"

	"github.com/Cameliuu/veil/win32"
	"golang.org/x/sys/windows"
)

type Window struct {
	height int
	width  int
	hWnd   windows.HWND
}

func wndProc(hwnd, msg, wp, lp uintptr) uintptr {
	switch uint32(msg) {
	case win32.Msg.WmDestroy:
		win32.PostQuitMessage()
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wp, lp)
}

func New(targetTitle string) (*Window, error) {
	gameWindow := win32.FindWindow(targetTitle)
	if gameWindow == 0 {
		return nil, fmt.Errorf("Could not find %q window", targetTitle)
	}

	rect := win32.GetRect(gameWindow)
	w := int(rect.Right - rect.Left)
	h := int(rect.Bottom - rect.Top)

	var instanceHandle windows.Handle
	windows.GetModuleHandleEx(0, nil, &instanceHandle)

	hwnd, err := win32.CreateWindow(
		"",
		int(rect.Left), int(rect.Top),
		w, h,
		syscall.NewCallback(wndProc),
		instanceHandle,
	)
	if err != nil {
		return nil, err
	}

	win32.ShowWindow(hwnd)
	win32.UpdateWindow(hwnd)

	return &Window{
		hWnd:   windows.HWND(hwnd),
		width:  w,
		height: h,
	}, nil
}
