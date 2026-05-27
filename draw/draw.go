package draw

import (
	"fmt"
	"log"

	"github.com/Cameliuu/veil/win32"
	"golang.org/x/sys/windows"
)

type Color struct {
	R, G, B uint8
}

func (c Color) ToBGR() uintptr {
	return uintptr(c.R) | uintptr(c.G)<<8 | uintptr(c.B)<<16
}

const (
	DT_TOP             uint32 = 0x00000000
	DT_LEFT            uint32 = 0x00000000
	DT_CENTER          uint32 = 0x00000001
	DT_RIGHT           uint32 = 0x00000002
	DT_VCENTER         uint32 = 0x00000004
	DT_BOTTOM          uint32 = 0x00000008
	DT_WORDBREAK       uint32 = 0x00000010
	DT_SINGLELINE      uint32 = 0x00000020
	DT_EXPANDTABS      uint32 = 0x00000040
	DT_TABSTOP         uint32 = 0x00000080
	DT_NOCLIP          uint32 = 0x00000100
	DT_EXTERNALLEADING uint32 = 0x00000200
	DT_CALCRECT        uint32 = 0x00000400
	DT_NOPREFIX        uint32 = 0x00000800
	DT_INTERNAL        uint32 = 0x00001000
	DT_EDITCONTROL     uint32 = 0x00002000

	DT_PATH_ELLIPSIS        uint32 = 0x00004000
	DT_END_ELLIPSIS         uint32 = 0x00008000
	DT_MODIFYSTRING         uint32 = 0x00010000
	DT_RTLREADING           uint32 = 0x00020000
	DT_WORD_ELLIPSIS        uint32 = 0x00040000
	DT_NOFULLWIDTHCHARBREAK uint32 = 0x00080000
	DT_HIDEPREFIX           uint32 = 0x00100000
	DT_PREFIXONLY           uint32 = 0x00200000
)

var (
	Red   = Color{255, 0, 0}
	Green = Color{0, 255, 0}
	Blue  = Color{0, 0, 255}
	White = Color{255, 255, 255}
	Black = Color{0, 0, 0}
)

func Box(hdc uintptr, rect win32.Rect, c Color) {
	// set outline colour
	pen := win32.CreatePen(c.ToBGR(), 2)
	defer win32.DeleteObject(pen)
	oldPen := win32.SelectObject(hdc, pen)
	defer win32.SelectObject(hdc, oldPen)

	// null brush = no fill, player visible through box
	oldBrush := win32.SelectObject(hdc, win32.GetNullBrush())
	defer win32.SelectObject(hdc, oldBrush)

	if err := win32.Rectangle(hdc, rect); err != nil {
		log.Printf("Could not draw box: %v", err)
	}
}

func TextOut(hdc uintptr, text string, x, y int32, Color Color) {
	err := win32.SetBkMode(hdc, win32.BK_TRANSPARENT)

	if err != nil {
		fmt.Printf("veil: could not set background mode %v", err)
	}

	err = win32.SetTextColor(hdc, Color.ToBGR())
	if err != nil {
		fmt.Printf("veil: could not set text color %v", err)
	}
	textPtr, _ := windows.UTF16PtrFromString(text)
	win32.TextOut(hdc,
		textPtr,
		uint32(len(text)),
		x,
		y)

}
func Box3D(hdc uintptr, rect win32.Rect, c Color, depth int) {
	pen := win32.CreatePen(c.ToBGR(), 2)
	defer win32.DeleteObject(pen)

	oldPen := win32.SelectObject(hdc, pen)
	defer win32.SelectObject(hdc, oldPen)

	x1, y1 := rect.Left, rect.Top
	x2, y2 := rect.Right, rect.Bottom

	dx := int32(depth)
	dy := int32(depth)

	bx1, by1 := x1-dx, y1-dy
	bx2, by2 := x2-dx, y2-dy

	drawLine := func(xa, ya, xb, yb int32) {
		win32.MoveToEx(hdc, win32.Point{X: xa, Y: ya})
		win32.LineTo(hdc, win32.Point{X: xb, Y: yb})
	}

	drawLine(x1, y1, x2, y1)
	drawLine(x2, y1, x2, y2)
	drawLine(x2, y2, x1, y2)
	drawLine(x1, y2, x1, y1)

	drawLine(bx1, by1, bx2, by1)
	drawLine(bx2, by1, bx2, by2)
	drawLine(bx2, by2, bx1, by2)
	drawLine(bx1, by2, bx1, by1)

	drawLine(x1, y1, bx1, by1)
	drawLine(x2, y1, bx2, by1)
	drawLine(x2, y2, bx2, by2)
	drawLine(x1, y2, bx1, by2)
}

func FilledBox(hdc uintptr, rect win32.Rect, c Color) {
	brush := win32.CreateSolidBrush(c.ToBGR())
	defer win32.DeleteObject(brush)
	oldBrush := win32.SelectObject(hdc, brush)
	defer win32.SelectObject(hdc, oldBrush)

	nullPen := win32.GetNullPen() // NULL_PEN = 8
	oldPen := win32.SelectObject(hdc, nullPen)
	defer win32.SelectObject(hdc, oldPen)

	if err := win32.Rectangle(hdc, rect); err != nil {
		log.Printf("Could not draw filled box: %v", err)
	}
}
func healthColor(hp int) Color {
	switch {
	case hp > 60:
		return Color{0, 255, 0}
	case hp > 30:
		return Color{255, 165, 0}
	default:
		return Color{255, 0, 0}
	}
}
func HealthBar(hdc uintptr, x, y, w, h, hp int) {
	// dark background
	FilledBox(hdc, win32.Rect{
		Left:   int32(x),
		Top:    int32(y),
		Right:  int32(x + w),
		Bottom: int32(y + h),
	}, Color{80, 0, 0})

	// health fill
	filledH := h * hp / 100
	FilledBox(hdc, win32.Rect{
		Left:   int32(x),
		Top:    int32(y + h - filledH), // fills from bottom up
		Right:  int32(x + w),
		Bottom: int32(y + h),
	}, healthColor(hp))
}
