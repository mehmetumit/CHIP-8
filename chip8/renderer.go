package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

// Monochrome display, 1 or 0
// width -> 64, height -> 32
// Width and height can be set differenlty on some interpreters
const WIDTH = 64
const HEIGHT = 32
const DISPLAY_PADDING = 90

type Display [WIDTH][HEIGHT]uint8

var DisplayScale int32

var (
	Window          *sdl.Window
	Renderer        *sdl.Renderer
	WindowWidth     int32 = 1024
	WindowHeight    int32 = 768
	PixelColor            = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	BackgroundColor       = sdl.Color{R: 0, G: 0, B: 0, A: 255}
)

func init() {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal("SDL initialization failed!", err)
	}
}
func StartDisplay(display *Display) {
	var err error
	Window, err = sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal("Window creation failed!", err)
	}

	Renderer, err = sdl.CreateRenderer(Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal("Failed to create renderer!", err)
	}
	ClearRenderer(display)

}

func setDrawColor(color *sdl.Color) {
	Renderer.SetDrawColor(color.R, color.B, color.G, color.A)
}
func updateRenderer() {
	Renderer.Present()
}
func ClearRenderer(display *Display) {
	log.Print("Display cleaning...")
	setDrawColor(&BackgroundColor)
	for i := 0; i < len(display); i++ {
		for j := 0; j < len(display[i]); j++ {
			display[i][j] = 0
		}
	}
	updateRenderer()
}
func RenderDisplay(display *Display) {
	setDrawColor(&BackgroundColor)
	Renderer.Clear()
	for j := 0; j < HEIGHT; j++ {
		for i := 0; i < WIDTH; i++ {
			pixelState := uint8ToBool(display[i][j])
			DrawPixel(int32(i), int32(j), pixelState)
		}
	}
	updateRenderer()
}
func uint8ToBool(num uint8) bool {
	if num == 0x1 {
		return true
	} else {
		return false
	}

}
func DrawPixel(x int32, y int32, isPixelOn bool) {
	if isPixelOn {
		setDrawColor(&PixelColor)
	} else {
		setDrawColor(&BackgroundColor)
	}
	pixelRect := sdl.Rect{X: DISPLAY_PADDING + int32(x*DisplayScale), Y: DISPLAY_PADDING + int32(y*DisplayScale), W: int32(DisplayScale), H: int32(DisplayScale)}
	Renderer.FillRect(&pixelRect)
}
