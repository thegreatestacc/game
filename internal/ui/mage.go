package ui

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type point3 struct{ x, y, z float64 }
type mageFace struct {
	points []point3
	tint   color.NRGBA
	depth  float64
}

// DrawTurningMage projects a small procedural model into the 2D portrait.
// Depth ordering makes the staff and face disappear behind the robe on a turn.
func DrawTurningMage(dst *ebiten.Image, cx, floorY, angle float64) {
	faces := make([]mageFace, 0, 160)
	add := func(c color.NRGBA, p ...point3) { faces = append(faces, mageFace{points: p, tint: c}) }
	type ring struct{ y, rx, rz float64 }
	surface := func(rings []ring, base color.NRGBA, hood bool) {
		const sides = 16
		for j := 0; j < len(rings)-1; j++ {
			for i := 0; i < sides; i++ {
				a, b := float64(i)*2*math.Pi/sides, float64(i+1)*2*math.Pi/sides
				p := func(r ring, t float64) point3 { return point3{r.rx * math.Sin(t), r.y, r.rz * math.Cos(t)} }
				c := base
				if hood && j == 1 && (i == 0 || i == 1 || i == 14 || i == 15) {
					c = Background
				}
				light := 0.72 + 0.28*math.Cos((a+b)/2+angle-0.6)
				c.R = uint8(float64(c.R) * light)
				c.G = uint8(float64(c.G) * light)
				c.B = uint8(float64(c.B) * light)
				add(c, p(rings[j], a), p(rings[j], b), p(rings[j+1], b), p(rings[j+1], a))
			}
		}
	}
	surface([]ring{{-155, 25, 23}, {-132, 48, 29}, {-88, 43, 30}, {0, 73, 43}}, color.NRGBA{R: 79, G: 103, B: 147, A: 255}, false)
	surface([]ring{{-235, 0, 0}, {-216, 34, 29}, {-174, 36, 33}, {-151, 23, 23}}, color.NRGBA{R: 67, G: 88, B: 125, A: 255}, true)
	surface([]ring{{-91, 44, 31}, {-84, 45, 32}}, Gold, false)
	// Gold seam on the front of the robe.
	add(Gold, point3{-2, -132, 30}, point3{2, -132, 30}, point3{3, 0, 44}, point3{-3, 0, 44})
	// Luminous eyes sit inside the dark hood opening.
	for _, x := range []float64{-13, 7} {
		add(Cyan, point3{x, -197, 32}, point3{x + 6, -197, 32}, point3{x + 6, -193, 33}, point3{x, -193, 33})
	}
	// Staff with a faceted crystal, attached to the same rotating model.
	for i := 0; i < 6; i++ {
		a, b := float64(i)*math.Pi/3, float64(i+1)*math.Pi/3
		p := func(t, y float64) point3 { return point3{-87 + 3*math.Sin(t), y, 12 + 3*math.Cos(t)} }
		add(Gold, p(a, 0), p(b, 0), p(b, -187), p(a, -187))
		q := func(t float64) point3 { return point3{-87 + 14*math.Sin(t), -205, 12 + 14*math.Cos(t)} }
		c := Cyan
		if i%2 == 0 {
			c = color.NRGBA{R: 166, G: 240, B: 239, A: 255}
		}
		add(c, point3{-87, -231, 12}, q(a), q(b))
		add(c, point3{-87, -182, 12}, q(b), q(a))
	}
	co, si := math.Cos(angle), math.Sin(angle)
	for i := range faces {
		for j, p := range faces[i].points {
			p.x, p.z = p.x*co+p.z*si, -p.x*si+p.z*co
			faces[i].points[j] = p
			faces[i].depth += p.z
		}
		faces[i].depth /= float64(len(faces[i].points))
	}
	sort.SliceStable(faces, func(i, j int) bool { return faces[i].depth < faces[j].depth })
	// Flattened ground rings make the vertical axis of rotation clear.
	for r := 0; r < 3; r++ {
		var path vector.Path
		for i := 0; i <= 64; i++ {
			t := float64(i) * math.Pi / 32
			x, y := float32(cx+(100-float64(r)*12)*math.Cos(t)), float32(floorY+10+(22-float64(r)*3)*math.Sin(t))
			if i == 0 {
				path.MoveTo(x, y)
			} else {
				path.LineTo(x, y)
			}
		}
		op := &vector.DrawPathOptions{AntiAlias: true}
		op.ColorScale.ScaleWithColor(Border)
		vector.StrokePath(dst, &path, &vector.StrokeOptions{Width: 1}, op)
	}
	for _, face := range faces {
		var path vector.Path
		for i, p := range face.points {
			x, y := float32(cx+p.x), float32(floorY+p.y+p.z*0.16)
			if i == 0 {
				path.MoveTo(x, y)
			} else {
				path.LineTo(x, y)
			}
		}
		path.Close()
		op := &vector.DrawPathOptions{AntiAlias: true}
		op.ColorScale.ScaleWithColor(face.tint)
		vector.FillPath(dst, &path, nil, op)
	}
}
