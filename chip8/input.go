package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

// 0x0 to 0xF -> store pressed or not
type KeyPad [16]bool

var keyMap = map[uint8]uint8{
	'1': 0,
	'2': 1,
	'3': 2,
	'4': 3,
	'Q': 4,
	'W': 5,
	'E': 6,
	'R': 7,
	'A': 8,
	'S': 9,
	'D': 10,
	'F': 11,
	'Z': 12,
	'X': 13,
	'C': 14,
	'V': 15,
}

func EventHandler(quitEvent func(), keyPad *KeyPad) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			log.Print("Quit Event Handled")
			Window.Destroy()
			sdl.Quit()
			quitEvent()
			break
		}
	}
}
