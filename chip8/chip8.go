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
	OpCode         OpCode
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
var instractionSetOperations = map[OpCode](*func()){}

// Clear the display
func OP_00E0() {}

// Return from subroutine
func OP_00EE() {}

// Jumps to address NNN
func OP_1NNN() {}

// Calls subroutine at NNN
func OP_2NNN() {}

// Skip the next instruction if VX equals NN (usually the next instruction is a jump to skip a code block)
func OP_3XNN() {}

// Skip the next instruction if VX does not equal NN (usually the next instruction is a jump to skip a code block).
func OP_4XNN() {}

// Skips the next instruction if VX equals VY (usually the next instruction is a jump to skip a code block)
func OP_5XY0() {}

// Sets VX to NN
func OP_6XNN() {}

// Adds NN to VX (carry flag is not changed)
func OP_7XNN() {}

// Sets VX to the value of VY
func OP_8XY0() {}

// Sets VX to VX or VY. (bitwise OR operation)
func OP_8XY1() {}

// Sets VX to VX and VY. (bitwise AND operation)
func OP_8XY2() {}

// Sets VX to VX xor VY
func OP_8XY3() {}

// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not
func OP_8XY4() {}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not
func OP_8XY5() {}

// Stores the least significant bit of VX in VF and then shifts VX to the right by 1
func OP_8XY6() {}

// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not
func OP_8XY7() {}

// Stores the most significant bit of VX in VF and then shifts VX to the left by 1
func OP_8XYE() {}

// Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
func OP_9XY0() {}

// Sets I to the address NNN
func OP_ANNN() {}

// Jumps to the address NNN plus V0
func OP_BNNN() {}

// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
func OP_CXNN() {}

// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels
func OP_DXYN() {}

// Skips the next instruction if the key stored in VX is pressed (usually the next instruction is a jump to skip a code block)
func OP_EX9E() {}

// Skips the next instruction if the key stored in VX is not pressed (usually the next instruction is a jump to skip a code block)
func OP_EXA1() {}

// Sets VX to the value of the delay timer
func OP_FX07() {}

// A key press is awaited, and then stored in VX (blocking operation, all instruction halted until next key event)
func OP_FX0A() {}

// Sets the delay timer to VX
func OP_FX15() {}

// Sets the sound timer to VX
func OP_FX18() {}

// Adds VX to I. VF is not affected
func OP_FX1E() {}

// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
func OP_FX29() {}

// Stores the binary-coded decimal representation of VX,
// with the hundreds digit in memory at location in I,
// the tens digit at location I+1, and the ones digit at location I+2
func OP_FX33() {}

// Stores from V0 to VX (including VX) in memory, starting at address I
// The offset from I is increased by 1 for each value written, but I itself is left unmodified
func OP_FX55() {}

// Fills from V0 to VX (including VX) with values from memory, starting at address I.
// The offset from I is increased by 1 for each value read, but I itself is left unmodified
func OP_FX65() {}

func fetch() {
	//01010101 00000000 | 00000000 10101010 -> OpCode is 2 byte
	chip8.Cpu.OpCode = OpCode(uint16(chip8.Cpu.Memory[chip8.Cpu.ProgramCounter])<<8 | uint16(chip8.Cpu.Memory[chip8.Cpu.ProgramCounter+1]))
}
func decodeAndExecute() {
	//chip8.Cpu.OpCode
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
func initSystem() {
	loadFonts()
	chip8.Cpu.ProgramCounter = ProgramCounter(START_ADDRESS)
	chip8.Cpu.OpCode = OpCode(0)
}
func Boot(romPath string) {
	err := loadRom(romPath)
	if err != nil {
		halt(err)
	}
	initSystem()
	loop()
}
func loop() {
	for true {
		fetch()
		decodeAndExecute()
	}
}
