package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
)

const (
	sampleRate = 44100
	frequency  = 880
)

var audioContext = audio.NewContext(sampleRate)

func createWave(ST uint8) []byte {
	size := sampleRate * int(ST) / 60
	wave := make([]byte, size+4-size%4)
	p := 0

	const length = int64(sampleRate / frequency)
	for i := 0; i+3 < len(wave); i += 4 {
		b := int16(math.Sin(2*math.Pi*float64(p)/float64(length)) * 32767)
		wave[i] = byte(b)
		wave[i+1] = byte(b >> 8)
		wave[i+2] = byte(b)
		wave[i+3] = byte(b >> 8)
		p++
	}
	fmt.Println(len(wave))
	return wave
}

func playSound(ST uint8) {
	player := audio.NewPlayerFromBytes(audioContext, createWave(ST))
	fmt.Println("Playing")
	player.Play()
	time.Sleep(2 * time.Second)
}

// func main() {
// 	playSound(10)
// }
