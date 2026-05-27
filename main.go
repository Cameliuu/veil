package main

import (
	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
)

var frameCount int

func callback(hdc uintptr) {
	frameCount++

	draw.Box3D(hdc, win32.Rect{
		Left:   60,
		Top:    120,
		Right:  300,
		Bottom: 200,
	}, draw.Red, 10)
}
func main() {

	window.Run("Counter-Strike", callback)

}
