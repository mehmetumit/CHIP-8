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
	'q': 4,
	'w': 5,
	'e': 6,
	'r': 7,
	'a': 8,
	's': 9,
	'd': 10,
	'f': 11,
	'z': 12,
	'x': 13,
	'c': 14,
	'v': 15,
}

func EventHandler(quitEvent func(), keyPad *KeyPad) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t:= event.(type) {
		case *sdl.QuitEvent:
			log.Print("Quit Event Handled")
			Window.Destroy()
			sdl.Quit()
			quitEvent()
			break
		case *sdl.KeyboardEvent:
			handleKeys(t.Keysym.Sym, keyPad, t.State)
		}
	}
}
func handleKeys(keyCode sdl.Keycode, keyPad *KeyPad, state uint8){
	log.Println("Key:", keyCode)
	if keyIndex, isExists := keyMap[uint8(keyCode)]; isExists{
		if state == sdl.PRESSED{
			keyPad[keyIndex] = true
		}else if state == sdl.RELEASED{
			keyPad[keyIndex] = false
		}
	}

}
