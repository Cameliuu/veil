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
