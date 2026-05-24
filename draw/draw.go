package draw

import (
	"log"

	"github.com/Cameliuu/veil/win32"
)

type Color struct {
	R, G, B uint8
}

func (c Color) ToBGR() uintptr {
	return uintptr(c.B) | uintptr(c.G)<<8 | uintptr(c.R)<<16
}

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
