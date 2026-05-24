package main

import (
	"fmt"

	"github.com/Cameliuu/veil/window"
)

func callback(hdc uintptr) {
	fmt.Printf("HDC Address: 0x%x\n", hdc)
}
func main() {

	window.Run("veil", callback)

}
