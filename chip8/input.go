package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)
type KeyMap map[uint8]uint8
//0x0 to 0xF -> store pressed or not
type KeyPad [16]bool

// Map default layout to custom
var keyMap = KeyMap{
	'1': '1',
	'2': '2',
	'3': '3',
	'C': '4',
	'4': 'Q',
	'5': 'W',
	'6': 'E',
	'D': 'R',
	'7': 'A',
	'8': 'S',
	'9': 'D',
	'E': 'F',
	'A': 'Z',
	'0': 'X',
	'B': 'C',
	'F': 'V',
}

func GetKeymap() KeyMap {
	return keyMap
}
func EventHandler(quitEvent func()) {
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
