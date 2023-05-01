package main

import (
	"flag"
	"fmt"

	"github.com/mehmetumit/CHIP-8/chip8"
)

func main() {
	var romPath string
	var displayScale int
	var speed uint
	flag.StringVar(&romPath, "path", "./roms/Instruction-Test.ch8", "The file path of rom")
	flag.UintVar(&speed, "speed", 3, "The emulation speed")
	flag.IntVar(&displayScale, "scale", 12, "The display scale")

	flag.Parse()
	fmt.Println(displayScale)
	fmt.Println(speed)
	chip8.Boot(romPath, int32(displayScale), uint8(speed))
}
