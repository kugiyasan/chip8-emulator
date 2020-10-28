package main

import (
	"fmt"
	"math/rand"
)

func (chip8 *Chip8) execute(instruction uint16) {
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
			chip8.display = [32]uint64{}
		case 0x00EE:
			// RET: Return from a subroutine
			chip8.SP--
			chip8.PC = chip8.stack[chip8.SP]
		default: // 0nnn
			// SYS addr: Jump to a machine code routine at nnn
			// It is ignored by modern interpreters
			fmt.Printf("unknown opcode %4X\n", instruction)
			// err := fmt.Errorf("unknown opcode %4X", instruction)
			// panic(err)
		}
	case 0x1: // 1nnn
		// JP addr: The interpreter sets the program counter to nnn
		chip8.PC = instruction&0x0FFF - 2
	case 0x2: // 2nnn
		// CALL addr: Call subroutine at nnn
		chip8.stack[chip8.SP] = chip8.PC
		chip8.SP++
		chip8.PC = instruction&0x0FFF - 2
	case 0x3: // 3xkk
		// SE Vx, byte: Skip next instruction if Vx = kk
		if chip8.V[x] == uint8(instruction&0xFF) {
			chip8.PC += 2
		}
	case 0x4: // 4xkk
		// SE Vx, byte: Skip next instruction if Vx != kk
		if chip8.V[x] != uint8(instruction&0xFF) {
			chip8.PC += 2
		}
	case 0x5: // 5xy0
		// SE Vx, byte: Skip next instruction if Vx = Vy
		y := instruction & 0xF0 >> 4
		if chip8.V[x] == chip8.V[y] {
			chip8.PC += 2
		}
	case 0x6: // 6xkk
		// LD Vx, byte: Set Vx = kk
		chip8.V[x] = uint8(instruction & 0xFF)
	case 0x7: // 7xkk
		// ADD Vx, byte: Adds the value kk to the value of
		// register Vx, then stores the result in Vx
		chip8.V[x] += uint8(instruction & 0xFF)
	case 0x8:
		y := instruction & 0xF0 >> 4
		switch instruction & 0xF {
		case 0x0: // 8xy0
			// LD Vx, Vy: Set Vx = Vy
			chip8.V[x] = chip8.V[y]
		case 0x1: // 8xy1
			// OR Vx, Vy: Set Vx = Vx OR Vy
			chip8.V[x] |= chip8.V[y]
		case 0x2: // 8xy2
			// AND Vx, Vy: Set Vx = Vx AND Vy
			chip8.V[x] &= chip8.V[y]
		case 0x3: // 8xy3
			// XOR Vx, Vy: Set Vx = Vx XOR Vy
			chip8.V[x] ^= chip8.V[y]
		case 0x4: // 8xy4
			// ADD Vx, Vy: Set Vx = Vx + Vy, set VF = carry
			a := chip8.V[x]
			chip8.V[x] += chip8.V[y]
			// ((c < a) != (b < 0)) where c := a + b
			if (chip8.V[x] < a) != (chip8.V[y] < 0) {
				chip8.V[0xF] = 1
			}
		case 0x5: // 8xy5
			// SUB Vx, Vy: Set Vx = Vx - Vy, set VF = NOT borrow
			if chip8.V[x] > chip8.V[y] {
				chip8.V[0xF] = 1
			}
			chip8.V[x] -= chip8.V[y]
		case 0x6: // 8xy6
			// SHR Vx {, Vy}: Set Vx = Vx SHR 1
			chip8.V[0xF] = chip8.V[x] & 0x1
			chip8.V[x] >>= 1
		case 0x7: // 8xy7
			// SUBN Vx, Vy: Set Vx = Vy - Vx, set VF = NOT borrow
			if chip8.V[y] > chip8.V[x] {
				chip8.V[0xF] = 1
			}
			chip8.V[y] -= chip8.V[x]
		case 0xE: // 8xyE
			// SHL Vx {, Vy}: Set Vx = Vx SHL 1
			chip8.V[0xF] = (chip8.V[x] & 0x80) >> 7
			chip8.V[x] <<= 1
		default:
			fmt.Printf("unknown opcode %4X\n", instruction)
		}

	case 0x9: // 9xy0
		// SNE Vx, Vy: Skip next instruction if Vx != Vy
		y := instruction & 0x00F0 >> 4
		if chip8.V[x] != chip8.V[y] {
			chip8.PC += 2
		}
	case 0xA: // Annn
		// LD I, addr: Set I = nnn
		chip8.I = instruction & 0x0FFF
	case 0xB: // Bnnn
		// JP V0, addr: Jump to location nnn + V0
		chip8.PC = instruction&0x0FFF + uint16(chip8.V[0x0])
	case 0xC: // Cxkk
		// RND Vx, byte: Set Vx = random byte AND kk
		chip8.V[x] = uint8(rand.Intn(1<<8)) & uint8(instruction&0xFF)
	case 0xD: // Dxyn
		// DRW Vx, Vy, nibble: Display n-byte sprite starting
		// at memory location I at (Vx, Vy), set VF = collision
		y := (instruction & 0x00F0) >> 4
		sprites := uint8(instruction & 0x000F)
		location := chip8.I
		shift := 56 - int(chip8.V[x])
		for n := uint8(0); n < sprites; n++ {
			if int(chip8.V[y]+n) >= len(chip8.display) {
				break
			}

			row := (uint64(chip8.ram[location]) & 0xFF)
			if shift > 0 {
				row <<= shift
			} else {
				row >>= -shift
			}
			if chip8.display[chip8.V[y]+n]&row != 0 {
				chip8.V[0xF] = 1
			}
			chip8.display[chip8.V[y]+n] ^= row
			location++
		}
	case 0xE:
		switch instruction & 0xFF {
		case 0x9E: // Ex9E
			// SKP Vx: Skip next instruction if key with the value of Vx is pressed
			if chip8.isPressed(chip8.V[x]) {
				chip8.PC += 2
			}
		case 0xA1: // ExA1
			// SKNP Vx:  Skip next instruction if key with the value of Vx is not pressed
			if !chip8.isPressed(chip8.V[x]) {
				chip8.PC += 2
			}
		default:
			fmt.Printf("unknown opcode %4X\n", instruction)
		}
	case 0xF:
		switch instruction & 0xFF {
		case 0x07: // Fx07
			// LD Vx, DT: Set Vx = delay timer value
			chip8.V[x] = chip8.DT
		case 0x0A: // Fx0A
			// LD Vx, K: Wait for a key press, store the value of the key in Vx
			chip8.V[x] = chip8.getKeyPress()
		case 0x15: // Fx15
			// LD DT, Vx: Set delay timer = Vx
			chip8.DT = uint8(x)
		case 0x18: // Fx18
			// LD ST, Vx: Set sound timer = Vx
			chip8.ST = uint8(x)
		case 0x1E: // Fx1E
			// ADD I, Vx: Set I = I + Vx
			chip8.I += uint16(chip8.V[x])
		case 0x29: // Fx29
			// LD F, Vx: Set I = location of sprite for digit Vx
			// The sprites will be located at the beginning
			chip8.I = uint16(5 * chip8.V[x])
		case 0x33: // Fx33
			// LD B, Vx: Store BCD representation of Vx
			// in memory locations I, I+1, and I+2
			BCD := chip8.V[x]
			chip8.ram[int(chip8.I)] = BCD / 100
			chip8.ram[int(chip8.I)+1] = BCD / 10 % 10
			chip8.ram[int(chip8.I)+2] = BCD % 10
		case 0x55: // Fx55
			// LD [I], Vx: Store registers V0 through Vx
			// in memory starting at location I
			location := chip8.I
			for register := uint16(0); register < x; register++ {
				chip8.ram[location] = chip8.V[register]
				location++
			}
		case 0x65: // Fx65
			// LD Vx, [I]: Read registers V0 through Vx
			// from memory starting at location I
			location := chip8.I
			for register := uint16(0); register <= x; register++ {
				chip8.V[register] = chip8.ram[location]
				location++
			}
		default:
			fmt.Printf("unknown opcode %4X\n", instruction)
		}
	default:
		// fmt.Printf("unknown opcode %4X\n", instruction)
		err := fmt.Errorf("unknown opcode %4X", instruction)
		panic(err)
	}
	chip8.PC += 2
}
