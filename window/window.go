package window

import (
	"fmt"
	"log"
	"syscall"
	"time"

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

		if hdc == 0 {
			return 0
		}
		if err != nil {
			log.Printf("veil: Could not get handle to device context :%v", err)
		}

		clientRect := win32.GetClientRect(uintptr(hwnd))
		brush := win32.CreateSolidBrush(win32.ColorKey)
		win32.FillRect(hdc, &clientRect, brush)
		win32.DeleteObject(brush)
		if onPaint != nil {
			onPaint(hdc)
		}
		return 0
	}
	return win32.DefWindowProc(hwnd, msg, wp, lp)
}

func Run(targetTitle string, callback func(hdc uintptr)) {
	onPaint = callback
	window, err := New(targetTitle)

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
		if win32.PeekMessage(&m) {
			if m.Message == 0x0012 { // WM_QUIT
				break
			}
			win32.DispatchMessage(&m)
		} else {
			time.Sleep(1 * time.Millisecond)
		}
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
	win32.SetWindowPos(hwnd, win32.HwndTopmost, 0, 0, 0, 0, win32.SwpNoSize|win32.SwpNoMove|win32.SwpNoActivate)
	win32.SetLayeredWindowAttributes(hwnd, win32.ColorKey, 0, win32.LwaColorKey)

	return &Window{
		hWnd:   windows.HWND(hwnd),
		width:  w,
		height: h,
	}, nil
}
