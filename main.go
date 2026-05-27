package main

import (
	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
)

var frameCount int

func callback(hdc uintptr) {
	frameCount++
	text := string("sadasdsa")
	draw.TextOut(hdc, text, 500, 500)
	draw.Box(hdc, win32.Rect{
		Left:   300,
		Top:    200,
		Right:  60,
		Bottom: 120,
	},
		draw.Red)
}
func main() {

	window.Run("Counter-Strike", callback)

}
