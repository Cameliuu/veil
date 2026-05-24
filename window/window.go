package window

import (
	"fmt"
	"log"
	"syscall"

	"github.com/Cameliuu/veil/win32"
	"golang.org/x/sys/windows"
)

type Window struct {
	height int
	width  int
	hWnd   windows.HWND
}

var onPaint func(hdc uintptr)

func wndProc(hwnd, msg, wp, lp uintptr) uintptr {
	switch uint32(msg) {
	case win32.WMsg.WmDestroy:
		win32.PostQuitMessage()
		return 0
	case win32.WMsg.WmTimer:
		//invalidate rect - this will mark window as dirty and trigger WmPaint
		win32.InvalidateRect(windows.HWND(hwnd))
		return 0
	case win32.WMsg.WmPaint:
		var ps win32.PaintStruct
		hdc, err := win32.BeginPaint(windows.HWND(hwnd), &ps)
		defer win32.EndPaint(windows.HWND(hwnd), &ps)
		if err != nil {
			log.Printf("veil: Could not get handle to device context :%v", err)
		}

		if onPaint != nil {
			onPaint(hdc)
		}
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wp, lp)
}

func Run(targetTitle string, callback func(hdc uintptr)) {
	onPaint = callback
	window, err := New("AssaultCube")

	if err != nil {
		log.Fatal(err)
	}
	win32.ShowWindow(window.hWnd)
	win32.UpdateWindow(window.hWnd)

	//TO-DO ACTUALLY CONFIG THE FPS
	isTimerSet, err := win32.SetTimer(window.hWnd, 1, 1000/60)
	if !isTimerSet {
		log.Fatalf("veil: could not set timer %v", err)
	}

	var m win32.Msg
	for {
		isMsgSuccessful, err := win32.GetMessage(&m)
		if !isMsgSuccessful {
			log.Fatalf("veil: message cannot be proccesed: %v", err)
			break
		}

		win32.DispatchMessage(&m)
	}
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

	return &Window{
		hWnd:   windows.HWND(hwnd),
		width:  w,
		height: h,
	}, nil
}
