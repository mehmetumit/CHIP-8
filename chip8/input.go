package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

// 0x0 to 0xF -> store pressed or not
type Keypad [16]bool

var keyMap = map[uint8]uint8{
	'1': 0x1,
	'2': 0x2,
	'3': 0x3,
	'4': 0xC,
	'q': 0x4,
	'w': 0x5,
	'e': 0x6,
	'r': 0xD,
	'a': 0x7,
	's': 0x8,
	'd': 0x9,
	'f': 0xE,
	'z': 0xA,
	'x': 0x0,
	'c': 0xB,
	'v': 0xF,
}

func EventHandler(quitEvent func(), keyPad *Keypad) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t:= event.(type) {
		case *sdl.QuitEvent:
			log.Print("Quit Event Handled")
			Window.Destroy()
			ClouseAudio()
			sdl.Quit()
			quitEvent()
			break
		case *sdl.KeyboardEvent:
			handleKeys(t.Keysym.Sym, keyPad, t.State)
		}
	}
}
func handleKeys(keyCode sdl.Keycode, keyPad *Keypad, state uint8){
	if keyIndex, isExists := keyMap[uint8(keyCode)]; isExists{
	log.Println("Key:", keyCode, "Index:",keyIndex)
		if state == sdl.PRESSED{
			keyPad[keyIndex] = true
		}else if state == sdl.RELEASED{
			keyPad[keyIndex] = false
		}
	}

}
