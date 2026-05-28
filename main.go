package main

import (
	"log"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/window"
)

var frameCount int

func callback(hdc uintptr) {
	frameCount++

	draw.Circle(hdc, 1920/2, 1080/2, 50, draw.Red)
}
func main() {
	w, err := window.New("Counter-Strike")

	if err != nil {
		log.Fatalf("veil: could not create window: %w", err)
	}

	w.Run(callback)
}
