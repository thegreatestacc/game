// Package animation interpolates the geometry and texture of directional sprites.
package animation

import (
	"image"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const columns, rows = 32, 64

type point struct{ x, y float32 }

type Turntable struct {
	Frames   []*ebiten.Image
	pairs    [][2][]point
	indices  []uint16
	vertices []ebiten.Vertex
	canvas   *ebiten.Image
	size     int
}

// NewTurntable creates correspondence meshes once, on load. Color and alpha
// matching aligns the silhouette and details at every height, so intermediate frames
// change shape instead of merely dissolving two stationary pictures.
func NewTurntable(sheet image.Image) *Turntable {
	b := sheet.Bounds()
	size := max((b.Dx()+3)/4, (b.Dy()+1)/2)
	t := &Turntable{size: size, canvas: ebiten.NewImage(size, size), vertices: make([]ebiten.Vertex, (columns+1)*(rows+1))}
	var images []*image.NRGBA
	for i := 0; i < 8; i++ {
		x, y := i%4, i/4
		r := image.Rect(b.Min.X+x*b.Dx()/4, b.Min.Y+y*b.Dy()/2, b.Min.X+(x+1)*b.Dx()/4, b.Min.Y+(y+1)*b.Dy()/2)
		frame := image.NewNRGBA(image.Rect(0, 0, size, size))
		draw.Draw(frame, image.Rect((size-r.Dx())/2, (size-r.Dy())/2, (size+r.Dx())/2, (size+r.Dy())/2), sheet, r.Min, draw.Src)
		t.Frames = append(t.Frames, ebiten.NewImageFromImage(frame))
		images = append(images, frame)
	}
	for i := range images {
		t.pairs = append(t.pairs, matchRows(images[i], images[(i+1)%len(images)]))
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			a := uint16(y*(columns+1) + x)
			b := a + 1
			c := a + columns + 1
			d := c + 1
			t.indices = append(t.indices, a, b, c, b, d, c)
		}
	}
	return t
}

// phase wraps continuously, including the last-to-first view. There are no
// held frames: every update changes geometry and textures across the interval.
func (t *Turntable) Draw(dst *ebiten.Image, cx, floorY, height, phase float64) {
	n := len(t.Frames)
	position := (phase - math.Floor(phase)) * float64(n)
	i := int(position)
	next := (i + 1) % n
	f := float32(position - float64(i))
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	scale := height / float64(t.size)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx-height/2, floorY-height)
	// Render exact keyframes directly, avoiding numerical degeneracy in an
	// identity-warped triangle mesh and retaining the original sprite quality.
	if f == 0 {
		dst.DrawImage(t.Frames[i], op)
		return
	}
	t.canvas.Clear()
	for pass, frame := range []int{i, next} {
		alpha := 1 - f
		if pass == 1 {
			alpha = f
		}
		if alpha == 0 {
			continue
		}
		for v, a := range t.pairs[i][0] {
			b := t.pairs[i][1][v]
			src := t.pairs[i][pass][v]
			t.vertices[v] = ebiten.Vertex{
				DstX: a.x + (b.x-a.x)*f, DstY: a.y + (b.y-a.y)*f,
				SrcX: src.x, SrcY: src.y,
				ColorR: 1, ColorG: 1, ColorB: 1, ColorA: alpha,
			}
		}
		// Add premultiplied colors in an isolated transparent target. Normal
		// source-over blending would darken the sprite halfway through a turn.
		t.canvas.DrawTriangles(t.vertices, t.indices, t.Frames[frame], &ebiten.DrawTrianglesOptions{Filter: ebiten.FilterLinear, Blend: ebiten.BlendLighter})
	}
	dst.DrawImage(t.canvas, op)
}
