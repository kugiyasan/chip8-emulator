package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestChip8_execute(t *testing.T) {
	type args struct {
		instruction uint16
	}

	chip8 := Chip8{}
	s, _ := strconv.ParseUint("D01F", 16, 16)

	tests := []struct {
		name  string
		chip8 *Chip8
		args  args
	}{
		{
			"Fill screen with white",
			&chip8,
			args{uint16(s)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instructions := []string{"A22A", "6040", "6108", "D01F", "610F", "D01F"}
			for _, s := range instructions {
				instruction, _ := strconv.ParseUint(s, 16, 16)
				tt.chip8.execute(uint16(instruction))
				fmt.Println(tt.chip8.display)
			}
		})
	}
}
