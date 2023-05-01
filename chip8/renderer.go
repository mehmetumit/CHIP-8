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
const BORDER_PADDING = 10

type Display [WIDTH][HEIGHT]uint8

var DisplayScale int32

var (
	Window             *sdl.Window
	Renderer           *sdl.Renderer
	WindowWidth        int32 = 940
	WindowHeight       int32 = 570
	PixelColor               = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	BackgroundColor          = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	DisplayBorderColor       = sdl.Color{R: 0, G: 100, B: 100, A: 255}
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
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE)
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
	Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
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
	DrawDisplayBorder()
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
	pixelRect := sdl.Rect{
		X: DISPLAY_PADDING + int32(x*DisplayScale),
		Y: DISPLAY_PADDING + int32(y*DisplayScale),
		W: int32(DisplayScale),
		H: int32(DisplayScale),
	}
	Renderer.FillRect(&pixelRect)
}
func DrawDisplayBorder() {
	setDrawColor(&BackgroundColor)
	setDrawColor(&DisplayBorderColor)
	borderRect := sdl.Rect{
		X: DISPLAY_PADDING - BORDER_PADDING,
		Y: DISPLAY_PADDING - BORDER_PADDING,
		W: int32(WIDTH*DisplayScale) + 2*BORDER_PADDING,
		H: int32(HEIGHT*DisplayScale) + 2*BORDER_PADDING,
	}
	Renderer.FillRect(&borderRect)

}
