package chip8

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"
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
type Opcode uint16

type CPU struct {
	//V0-VF
	Registers      Registers
	IndexRegister  IndexRegister
	Memory         Memory
	ProgramCounter ProgramCounter
	ProgramStack   ProgramStack
	StackPointer   StackPointer
	Opcode         Opcode
}
type Chip8 struct {
	Cpu        CPU
	Display    Display
	Keypad     Keypad
	DelayTimer DelayTimer
	SoundTimer SoundTimer
	Speed      uint8
}

const START_ADDRESS = uint16(0x200)
const FONTSET_START_ADDRESS = uint16(0x050)
const FONTSET_END_ADDRESS = uint16(0x0A0)

var chip8 = &Chip8{
	Cpu: CPU{
		ProgramCounter: ProgramCounter(START_ADDRESS),
	},
}

/*
****Opcodes****
00E0 00EE (0NNN)-> not necessary for most roms
1NNN
2NNN
3XNN
4XNN
5XY0
6XNN
7XNN
8XY0 8XY1 8XY2 8XY3 8XY4 8XY5 8XY6 8XY7 8XYE
9XY0
ANNN
BNNN
CXNN
DXYN
EX9E EXA1
FX07 FX0A FX15 FX18 FX1E FX29 FX33 FX55 FX65
*/
var opcodeTable = map[Opcode](*func()){}

// Clear the display
func OP_00E0() {
	ClearRenderer(&chip8.Display)
}

// Return from subroutine
func OP_00EE() {
	chip8.Cpu.StackPointer -= 1
	//Return
	chip8.Cpu.ProgramCounter = chip8.Cpu.ProgramStack[chip8.Cpu.StackPointer]
}

// Jump to address NNN
func OP_1NNN() {
	chip8.Cpu.ProgramCounter = ProgramCounter(chip8.Cpu.Opcode & 0x0FFF)
}

// Call subroutine at NNN
func OP_2NNN() {
	address := (chip8.Cpu.Opcode & 0x0FFF)
	//Save state in stack
	chip8.Cpu.ProgramStack[chip8.Cpu.StackPointer] = chip8.Cpu.ProgramCounter
	chip8.Cpu.StackPointer += 1

	//Call
	chip8.Cpu.ProgramCounter = ProgramCounter(address)
}

// Skip the next instruction if VX equals NN (usually the next instruction is a jump to skip a code block)
func OP_3XNN() {
	registerIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	val := Register(chip8.Cpu.Opcode & 0x00FF)
	if chip8.Cpu.Registers[registerIndex] == val {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Skip the next instruction if VX does not equal NN (usually the next instruction is a jump to skip a code block).
func OP_4XNN() {
	registerIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	val := Register(chip8.Cpu.Opcode & 0x00FF)
	if chip8.Cpu.Registers[registerIndex] != val {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Skip the next instruction if VX equals VY (usually the next instruction is a jump to skip a code block)
func OP_5XY0() {
	registerIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	val := Register(chip8.Cpu.Opcode & 0x00FF)
	if chip8.Cpu.Registers[registerIndex] == val {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Set VX to NN
func OP_6XNN() {
	registerIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	val := Register(chip8.Cpu.Opcode & 0x00FF)
	chip8.Cpu.Registers[registerIndex] = val
}

// Add NN to VX (carry flag is not changed)
func OP_7XNN() {
	registerIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	val := Register(chip8.Cpu.Opcode & 0x00FF)
	chip8.Cpu.Registers[registerIndex] += val
}

// Set VX to the value of VY
func OP_8XY0() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] = chip8.Cpu.Registers[regYIndex]
}

// Set VX to VX or VY. (bitwise OR operation)
func OP_8XY1() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] |= chip8.Cpu.Registers[regYIndex]
}

// Set VX to VX and VY. (bitwise AND operation)
func OP_8XY2() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] &= chip8.Cpu.Registers[regYIndex]
}

// Set VX to VX xor VY
func OP_8XY3() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	chip8.Cpu.Registers[regXIndex] ^= chip8.Cpu.Registers[regYIndex]
}

// Add VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not
func OP_8XY4() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	newRegX := chip8.Cpu.Registers[regXIndex] + chip8.Cpu.Registers[regYIndex]
	//Overflow detection
	if newRegX < chip8.Cpu.Registers[regXIndex] {
		chip8.Cpu.Registers[0xF] = 1
	} else {
		chip8.Cpu.Registers[0xF] = 0
	}
	//Set register X
	chip8.Cpu.Registers[regXIndex] = newRegX
}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not
func OP_8XY5() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
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
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	lsb := Register(chip8.Cpu.Opcode & 1)
	chip8.Cpu.Registers[regXIndex] >>= 1
	chip8.Cpu.Registers[0xF] = lsb
}

// Set VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not
func OP_8XY7() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
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
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8

	chip8.Cpu.Registers[0xF] = chip8.Cpu.Registers[regXIndex] & 0xF0
	chip8.Cpu.Registers[regXIndex] <<= 1
}

// Skip the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
func OP_9XY0() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	if chip8.Cpu.Registers[regXIndex] != chip8.Cpu.Registers[regYIndex] {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Set I to the address NNN
func OP_ANNN() {
	chip8.Cpu.IndexRegister = IndexRegister(chip8.Cpu.Opcode & 0x0FFF)
}

// Jump to the address NNN plus V0
func OP_BNNN() {
	address := chip8.Cpu.Opcode & 0x0FFF
	chip8.Cpu.ProgramCounter = ProgramCounter(address + Opcode(chip8.Cpu.Registers[0]))
}

// Set VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
func OP_CXNN() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	chip8.Cpu.Registers[regXIndex] = Register((chip8.Cpu.Opcode & 0x00FF) & Opcode(rand.Intn(256)))
}

/*
Draw a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height
of N pixels. Each row of 8 pixels is read as bit-coded starting from memory
location I; I value does not change after the execution of this instruction.
Sprite pixels are XOR'd with corresponding screen pixels. The carry flag (VF) is
set to 1 if any screen pixels and sprite pixels are on at the same position otherwise
set to 0. This is used for collision detection.
*/
func OP_DXYN() {
	log.Print("Draw sprite instruction called!")
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	regYIndex := (chip8.Cpu.Opcode & 0x00F0) >> 4
	pixelNum := chip8.Cpu.Opcode & 0x000F
	startAddress := chip8.Cpu.IndexRegister
	posX := uint8(chip8.Cpu.Registers[regXIndex])
	posY := uint8(chip8.Cpu.Registers[regYIndex])
	isCollided := uint8(0)
	chip8.Cpu.Registers[0xF] = 0
	//Iterate over sprite in the memory
	for i := uint8(0); i < uint8(pixelNum); i++ {
		//8 pixels are loaded
		pixelBits := chip8.Cpu.Memory[uint16(startAddress)+uint16(i)]
		log.Printf("Pixel bits: 0x%X", pixelBits)
		for j := uint8(0); j < 8; j++ {
			//Get left most bit
			bit := uint8((pixelBits & 0x80) >> 7)
			//Collision
			if bit == 1 && chip8.Display[(posX+j)%WIDTH][(posY+i)%HEIGHT] == 1 {
				isCollided = 1
			}
			//Limit indicies to prevent overflow
			chip8.Display[(posX+j)%WIDTH][(posY+i)%HEIGHT] ^= bit
			pixelBits = pixelBits << 1
			log.Printf("Bit: 0x%X", bit)
			log.Printf("Pixel bits in loop: 0x%X", pixelBits)
		}
	}
	//Set the flip flag
	chip8.Cpu.Registers[0xF] = Register(isCollided)
}

// Skip the next instruction if the key stored in VX is pressed (usually the next instruction is a jump to skip a code block)
func OP_EX9E() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	key := uint8(chip8.Cpu.Registers[regXIndex])
	//Key pressed
	if chip8.Keypad[key] {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Skip the next instruction if the key stored in VX is not pressed (usually the next instruction is a jump to skip a code block)
func OP_EXA1() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	key := uint8(chip8.Cpu.Registers[regXIndex])
	//Key not pressed
	if !chip8.Keypad[key] {
		chip8.Cpu.ProgramCounter += 2
	}
}

// Set VX to the value of the delay timer
func OP_FX07() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	chip8.Cpu.Registers[regXIndex] = Register(chip8.DelayTimer)
}

// A key press is awaited, and then stored in VX (blocking operation, all instruction halted until next key event)
func OP_FX0A() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	//Wait for key press
	waitKeyPress := true
	for i, key := range chip8.Keypad {
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
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	chip8.DelayTimer = DelayTimer(chip8.Cpu.Registers[regXIndex])
}

// Set the sound timer to VX
func OP_FX18() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	chip8.SoundTimer = SoundTimer(chip8.Cpu.Registers[regXIndex])
	PlayAudio()
}

// Add VX to I. VF is not affected
func OP_FX1E() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	chip8.Cpu.IndexRegister += IndexRegister(chip8.Cpu.Registers[regXIndex])
}

// Set I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
func OP_FX29() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
	fontLocation := FONTSET_START_ADDRESS + uint16(chip8.Cpu.Registers[regXIndex]*5)
	chip8.Cpu.IndexRegister = IndexRegister(fontLocation)
}

/*
Store the binary-coded decimal(BCD) representation of VX,
with the hundreds digit in memory at location in I,
the tens digit at location I+1, and the ones digit at location I+2
*/
func OP_FX33() {
	regXIndex := (chip8.Cpu.Opcode & 0x0F00) >> 8
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
	regXIndex := uint8((chip8.Cpu.Opcode & 0x0F00) >> 8)
	for i := uint8(0); i <= regXIndex; i++ {
		chip8.Cpu.Memory[startAddress+IndexRegister(i)] = uint8(chip8.Cpu.Registers[i])
	}
}

// Fill from V0 to VX (including VX) with values from memory, starting at address I.
// The offset from I is increased by 1 for each value read, but I itself is left unmodified
func OP_FX65() {
	startAddress := chip8.Cpu.IndexRegister
	regXIndex := uint8((chip8.Cpu.Opcode & 0x0F00) >> 8)
	for i := uint8(0); i <= regXIndex; i++ {
		chip8.Cpu.Registers[i] = Register(chip8.Cpu.Memory[startAddress+IndexRegister(i)])
	}
}

func fetch() {
	//01010101 00000000 | 00000000 10101010 -> Opcodes are 2 byte each
	chip8.Cpu.Opcode = Opcode(uint16(chip8.Cpu.Memory[chip8.Cpu.ProgramCounter])<<8 | uint16(chip8.Cpu.Memory[chip8.Cpu.ProgramCounter+1]))
	log.Printf("Fetched Opcode: 0x%X", chip8.Cpu.Opcode)
}
func decodeAndExecute() {
	opcode := chip8.Cpu.Opcode
	firstNum := uint8((opcode & 0xF000) >> 12)
	lastTwoNum := uint8((opcode & 0x00F0) | (opcode & 0x000F))
	lastNum := uint8(opcode & 0x000F)
	switch firstNum {
	case 0x0: // 00E0 00EE
		switch lastNum {
		case 0x0: //00E0
			OP_00E0()
		case 0xE: //00EE
			OP_00EE()
		}
	case 0x1: // 1NNN
		OP_1NNN()
	case 0x2: // 2NNN
		OP_2NNN()
	case 0x3: // 3XNN
		OP_3XNN()
	case 0x4: // 4XNN
		OP_4XNN()
	case 0x5: // 5XY0
		OP_5XY0()
	case 0x6: // 6XNN
		OP_6XNN()
	case 0x7: // 7XNN
		OP_7XNN()
	case 0x8: // 8XY0 8XY1 8XY2 8XY3 8XY4 8XY5 8XY6 8XY7 8XYE
		switch lastNum {
		case 0x0: //8XY0
			OP_8XY0()
		case 0x1: //8XY1
			OP_8XY1()
		case 0x2: //8XY2
			OP_8XY2()
		case 0x3: //8XY3
			OP_8XY3()
		case 0x4: //8XY4
			OP_8XY4()
		case 0x5: //8XY5
			OP_8XY5()
		case 0x6: //8XY6
			OP_8XY6()
		case 0x7: //8XY7
			OP_8XY7()
		case 0xE: //8XYE
			OP_8XYE()
		}
	case 0x9: // 9XY0
		OP_9XY0()
	case 0xA: // ANNN
		OP_ANNN()
	case 0xB: // BNNN
		OP_BNNN()
	case 0xC: // CXNN
		OP_CXNN()
	case 0xD: // DXYN
		OP_DXYN()
	case 0xE: //EX9E EXA1
		switch lastTwoNum {
		case 0x9E: // EX9E
			OP_EX9E()
		case 0xA1: //  EXA1
			OP_EXA1()
		}
	case 0xF: // FX07 FX0A FX15 FX18 FX1E FX29 FX33 FX55 FX65
		switch lastTwoNum {
		case 0x07: //FX07
			OP_FX07()
		case 0x0A: //FX0A
			OP_FX0A()
		case 0x15: //FX15
			OP_FX15()
		case 0x18: //FX18
			OP_FX18()
		case 0x1E: //FX1E
			OP_FX1E()
		case 0x29: //FX29
			OP_FX29()
		case 0x33: //FX33
			OP_FX33()
		case 0x55: //FX55
			OP_FX55()
		case 0x65: //FX65
			OP_FX65()

		}

	}
}
func halt() {
	log.Print("Halting...")
	os.Exit(1)
}

func checkRomSize(romData *[]byte) error {
	log.Println("Rom size:", len(*romData), "byte")
	if int(START_ADDRESS)+len(chip8.Cpu.Memory)+len(*romData) < 0 {
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
func Boot(romPath string, displayScale int32, speed uint8) {
	log.SetFlags(4)
	err := loadRom(romPath)
	if err != nil {
		panic(err)
	}
	loadFonts()
	chip8.Cpu.ProgramCounter = ProgramCounter(START_ADDRESS)
	chip8.Cpu.Opcode = Opcode(0)
	chip8.Speed = speed
	DisplayScale = displayScale
	StartDisplay(&chip8.Display)
	loop()
}
func loop() {
	start := time.Now()
	for true {
		if time.Since(start).Milliseconds() >= int64(chip8.Speed) {
			log.Print("Running...")
			cycle()
			start = time.Now()
			EventHandler(halt, &chip8.Keypad)
		}
	}
}
func cycle() {
	log.Println("Delay Timer:", chip8.DelayTimer)
	log.Println("Sound Timer:", chip8.SoundTimer)
	if chip8.DelayTimer > 0 {
		chip8.DelayTimer -= 1
	}
	if chip8.SoundTimer > 0 {
		chip8.SoundTimer -= 1
		if chip8.SoundTimer <= 0{
			PauseAudio()
		}
	}
	fetch()
	if chip8.Cpu.ProgramCounter+2 < ProgramCounter(len(chip8.Cpu.Memory)) {
		chip8.Cpu.ProgramCounter += 2
	} else {
		log.Fatal("Reached to end of memory!!!")
	}
	decodeAndExecute()
	RenderDisplay(&chip8.Display)
}
