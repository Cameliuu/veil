package win32

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Rect struct {
	Left, Top, Right, Bottom int32
}

type Messages struct {
	WmDestroy uint32
	WmPaint   uint32
	WmTimer   uint32
}

var Msg = Messages{
	WmDestroy: 0x0002,
	WmPaint:   0x000F,
	WmTimer:   0x0113,
}

type WndClassEx struct {
	cbSize     uint32         //the size in bytes of the structure
	style      uint32         // the class styles. Cand be OR'd togheter https://learn.microsoft.com/en-us/windows/win32/winmsg/window-class-styles
	wndproc    uintptr        //pointer to the window procedure
	cbClsExtra int32          // number of extra bytes for the window-class structure - initialized to 0
	cbWndExtra int32          // number of extra bytes for the window procredure - also initialzed to 0
	instance   windows.Handle // Handle to the instance that contains the windows procedure
	icon       windows.Handle // Handle to the icon
	cursor     windows.Handle // Handle to the cursor
	brush      windows.Handle // handle to the class background brush
	menuName   *uint16
	className  *uint16        //pointer to a null terminated character stric for the class name of the window
	iconSmall  windows.Handle // pointer to a small incon associated with the window class
}

const className = "veil_overlay_class"

var (
	user32DLL           = windows.NewLazyDLL("user32.DLL")
	procFindWindowW     = user32DLL.NewProc("FindWindowW")
	procGetWindowRect   = user32DLL.NewProc("GetWindowRect")
	procRegisterClassEx = user32DLL.NewProc("RegisterClassExW")
	procCreateWindowExW = user32DLL.NewProc("CreateWindowExW")
	procPostQuitMessage = user32DLL.NewProc("PostQuitMessage")
	procDefWindowProcW  = user32DLL.NewProc("DefWindowProcW")
	procShowWindow      = user32DLL.NewProc("ShowWindow")
	procUpdateWindow    = user32DLL.NewProc("UpdateWindow")
)

func GetRect(hWnd uintptr) Rect {
	var r Rect
	procGetWindowRect.Call(hWnd, uintptr(unsafe.Pointer(&r)))
	return r
}

func RegisterClassEx(wc *WndClassEx) (bool, error) {
	atom, _, err := procRegisterClassEx.Call(uintptr(unsafe.Pointer(wc)))

	if atom == 0 {
		return false, err
	}

	return true, nil
}
func PostQuitMessage() {
	procPostQuitMessage.Call(0)
}

func DefWindowProc(hwnd, msg, wp, lp uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wp, lp)
	return r
}

const (
	wsExLayered     = 0x00080000
	wsExTransparent = 0x00000020
	wsExTopmost     = 0x00000008
	wsExNoActivate  = 0x08000000
	wsPopup         = 0x80000000

	lwaColorKey = 0x00000001
	colorKey    = 0x00FF00FF

	wmDestroy = 0x0002
	wmPaint   = 0x000F
	wmTimer   = 0x0113

	timerID   = 1
	nullBrush = 5
	idcArrow  = 32512
)

func CreateWindowEx(exStyle uint32, className, windowName *uint16, style uint32, x, y, w, h int, instance windows.Handle) (uintptr, error) {

	hwnd, _, _ := procCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(style),
		uintptr(x), uintptr(y),
		uintptr(w), uintptr(h),
		0, // parent
		0, // menu
		uintptr(instance),
		0, // lpParam
	)

	if hwnd == 0 {
		return 0, fmt.Errorf("veil: failed to create window")
	}
	return hwnd, nil
}

func CreateWindow(title string, x, y, w, h int, wndProc uintptr, instanceHandle windows.Handle) (uintptr, error) {

	ptrToClassName, _ := windows.UTF16PtrFromString(className)
	ptrToTitle, _ := windows.UTF16PtrFromString(title)

	//First register the class
	wc := WndClassEx{
		cbSize:    uint32(unsafe.Sizeof(WndClassEx{})),
		wndproc:   wndProc,
		instance:  instanceHandle,
		className: ptrToClassName,
	}

	isWcRegistered, err := RegisterClassEx(&wc)
	if !isWcRegistered {
		return 0, fmt.Errorf("veil: could not register window class: %v", err)
	}

	hWnd, err := CreateWindowEx(
		wsExLayered|
			wsExTransparent|
			wsExNoActivate|
			wsExLayered,
		ptrToClassName,
		ptrToTitle,
		0,
		x,
		y,
		w,
		h,
		instanceHandle)

	if err != nil {
		return 0, err
	}

	return hWnd, nil
}

func FindWindow(title string) uintptr {
	p, _ := windows.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(p)))
	return hwnd
}
func ShowWindow(hwnd uintptr) {
	procShowWindow.Call(hwnd, 5) // 5 = SW_SHOW
}

func UpdateWindow(hwnd uintptr) {
	procUpdateWindow.Call(hwnd)
}
