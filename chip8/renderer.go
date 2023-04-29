package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

// Monochrome display, 1 or 0
// width -> 64, height -> 32
const WIDTH = 64
const HEIGHT = 32

type Display [WIDTH * HEIGHT]uint8

var DisplayScale uint8

var (
	Window       *sdl.Window
	Surface      *sdl.Surface
	WindowWidth  int32 = 800
	WindowHeight int32 = 600
	PixelColor         = sdl.Color{R: 255, G: 255, B: 255, A: 255}
)

func init() {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal("SDL initialization failed!", err)
	}
}
func StartDisplay() {
	var err error
	Window, err = sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal("Window creation failed!", err)
	}

	Surface, err = Window.GetSurface()
	if err != nil {
		panic(err)
	}
	Surface.FillRect(nil, 0)

	rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	pixel := sdl.MapRGBA(Surface.Format, PixelColor.R, PixelColor.G, PixelColor.B, PixelColor.A)
	Surface.FillRect(&rect, pixel)
	Window.UpdateSurface()
}

func ClearDisplay() {

}
func Draw(x uint8, y uint8, n uint8) {

}
