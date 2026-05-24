package main

import (
	"log"

	"github.com/Cameliuu/veil/window"
)

func main() {
	_, err := window.New("AssaultCube")

	if err != nil {
		log.Fatal(err)
	}

}
