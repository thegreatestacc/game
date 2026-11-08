package ui

import "github.com/hajimehoshi/ebiten/v2"

// Input is sampled once per Update, before any screen transition.
type Input struct {
	Interact      bool
	MoveX, MoveY  float64
	X, Y          int
	Click, Escape bool
}

type Button struct {
	X, Y, W, H int
	Label      string
	hover      bool
}

func (b *Button) Update(in Input) bool {
	b.hover = in.X >= b.X && in.X < b.X+b.W && in.Y >= b.Y && in.Y < b.Y+b.H
	return b.hover && in.Click
}

func (b Button) Draw(dst *ebiten.Image, p Painter) {
	c := Panel
	if b.hover {
		c = Border
	}
	Box(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), Gold)
	Box(dst, float32(b.X+1), float32(b.Y+1), float32(b.W-2), float32(b.H-2), c)
	p.Center(dst, b.Label, float64(b.X+b.W/2), float64(b.Y)+(float64(b.H)-28)/2, 24, Ink)
}
