package main

import (
	"fmt"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
)

var frameCount int

func callback(hdc uintptr) {
	frameCount++
	fmt.Println("draw frame:", frameCount)
	draw.Box(hdc, win32.Rect{
		Left:   300,
		Top:    200,
		Right:  60,
		Bottom: 120,
	},
		draw.Red)
}
func main() {

	window.Run("AssaultCube", callback)

}
