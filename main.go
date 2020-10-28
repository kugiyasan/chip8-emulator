package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten"
)

func main() {
	filename := flag.String("filename", "./roms/Pong (alt).ch8", "rom to load")
	flag.Parse()

	g := &Game{Chip8{}}
	g.load(*filename)

	go g.run()

	ebiten.SetWindowSize(screenWidth*scale, screenHeight*scale)
	ebiten.SetWindowTitle("Chip8 Emulator")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
