package main

import (
	"fmt"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
)

func callback(hdc uintptr) {
	fmt.Println("onPaint called, hdc:", hdc)

	draw.Box(hdc, win32.Rect{
		Left:   300,
		Top:    200,
		Right:  60,
		Bottom: 120,
	},
		draw.Green)
}
func main() {

	window.Run("AssaultCube", callback)

}
