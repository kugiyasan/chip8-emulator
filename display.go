package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

var scale = 16
var emptyImage, displayErr = ebiten.NewImage(scale, scale, ebiten.FilterNearest)

func init() {
	emptyImage.Fill(color.White)
}

// Game implements ebiten.Game interface
type Game struct{ Chip8 }

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
func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

var hue uint = 0

func getRGBFromHue(hue uint) (uint8, uint8, uint8) {
	// https://www.rapidtables.com/convert/color/hsv-to-rgb.html
	a := 1 - math.Abs(math.Mod(float64(hue)/60, 2)-1)
	x := uint8(a * 255)

	switch {
	case hue < 60:
		return 255, x, 0
	case hue < 120:
		return x, 255, 0
	case hue < 180:
		return 0, 255, x
	case hue < 240:
		return 0, x, 255
	case hue < 300:
		return x, 0, 255
	case hue < 360:
		return 255, 0, x
	}
	return 255, 255, 255
}

func updateFillColor() {
	var red, green, blue uint8 = getRGBFromHue(hue)

	emptyImage.Fill(color.RGBA{red, green, blue, 255})
	hue = (hue + 1) % 360
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	updateFillColor()

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
	// return 64, 32
}
