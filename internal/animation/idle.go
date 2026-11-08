package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const BreathPeriod = 240 // Four seconds at the game's fixed 60 Hz.

// Idle keeps animation time in Update, independent of render frequency.
type Idle struct {
	walk     Walk
	Tick     int
	Weight   float64
	vertices []ebiten.Vertex
	indices  []uint16
}

func (a *Idle) Update(moving bool) {
	a.walk.Update(moving)
	a.Tick = (a.Tick + 1) % BreathPeriod
	if moving {
		a.Weight = math.Max(0, a.Weight-1.0/12)
	} else {
		a.Weight = math.Min(1, a.Weight+1.0/12)
	}
}

// offset moves the torso and hands locally; hips, legs and feet stay anchored.
// Coordinates are relative to the full square starter sprite cell.
func (a *Idle) offset(x, y float64) (float64, float64) {
	phase := float64(a.Tick) * 2 * math.Pi / BreathPeriod
	breath := (1 - math.Cos(phase)) * 0.5 * a.Weight
	upper := 1 - smooth(0.40, 0.67, y)
	chest := math.Exp(-math.Pow((y-0.43)/0.13, 2))
	side := math.Min(1, math.Abs(x-0.5)/0.16)
	hand := math.Exp(-math.Pow((y-0.55)/0.10, 2)) * side
	direction := math.Tanh((x - 0.5) * 12)
	dx := direction * (0.008*chest + 0.005*hand) * breath
	dy := -0.012*upper*breath + 0.004*hand*math.Sin(phase)*a.Weight
	return dx, dy
}

func smooth(lo, hi, x float64) float64 {
	t := math.Max(0, math.Min(1, (x-lo)/(hi-lo)))
	return t * t * (3 - 2*t)
}

func (a *Idle) Draw(dst, src *ebiten.Image, cx, floorY, height float64) {
	const cols, rows = 20, 32
	if a.vertices == nil {
		a.vertices = make([]ebiten.Vertex, (cols+1)*(rows+1))
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				i := uint16(y*(cols+1) + x)
				j := i + cols + 1
				a.indices = append(a.indices, i, i+1, j, i+1, j+1, j)
			}
		}
	}
	b := src.Bounds()
	width := height * float64(b.Dx()) / float64(b.Dy())
	for row := 0; row <= rows; row++ {
		for col := 0; col <= cols; col++ {
			x, y := float64(col)/cols, float64(row)/rows
			dx, dy := a.offset(x, y)
			wx, wy := a.walk.offset(x, y)
			dx += wx
			dy += wy
			a.vertices[row*(cols+1)+col] = ebiten.Vertex{
				DstX: float32(cx - width/2 + (x+dx)*width), DstY: float32(floorY - height + (y+dy)*height),
				SrcX: float32(b.Min.X) + float32(x*float64(b.Dx())), SrcY: float32(b.Min.Y) + float32(y*float64(b.Dy())),
				ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
			}
		}
	}
	dst.DrawTriangles(a.vertices, a.indices, src, &ebiten.DrawTrianglesOptions{Filter: ebiten.FilterNearest})
}
