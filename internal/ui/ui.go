package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	Background = color.NRGBA{R: 13, G: 18, B: 29, A: 255}
	Panel      = color.NRGBA{R: 22, G: 30, B: 44, A: 255}
	Border     = color.NRGBA{R: 49, G: 65, B: 82, A: 255}
	Ink        = color.NRGBA{R: 232, G: 232, B: 224, A: 255}
	Muted      = color.NRGBA{R: 143, G: 159, B: 174, A: 255}
	Gold       = color.NRGBA{R: 208, G: 174, B: 108, A: 255}
	Cyan       = color.NRGBA{R: 114, G: 204, B: 214, A: 255}
)

type Painter struct{ Font *text.GoTextFaceSource }

func (p Painter) Text(dst *ebiten.Image, s string, x, y, size float64, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, s, &text.GoTextFace{Source: p.Font, Size: size}, op)
}

func (p Painter) Center(dst *ebiten.Image, s string, x, y, size float64, c color.Color) {
	w, _ := text.Measure(s, &text.GoTextFace{Source: p.Font, Size: size}, size*1.3)
	p.Text(dst, s, x-w/2, y, size, c)
}

func Box(dst *ebiten.Image, x, y, w, h float32, c color.Color) {
	vector.DrawFilledRect(dst, x, y, w, h, c, false)
}

func Frame(dst *ebiten.Image, x, y, w, h float32) {
	Box(dst, x, y, w, h, Panel)
	vector.StrokeRect(dst, x, y, w, h, 1, Border, false)
}

func Backdrop(dst *ebiten.Image) {
	dst.Fill(Background)
	for i := 0; i < 60; i++ {
		x := float32((i*173 + 47) % 1280)
		y := float32((i*97 + 23) % 720)
		vector.DrawFilledCircle(dst, x, y, 1, Border, true)
	}
	vector.StrokeLine(dst, 48, 660, 1232, 660, 1, Border, false)
}

func Sigil(dst *ebiten.Image, x, y, r float32) {
	vector.StrokeCircle(dst, x, y, r, 1, Border, true)
	vector.StrokeCircle(dst, x, y, r*0.8, 1, Border, true)
	for i := 0; i < 6; i++ {
		a := float64(i)*math.Pi/3 - math.Pi/2
		b := a + 2*math.Pi/3
		vector.StrokeLine(dst, x+r*float32(math.Cos(a)), y+r*float32(math.Sin(a)), x+r*float32(math.Cos(b)), y+r*float32(math.Sin(b)), 1, Border, true)
	}
}
