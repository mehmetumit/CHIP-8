package chip8

import (
	"errors"
	"log"
	"math/rand"
)

// 16 Register each of them 8 bit
// 0x00-0xFF
type Registers [16]Register
type Register uint8

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

// 0x0000
type OpCode uint16

type CPU struct {
	//V0-VF
	Registers      Registers
	IndexRegister  IndexRegister
	Memory         Memory
	ProgramCounter ProgramCounter
	ProgramStack   ProgramStack
	StackPointer   StackPointer
	OpCode         OpCode
}
type Chip8 struct {
	Cpu        CPU
	Display    Display
	KeyPad     KeyPad
	DelayTimer DelayTimer
	SoundTimer SoundTimer
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
func OP_00E0() {
	ClearDisplay()
}

// Return from subroutine
func OP_00EE() {
	chip8.Cpu.StackPointer -= 1
	//Return
	chip8.Cpu.ProgramCounter = chip8.Cpu.ProgramStack[chip8.Cpu.StackPointer]
}

// Jump to address NNN
func OP_1NNN() {
	chip8.Cpu.ProgramCounter = ProgramCounter(chip8.Cpu.OpCode & 0x0FFF)
}

// Call subroutine at NNN
func OP_2NNN() {
	address := (chip8.Cpu.OpCode & 0x0FFF)
	//Save state in stack
	chip8.Cpu.ProgramStack[chip8.Cpu.StackPointer] = chip8.Cpu.ProgramCounter
	chip8.Cpu.StackPointer += 1

	//Call
	chip8.Cpu.ProgramCounter = ProgramCounter(address)
}

// Skip the next instruction if VX equals NN (usually the next instruction is a jump to skip a code block)
func OP_3XNN() {
	registerIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	val := Register(chip8.Cpu.OpCode & 0x00FF)
	if chip8.Cpu.Registers[registerIndex] == val {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Skip the next instruction if VX does not equal NN (usually the next instruction is a jump to skip a code block).
func OP_4XNN() {
	registerIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	val := Register(chip8.Cpu.OpCode & 0x00FF)
	if chip8.Cpu.Registers[registerIndex] != val {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Skip the next instruction if VX equals VY (usually the next instruction is a jump to skip a code block)
func OP_5XY0() {
	registerIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	val := Register(chip8.Cpu.OpCode & 0x00FF)
	if chip8.Cpu.Registers[registerIndex] == val {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Set VX to NN
func OP_6XNN() {
	registerIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	val := Register(chip8.Cpu.OpCode & 0x00FF)
	chip8.Cpu.Registers[registerIndex] = val
}

// Add NN to VX (carry flag is not changed)
func OP_7XNN() {
	registerIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	val := Register(chip8.Cpu.OpCode & 0x00FF)
	chip8.Cpu.Registers[registerIndex] += val
}

// Set VX to the value of VY
func OP_8XY0() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] = chip8.Cpu.Registers[regYIndex]
}

// Set VX to VX or VY. (bitwise OR operation)
func OP_8XY1() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] |= chip8.Cpu.Registers[regYIndex]
}

// Set VX to VX and VY. (bitwise AND operation)
func OP_8XY2() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] &= chip8.Cpu.Registers[regYIndex]
}

// Set VX to VX xor VY
func OP_8XY3() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] ^= chip8.Cpu.Registers[regYIndex]
}

// Add VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not
func OP_8XY4() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] ^= chip8.Cpu.Registers[regYIndex]
}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not
func OP_8XY5() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	//Negative result, set last register to 0
	if chip8.Cpu.Registers[regXIndex] < chip8.Cpu.Registers[regYIndex] {
		chip8.Cpu.Registers[0xF] = 0
	} else {
		chip8.Cpu.Registers[0xF] = 1
	}
	chip8.Cpu.Registers[regXIndex] -= chip8.Cpu.Registers[regYIndex]
}

// Store the least significant bit of VX in VF and then shifts VX to the right by 1
// Ignore VY like CHIP-48 and SCHIP implementations
func OP_8XY6() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	lsb := Register(chip8.Cpu.OpCode & 1)
	chip8.Cpu.Registers[regXIndex] >>= 1
	chip8.Cpu.Registers[0xF] = lsb
}

// Set VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not
func OP_8XY7() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	//Negative result, set last register to 0
	if chip8.Cpu.Registers[regYIndex] < chip8.Cpu.Registers[regXIndex] {
		chip8.Cpu.Registers[0xF] = 0
	} else {
		chip8.Cpu.Registers[0xF] = 1
	}
	chip8.Cpu.Registers[regXIndex] = chip8.Cpu.Registers[regYIndex] - chip8.Cpu.Registers[regXIndex]

}

// Store the most significant bit of VX in VF and then shifts VX to the left by 1
func OP_8XYE() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8

	chip8.Cpu.Registers[0xF] = chip8.Cpu.Registers[regXIndex] & 0xF0
	chip8.Cpu.Registers[regXIndex] <<= 1
}

// Skip the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
func OP_9XY0() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	if chip8.Cpu.Registers[regXIndex] != chip8.Cpu.Registers[regYIndex] {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Set I to the address NNN
func OP_ANNN() {
	chip8.Cpu.IndexRegister = IndexRegister(chip8.Cpu.OpCode & 0x0FFF)
}

// Jump to the address NNN plus V0
func OP_BNNN() {
	address := chip8.Cpu.OpCode & 0x0FFF
	chip8.Cpu.ProgramCounter = ProgramCounter(address + OpCode(chip8.Cpu.Registers[0]))
}

// Set VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
func OP_CXNN() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	chip8.Cpu.Registers[regXIndex] = Register((chip8.Cpu.OpCode & 0x00FF) & OpCode(rand.Intn(256)))
}

// Draw a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels
// TODO
func OP_DXYN() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.OpCode & 0x00F0) >> 4
	pixelNum := chip8.Cpu.OpCode & 0x000F
	//TODO
	Draw(uint8(chip8.Cpu.Registers[regXIndex]), uint8(chip8.Cpu.Registers[regYIndex]), uint8(pixelNum))
}

// Skip the next instruction if the key stored in VX is pressed (usually the next instruction is a jump to skip a code block)
func OP_EX9E() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	key := uint8(chip8.Cpu.Registers[regXIndex])
	//Key pressed
	if chip8.KeyPad[key] {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Skip the next instruction if the key stored in VX is not pressed (usually the next instruction is a jump to skip a code block)
func OP_EXA1() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	key := uint8(chip8.Cpu.Registers[regXIndex])
	//Key not pressed
	if !chip8.KeyPad[key] {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Set VX to the value of the delay timer
func OP_FX07() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	chip8.Cpu.Registers[regXIndex] = Register(chip8.DelayTimer)
}

// A key press is awaited, and then stored in VX (blocking operation, all instruction halted until next key event)
func OP_FX0A() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	//Wait for key press
	waitKeyPress := true
	for i, key := range chip8.KeyPad {
		if key {
			chip8.Cpu.Registers[regXIndex] = Register(i)
			waitKeyPress = false
			break
		}
	}
	//One cycle is 2
	if waitKeyPress {
		chip8.Cpu.ProgramCounter -= 2
	}
}

// Set the delay timer to VX
func OP_FX15() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	chip8.DelayTimer = DelayTimer(chip8.Cpu.Registers[regXIndex])
}

// Set the sound timer to VX
func OP_FX18() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	chip8.SoundTimer = SoundTimer(chip8.Cpu.Registers[regXIndex])
}

// Add VX to I. VF is not affected
func OP_FX1E() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	chip8.Cpu.IndexRegister += IndexRegister(chip8.Cpu.Registers[regXIndex])
}

// Set I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
func OP_FX29() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	fontLocation := FONTSET_START_ADDRESS + uint16(chip8.Cpu.Registers[regXIndex]*5)
	chip8.Cpu.IndexRegister = IndexRegister(fontLocation)
}

// Store the binary-coded decimal(BCD) representation of VX,
// with the hundreds digit in memory at location in I,
// the tens digit at location I+1, and the ones digit at location I+2
func OP_FX33() {
	regXIndex := (chip8.Cpu.OpCode & 0x0F00) >> 8
	num := uint8(chip8.Cpu.Registers[regXIndex])
	//255 -> 0010 0101 0101
	//Ones
	chip8.Cpu.Memory[chip8.Cpu.IndexRegister+2] = num % 10
	num /= 10
	//Tens
	chip8.Cpu.Memory[chip8.Cpu.IndexRegister+1] = num % 10
	num /= 10
	//Hundreds
	chip8.Cpu.Memory[chip8.Cpu.IndexRegister] = num % 10
	num /= 10
}

// Store from V0 to VX (including VX) in memory, starting at address I
// The offset from I is increased by 1 for each value written, but I itself is left unmodified
func OP_FX55() {
	startAddress := chip8.Cpu.IndexRegister
	regXIndex := uint8((chip8.Cpu.OpCode & 0x0F00) >> 8)
	for i := uint8(0); i <= regXIndex; i++ {
		chip8.Cpu.Memory[startAddress+IndexRegister(i)] = uint8(chip8.Cpu.Registers[i])
	}
}

// Fill from V0 to VX (including VX) with values from memory, starting at address I.
// The offset from I is increased by 1 for each value read, but I itself is left unmodified
func OP_FX65() {
	startAddress := chip8.Cpu.IndexRegister
	regXIndex := uint8((chip8.Cpu.OpCode & 0x0F00) >> 8)
	for i := uint8(0); i <= regXIndex; i++ {
		chip8.Cpu.Registers[i] = Register(chip8.Cpu.Memory[startAddress+IndexRegister(i)])
	}
}

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
			//Push rom into memory
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
func Boot(romPath string, displayScale uint8, speed uint8) {
	err := loadRom(romPath)
	if err != nil {
		halt(err)
	}
	loadFonts()
	chip8.Cpu.ProgramCounter = ProgramCounter(START_ADDRESS)
	chip8.Cpu.OpCode = OpCode(0)
	DisplayScale = displayScale
	loop()
}
func loop() {
	for true {
		cycle()
	}
}
func cycle() {
	fetch()
	chip8.Cpu.ProgramCounter += 2
	decodeAndExecute()
	if chip8.DelayTimer > 0 {
		chip8.DelayTimer -= 1
	}
	if chip8.SoundTimer > 0 {
		chip8.SoundTimer -= 1
	}
}
