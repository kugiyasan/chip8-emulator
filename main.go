package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func getKeyPress() uint8 {
	return 0xF
}

func chip8() {
	ram := make([]byte, 4096)
	// RAM address: 0x000 to 0xFFF
	// Interpreter location: 0x000 to 0x1FF
	// Sprite location: 0x000 to 0x050

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
			ram[5*d+row] = digits[d][row]
		}
	}

	V := make([]uint8, 16)
	// 16x 8-bit registers: V0 to VF
	// VF stores flags

	var I uint16
	// 16-bit register I
	// store memory address, 12 bits are usually used

	var PC uint16 = 0x200
	// 16-bit PC program counter

	var SP uint8
	// 8-bit SP stack pointer

	stack := make([]uint16, 16)
	// store the address that the interpreter
	// shoud return to when finished with a subroutine

	keyboard := make([]bool, 16)

	display := make([]uint64, 32)
	// storing the display image (64x32) by row in a uint64
	// should let the array length change to 64x48 or 64x64
	// should support 128x64 too

	var DT uint8
	// DT: delay timer
	var ST uint8
	// ST: sound timer

	code, err := os.Open("./roms/IBM Logo.ch8")
	if err != nil {
		fmt.Println(err)
	}

	bytesRead, err := code.Read(ram[0x200:])
	if bytesRead == 0 {
		panic("The file contains nothing")
	} else if err != nil {
		panic(err)
	}

	fmt.Printf("%X\n", ram[0x200:0x202])

	for {
		time.Sleep(100 * time.Millisecond)
		instruction := uint16(ram[PC])<<8 | uint16(ram[PC+1])
		fmt.Printf("Opcode: %4X PC: %3X\n", instruction, PC)

		// x is used a lot in opcodes
		x := instruction & 0xF00 >> 8

		// nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
		// n or nibble - A 4-bit value, the lowest 4 bits of the instruction
		// x - A 4-bit value, the lower 4 bits of the high byte of the instruction
		// y - A 4-bit value, the upper 4 bits of the low byte of the instruction
		// kk or byte - An 8-bit value, the lowest 8 bits of the instruction
		switch instruction >> 12 {
		case 0x0:
			switch instruction {
			case 0x00E0:
				// CLS: Clear the display
				display = make([]uint64, 32)
			case 0x00EE:
				// RET: Return from a subroutine
				PC = stack[SP]
				SP--
			default: // 0nnn
				// SYS addr: Jump to a machine code routine at nnn
				// It is ignored by modern interpreters
			}
		case 0x1: // 1nnn
			// JP addr: The interpreter sets the program counter to nnn
			PC = instruction & 0x0FFF
		case 0x2: // 2nnn
			// CALL addr: Call subroutine at nnn
			stack[SP] = PC
			SP++
			PC = instruction & 0x0FFF
		case 0x3: // 3xkk
			// SE Vx, byte: Skip next instruction if Vx = kk
			if V[x] == uint8(instruction&0xFF) {
				PC += 2
			}
		case 0x4: // 4xkk
			// SE Vx, byte: Skip next instruction if Vx != kk
			if V[x] != uint8(instruction&0xFF) {
				PC += 2
			}
		case 0x5: // 5xy0
			// SE Vx, byte: Skip next instruction if Vx = Vy
			y := instruction & 0xF0 >> 4
			if V[x] == V[y] {
				PC += 2
			}
		case 0x6: // 6xkk
			// LD Vx, byte: Set Vx = kk
			V[x] = uint8(instruction & 0xFF)
		case 0x7: // 7xkk
			// ADD Vx, byte: Adds the value kk to the value of
			// register Vx, then stores the result in Vx
			V[x] += uint8(instruction & 0xFF)
		case 0x8:
			y := instruction & 0xF0 >> 4
			switch instruction & 0xF {
			case 0x0: // 8xy0
				// LD Vx, Vy: Set Vx = Vy
				V[x] = V[y]
			case 0x1: // 8xy1
				// OR Vx, Vy: Set Vx = Vx OR Vy
				V[x] |= V[y]
			case 0x2: // 8xy2
				// AND Vx, Vy: Set Vx = Vx AND Vy
				V[x] &= V[y]
			case 0x3: // 8xy3
				// XOR Vx, Vy: Set Vx = Vx XOR Vy
				V[x] ^= V[y]
			case 0x4: // 8xy4
				// ADD Vx, Vy: Set Vx = Vx + Vy, set VF = carry
				a := V[x]
				V[x] += V[y]
				// ((c < a) != (b < 0)) where c := a + b
				if (V[x] < a) != (V[y] < 0) {
					V[0xF] = 1
				}
			case 0x5: // 8xy5
				// SUB Vx, Vy: Set Vx = Vx - Vy, set VF = NOT borrow
				if V[x] > V[y] {
					V[0xF] = 1
				}
				V[x] -= V[y]
			case 0x6: // 8xy6
				// SHR Vx {, Vy}: Set Vx = Vx SHR 1
				V[0xF] = V[x] & 0x1
				V[x] >>= 1
			case 0x7: // 8xy7
				// SUBN Vx, Vy: Set Vx = Vy - Vx, set VF = NOT borrow
				if V[y] > V[x] {
					V[0xF] = 1
				}
				V[y] -= V[x]
			case 0x8: // 8xy8
				// SHL Vx {, Vy}: Set Vx = Vx SHL 1
				V[0xF] = V[x] & 0x80
				V[x] <<= 1
			}

		case 0x9: // 9xy0
			// SNE Vx, Vy: Skip next instruction if Vx != Vy
			if V[x] != V[instruction&0x00F0] {
				PC += 2
			}
		case 0xA: // Annn
			// LD I, addr: Set I = nnn
			I = instruction & 0x0FFF
		case 0xB: // Bnnn
			// JP V0, addr: Jump to location nnn + V0
			PC = instruction&0x0FFF + uint16(V[0x0])
		case 0xC: // Cxkk
			// RND Vx, byte: Set Vx = random byte AND kk
			V[x] = uint8(rand.Intn(1<<8)) | uint8(instruction&0xFF)
		case 0xD: // Dxyn
			// DRW Vx, Vy, nibble: Display n-byte sprite starting
			// at memory location I at (Vx, Vy), set VF = collision
			y := instruction & 0x00F0 >> 4
			nibble := int(instruction & 0x000F)
			location := I
			for n := 0; n < nibble; n++ {
				row := uint64(ram[location]) << (8 * (7 - y))
				if display[x]&row != 0 {
					V[0xF] = 1
				}
				display[x] ^= row
				location++
				y++
			}
		case 0xE:
			switch instruction & 0xFF {
			case 0x9E: // Ex9E
				// SKP Vx: Skip next instruction if key with the value of Vx is pressed
				if keyboard[x] {
					PC += 2
				}
			case 0xA1: // ExA1
				// SKNP Vx:  Skip next instruction if key with the value of Vx is not pressed
				if !keyboard[x] {
					PC += 2
				}
			}
		case 0xF:
			switch instruction & 0xFF {
			case 0x07: // Fx07
				// LD Vx, DT: Set Vx = delay timer value
				V[x] = DT
			case 0x0A: // Fx0A
				// LD Vx, K: Wait for a key press, store the value of the key in Vx
				V[x] = getKeyPress()
			case 0x15: // Fx15
				// LD DT, Vx: Set delay timer = Vx
				DT = uint8(x)
			case 0x18: // Fx18
				// LD ST, Vx: Set sound timer = Vx
				ST = uint8(x)
			case 0x1E: // Fx1E
				// ADD I, Vx: Set I = I + Vx
				I += uint16(V[x])
			case 0x29: // Fx29
				// LD F, Vx: Set I = location of sprite for digit Vx
				// The sprites will be located at the beginning
				I = uint16(5 * x)
			case 0x33: // Fx33
				// LD B, Vx: Store BCD representation of Vx in memory locations I, I+1, and I+2
				BCD := fmt.Sprintf("%d", V[x])
				for i := range BCD {
					ram[int(I)+i] = BCD[i]
				}
			case 0x55: // Fx55
				// LD [I], Vx: Store registers V0 through Vx in memory starting at location I
				location := I
				for register := uint16(0); register < x; register++ {
					ram[location] = V[register]
					location++
				}
			case 0x65: // Fx65
				// LD Vx, [I]: Read registers V0 through Vx from memory starting at location I
				location := I
				for register := uint16(0); register < x; register++ {
					V[register] = ram[location]
					location++
				}
			}
		}
		if ST > 0 {
			// buzz
		} else {
			// stop buzzing
		}
		for DT > 0 {
			DT--
			time.Sleep(1 / 60 * time.Second)
		}
		PC += 2
	}
	fmt.Println(display)
}

func main() {
	chip8()
}
