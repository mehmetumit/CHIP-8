package main

import (
	"flag"
	"github.com/mehmetumit/CHIP-8/chip8"
)

func main() {
	var romPath string
	var displayScale int32
	var speed uint8
	flag.StringVar(&romPath, "path", "./roms/Instruction-Test.ch8", "The file path of rom")
	speed = uint8(*flag.Uint("speed", 3, "The emulation speed"))
	displayScale = int32(*flag.Int("scale", 12, "The display scale"))
	flag.Parse()
	chip8.Boot(romPath, displayScale, speed)
}
