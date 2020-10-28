package main

import (
	"fmt"
	"os"
	"time"
)

// Registers of the chip8 CPU
type Registers struct {
	ram [4096]byte
	// RAM address: 0x000 to 0xFFF
	// Interpreter location: 0x000 to 0x1FF
	// Sprite location: 0x000 to 0x050
	V     [16]uint8 // V0 to VF; VF stores flags
	I     uint16
	PC    uint16 // program counter
	SP    uint8  // stack pointer
	stack [16]uint16
	// store the address that the interpreter
	// shoud return to when finished with a subroutine
	DT uint8 // delay time
	ST uint8 // sound timer
}

// Chip8 is an emulator
type Chip8 struct {
	Registers
	display [32]uint64
	// display  [64][32]bool
	// storing the display image (64x32) by row in a uint64
	// should let the array length change to 64x48 or 64x64
	// should support 128x64 too
}

func (chip8 *Chip8) loadDigits() {
	digits := [][5]byte{
		{0xF0, 0x90, 0x90, 0x90, 0xF0},
		{0x20, 0x60, 0x20, 0x20, 0x70},
		{0xF0, 0x10, 0xF0, 0x80, 0xF0},
		{0xF0, 0x10, 0xF0, 0x10, 0xF0},
		{0x90, 0x90, 0xF0, 0x10, 0x10},
		{0xF0, 0x80, 0xF0, 0x10, 0xF0},
		{0xF0, 0x80, 0xF0, 0x90, 0xF0},
		{0xF0, 0x10, 0x20, 0x40, 0x40},
		{0xF0, 0x90, 0xF0, 0x90, 0xF0},
		{0xF0, 0x90, 0xF0, 0x10, 0xF0},
		{0xF0, 0x90, 0xF0, 0x90, 0x90},
		{0xE0, 0x90, 0xE0, 0x90, 0xE0},
		{0xF0, 0x80, 0x80, 0x80, 0xF0},
		{0xE0, 0x90, 0x90, 0x90, 0xE0},
		{0xF0, 0x80, 0xF0, 0x80, 0xF0},
		{0xF0, 0x80, 0xF0, 0x80, 0x80},
	}
	for d := range digits {
		for row := range digits[d] {
			chip8.ram[5*d+row] = digits[d][row]
		}
	}
}

func (chip8 *Chip8) loadRom(path string) {
	code, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	bytesRead, err := code.Read(chip8.ram[0x200:])
	if bytesRead == 0 {
		panic("The file is empty")
	} else if err != nil {
		panic(err)
	}
}

func (chip8 *Chip8) load(path string) {
	chip8.loadDigits()
	chip8.PC = 0x200
	chip8.loadRom(path)
}

func (chip8 *Chip8) getNextInstruction() uint16 {
	byte1 := uint16(chip8.ram[chip8.PC]) << 8
	byte2 := uint16(chip8.ram[chip8.PC+1])
	return byte1 | byte2
}

func (chip8 *Chip8) run() {
	for {
		for chip8.DT > 0 {
			chip8.DT--
			time.Sleep(16 * time.Millisecond) // 60Hz
		}
		if chip8.ST > 0 {
			playSound(chip8.ST)
			chip8.ST = 0
		}
		instruction := chip8.getNextInstruction()
		fmt.Printf("Opcode: %4X PC: %3X\n", instruction, chip8.PC)

		if instruction&0xF000 == 1<<12 && instruction&0x0FFF == chip8.PC {
			fmt.Println("infinite loop, shutting down the emulator")
			return
		}

		chip8.execute(instruction)

		fmt.Printf("%3X %X %X %v\n", chip8.I, chip8.SP, chip8.stack, chip8.V)
		time.Sleep(2 * time.Millisecond) // 500Hz
	}
}
