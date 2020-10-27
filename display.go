package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

var scale = 16
var emptyImage = ebiten.NewImage(scale, scale)

func init() {
	emptyImage.Fill(color.White)
}

// Game implements ebiten.Game interface.
type Game struct {
	Chip8
	// display [64][32]bool
	// display [32]uint64
}

func rect(x0, y0 float32, offset uint16) ([]ebiten.Vertex, []uint16) {
	var r, g, b, a float32 = 1, 1, 1, 1
	x1 := x0 + float32(scale)
	y1 := y0 + float32(scale)

	return []ebiten.Vertex{
			{
				DstX:   x0,
				DstY:   y0,
				SrcX:   1,
				SrcY:   1,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			},
			{
				DstX:   x1,
				DstY:   y0,
				SrcX:   1,
				SrcY:   1,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			},
			{
				DstX:   x0,
				DstY:   y1,
				SrcX:   1,
				SrcY:   1,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			},
			{
				DstX:   x1,
				DstY:   y1,
				SrcX:   1,
				SrcY:   1,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			},
		}, []uint16{0 + offset, 1 + offset, 2 + offset,
			1 + offset, 2 + offset, 3 + offset}
}

func rects(x uint64, y int) ([]ebiten.Vertex, []uint16) {
	var vertices []ebiten.Vertex
	var indices []uint16
	for j := 0; j < 64; j++ {
		if (x<<uint(j))&(1<<63) == 1<<63 {
			v, i := rect(float32(j*scale), float32(y*scale), uint16(len(vertices)))
			vertices = append(vertices, v...)
			indices = append(indices, i...)
		}
	}
	return vertices, indices
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// n := rand.Intn(31)
	// g.display[n] = rand.Uint64()

	instruction := g.getNextInstruction()
	fmt.Printf("Opcode: %4X PC: %3X\n", instruction, g.PC)

	if instruction == 0x1228 {
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	g.execute(instruction)
	if g.ST > 0 {
		// buzz
	} else {
		// stop buzzing
	}
	for g.DT > 0 {
		g.DT--
		// time.Sleep(1 / 60 * time.Second)
	}
	g.PC += 2
	fmt.Printf("%3X %X %v %v\n", g.I, g.SP, g.stack, g.V)
	time.Sleep(100 * time.Millisecond)
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < screenHeight; y++ {
		// v, i := rect(float32(x*scale), float32(y*scale))
		v, i := rects(g.display[y], y)
		screen.DrawTriangles(v, i, emptyImage, nil)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 64 * scale, 32 * scale
}

func main() {
	game := &Game{Chip8{}}
	game.load("./roms/IBM Logo.ch8")
	ebiten.SetWindowSize(screenWidth*scale, screenHeight*scale)
	ebiten.SetWindowTitle("Chip8 Emulator")
	// ebiten.SetMaxTPS(6)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
