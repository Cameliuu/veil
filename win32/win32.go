package win32

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
=========================================================================== STRUCTS ======================================================
*/
type Point struct{ X, Y int32 }

type Rect struct {
	Left, Top, Right, Bottom int32
}
type Msg struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      Point
}

type Messages struct {
	WmDestroy uint32
	WmPaint   uint32
	WmTimer   uint32
}

var WMsg = Messages{
	WmDestroy: 0x0002,
	WmPaint:   0x000F,
	WmTimer:   0x0113,
}

type Options struct {
	FPS           int
	UseClientRect bool
}
type PaintStruct struct {
	Hdc       uintptr
	Erase     int32
	RcPaint   [4]int32
	Restore   int32
	IncUpdate int32
	Reserved  [32]byte
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

/*
================================================================================= CONSTANTS =================================================
*/
const className = "veil_overlay_class"
const CLR_INVALID = 0xFFFFFFFF
const (
	BK_TRANSPARENT = 1
	BK_OPAQUE      = 2
)

var (
	user32DLL                      = windows.NewLazyDLL("user32.DLL")
	gdiDLL                         = windows.NewLazyDLL("gdi32.dll")
	procFindWindowW                = user32DLL.NewProc("FindWindowW")
	procGetWindowRect              = user32DLL.NewProc("GetWindowRect")
	procGetClientRect              = user32DLL.NewProc("GetClientRect")
	procRegisterClassEx            = user32DLL.NewProc("RegisterClassExW")
	procCreateWindowExW            = user32DLL.NewProc("CreateWindowExW")
	procPostQuitMessage            = user32DLL.NewProc("PostQuitMessage")
	procDefWindowProcW             = user32DLL.NewProc("DefWindowProcW")
	procShowWindow                 = user32DLL.NewProc("ShowWindow")
	procUpdateWindow               = user32DLL.NewProc("UpdateWindow")
	procSetTimer                   = user32DLL.NewProc("SetTimer")
	procGetMessage                 = user32DLL.NewProc("GetMessageW")
	procPeekMessageW               = user32DLL.NewProc("PeekMessageW")
	procDispatchMessageW           = user32DLL.NewProc("DispatchMessageW")
	procInvalidateRect             = user32DLL.NewProc("InvalidateRect")
	procBeginPaint                 = user32DLL.NewProc("BeginPaint")
	procEndPaint                   = user32DLL.NewProc("EndPaint")
	procRectangle                  = gdiDLL.NewProc("Rectangle")
	procFillRect                   = user32DLL.NewProc("FillRect")
	procCreateSolidBrush           = gdiDLL.NewProc("CreateSolidBrush")
	procDeleteObject               = gdiDLL.NewProc("DeleteObject")
	procSelectObject               = gdiDLL.NewProc("SelectObject")
	procGetStockObject             = gdiDLL.NewProc("GetStockObject")
	procSetLayeredWindowAttributes = user32DLL.NewProc("SetLayeredWindowAttributes")
	procSetWindowPos               = user32DLL.NewProc("SetWindowPos")
	procCreatePen                  = gdiDLL.NewProc("CreatePen")
	procDrawText                   = user32DLL.NewProc("DrawTextW")
	procTextOut                    = gdiDLL.NewProc("TextOutW")
	procSetBkMode                  = gdiDLL.NewProc("SetBkMode")
	procSetTextColor               = gdiDLL.NewProc("SetTextColor")
)

func SetTextColor(hdc uintptr, color uintptr) error {
	r, _, err := procSetTextColor.Call(hdc, color)

	if r == CLR_INVALID {
		return err
	}
	return nil
}
func SetBkMode(hdc uintptr, mode uint32) error {
	r, _, err := procSetBkMode.Call(hdc, uintptr(mode))

	if r == 0 {
		return err
	}
	return nil
}

func SetWindowPos(hwnd, hwndInsertAfter uintptr, x, y, w, h int, flags uint32) {
	procSetWindowPos.Call(
		hwnd,
		hwndInsertAfter,
		uintptr(x),
		uintptr(y),
		uintptr(w),
		uintptr(h),
		uintptr(flags),
	)
}
func SetLayeredWindowAttributes(hwnd, colorKey, alpha, flags uintptr) {
	procSetLayeredWindowAttributes.Call(hwnd, colorKey, alpha, flags)
}
func DeleteObject(hObj uintptr) error {
	r, _, err := procDeleteObject.Call(hObj)

	if r == 0 {
		return err
	}
	return nil
}
func CreatePen(color uintptr, width int) uintptr {
	pen, _, _ := procCreatePen.Call(0, uintptr(width), color)
	return pen
}
func TextOut(hdc uintptr, text *uint16, length uint32, x, y int32) {
	textPtr := unsafe.Pointer(text)
	procTextOut.Call(hdc,
		uintptr(x),
		uintptr(y),
		uintptr(textPtr),
		uintptr(length))
}

func SelectObject(hdc, obj uintptr) uintptr {
	old, _, _ := procSelectObject.Call(hdc, obj)
	return old
}
func GetNullPen() uintptr {
	brush, _, _ := procGetStockObject.Call(8)
	return brush
}
func GetNullBrush() uintptr {
	brush, _, _ := procGetStockObject.Call(5)
	return brush
}
func FillRect(hdc uintptr, rect *Rect, brush uintptr) error {

	r, _, err := procFillRect.Call(hdc,
		uintptr(unsafe.Pointer(rect)),
		brush)
	if r == 0 {
		return err
	}

	return nil
}
func CreateSolidBrush(color uintptr) uintptr {
	brush, _, _ := procCreateSolidBrush.Call(color)
	return brush
}
func GetClientRect(hwnd uintptr) Rect {
	var r Rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	return r
}
func Rectangle(hdc uintptr, rect Rect) error {
	r, _, err := procRectangle.Call(hdc,
		uintptr(rect.Left),
		uintptr(rect.Top),
		uintptr(rect.Right),
		uintptr(rect.Bottom))
	if r == 0 {
		return err
	}
	return nil
}
func GetRect(hWnd uintptr) Rect {
	var r Rect
	procGetWindowRect.Call(hWnd, uintptr(unsafe.Pointer(&r)))
	return r
}

func InvalidateRect(hWnd windows.HWND) {
	procInvalidateRect.Call(uintptr(hWnd),
		0,
		1)
}

func BeginPaint(hWnd windows.HWND, paint *PaintStruct) (uintptr, error) {
	hdc, _, err := procBeginPaint.Call(uintptr(hWnd),
		uintptr(unsafe.Pointer(paint)))
	if hdc == 0 {
		return 0, err
	}

	return hdc, nil
}
func DrawText(hdc uintptr, text *uint16, length uint32, rect *Rect, format uint32) error {
	ret, _, err := procDrawText.Call(
		hdc,
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(&rect)),
		uintptr(format),
	)

	if ret == 0 {
		return err
	}

	return nil
}
func EndPaint(hWnd windows.HWND, paint *PaintStruct) {
	procEndPaint.Call(uintptr(hWnd), uintptr(unsafe.Pointer(paint)))
}

func PeekMessage(msg *Msg) bool {
	r, _, _ := procPeekMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		0, 0, 0,
		1,
	)
	return r != 0
}
func GetMessage(msg *Msg) (bool, error) {
	r, _, err := procGetMessage.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)
	if r == ^uintptr(0) { // -1 = real error
		return false, err
	}
	return r != 0, nil // 0 = WM_QUIT, non-zero = normal message
}
func DispatchMessage(msg *Msg) {
	procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
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

func SetTimer(hWnd windows.HWND, idEvent uint32, elapsedTime uint32) (bool, error) {
	isTimerSet, _, err := procSetTimer.Call(
		uintptr(hWnd),
		uintptr(idEvent),
		uintptr(elapsedTime),
	)

	if isTimerSet == 0 {
		return false, err
	}

	return true, nil
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

	LwaColorKey = 0x00000001
	ColorKey    = 0x00FF00FF

	wmDestroy = 0x0002
	wmPaint   = 0x000F
	wmTimer   = 0x0113

	timerID   = 1
	nullBrush = 5
	idcArrow  = 32512
)
const (
	HwndTopmost   = ^uintptr(0)                   // -1, always on top
	HwndNoTopmost = uintptr(18446744073709551614) // -2, removes topmost
	SwpNoSize     = 0x0001
	SwpNoMove     = 0x0002
	SwpNoActivate = 0x0010
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
func ShowWindow(hwnd windows.HWND) {
	procShowWindow.Call(uintptr(hwnd), 5) // 5 = SW_SHOW
}

func UpdateWindow(hwnd windows.HWND) {
	procUpdateWindow.Call(uintptr(hwnd))
}
