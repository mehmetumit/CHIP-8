package chip8

import (
	"errors"
	"log"
)

// 16 Register each of them 8 bit
// 0x00-0xFF
type Register [16]uint8

// 16 bit, because max memory address(0xFFF) too big for an 8 bit register
type IndexRegister uint16

// 16 bit, memory address of next instruction(8 bit not enough)
type ProgramCounter uint16

// Keep track of execution order
type ProgramStack [16]ProgramCounter

// Similar to program counter but for program stack
type StackPointer uint8
type DelayTimer uint8
type SoundTimer uint8

// 4KB memory
// 0x000-0xFFF -> Address space
// 0x000-0x1FF -> Reserved for CHIP-8 interpreter, not used for now
// 0x050-0x0A0 -> For 16 built-in characters (0 to F)(ROMs will bee looking for these characters)
// 0x200-0xFFF -> Instructions from the ROM. May not be full
type Memory [4 * 1024]uint8
type OpCode uint16

type CPU struct {
	//V0-VF
	Registers      Register
	IndexRegister  IndexRegister
	Memory         Memory
	ProgramCounter ProgramCounter
	ProgramStack   ProgramStack
	StackPointer   StackPointer
	// OpCode         OpCode
}
type Chip8 struct {
	Cpu     CPU
	Display Display
	KeyMap  KeyMap
}

const START_ADDRESS = uint16(0x200)
const FONTSET_START_ADDRESS = uint16(0x050)
const FONTSET_END_ADDRESS = uint16(0x0A0)

var chip8 = &Chip8{
	Cpu: CPU{
		ProgramCounter: ProgramCounter(START_ADDRESS),
	},
}

func fetch() OpCode {
	opCode := OpCode(0xffff)
	return opCode

}
func decode(opCode OpCode) func() {
	return func() {
		log.Print("Decode!")
	}

}
func execute(command func()) {
	command()
}
func halt(e error) {
	panic(e)
}

func checkRomSize(romData *[]byte) error {
	if len(*romData)-int(START_ADDRESS)-len(chip8.Cpu.Memory) < 0 {
		return errors.New("Rom is too large to fit into memory!")
	}
	return nil

}
func loadRom(filePath string) error {
	romData, err := ReadFile(filePath)
	if err == nil {
		err = checkRomSize(&romData)
		if err == nil {
			copy(chip8.Cpu.Memory[START_ADDRESS:], romData)
			log.Printf(`%v rom loaded successfully!`, filePath)
		}
	}
	return err
}
func loadFonts() {
	copy(chip8.Cpu.Memory[FONTSET_START_ADDRESS:FONTSET_END_ADDRESS], Fontset[:])
	log.Print("Fontset loaded successfully!")

}
func Boot(romPath string) {
	err := loadRom(romPath)
	if err != nil {
		halt(err)
	}
	loadFonts()
	loop()
}
func loop() {
	for true {
		opCode := fetch()
		command := decode(opCode)
		execute(command)
	}
}
