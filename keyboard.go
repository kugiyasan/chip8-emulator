package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten"
)

func (chip8 *Chip8) getKeyPress() uint8 {
	for {
		for i := uint8(0); i < 16; i++ {
			if chip8.isPressed(i) {
				fmt.Printf("KEY PRESSED %X\n", i)
				return uint8(i)
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (chip8 *Chip8) isPressed(key uint8) bool {
	keys := []ebiten.Key{
		ebiten.KeyKP0,
		ebiten.KeyKP1,
		ebiten.KeyKP2,
		ebiten.KeyKP3,
		ebiten.KeyKP4,
		ebiten.KeyKP5,
		ebiten.KeyKP6,
		ebiten.KeyKP7,
		ebiten.KeyKP8,
		ebiten.KeyKP9,
		ebiten.KeyA,
		ebiten.KeyB,
		ebiten.KeyC,
		ebiten.KeyD,
		ebiten.KeyE,
		ebiten.KeyF,
	}
	return ebiten.IsKeyPressed(keys[key])
}
