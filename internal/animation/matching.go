package animation

import (
	"image"
	"math"
)

// matchRows uses monotone sequence alignment to match color/alpha features
// along each scanline. Unlike an unconstrained displacement field, the
// correspondence cannot cross itself or turn the portrait mesh inside out.
func matchRows(a, b *image.NRGBA) [2][]point {
	const samples = 128
	type feature [4]float64
	getRow := func(img *image.NRGBA, row int) [samples]feature {
		var values [samples]feature
		w, h := img.Bounds().Dx(), img.Bounds().Dy()
		for i := range values {
			x := min(w-1, i*w/(samples-1))
			y := row * h / rows
			for dy := -3; dy <= 3; dy++ {
				c := img.NRGBAAt(x, max(0, min(h-1, y+dy)))
				alpha := float64(c.A) / 255
				values[i][0] += float64(c.R) / 255 * alpha / 7
				values[i][1] += float64(c.G) / 255 * alpha / 7
				values[i][2] += float64(c.B) / 255 * alpha / 7
				values[i][3] += alpha / 7
			}
		}
		return values
	}
	n := (rows + 1) * (columns + 1)
	result := [2][]point{make([]point, n), make([]point, n)}
	cost := make([]float64, samples*samples)
	parent := make([]byte, len(cost))
	for row := 0; row <= rows; row++ {
		av, bv := getRow(a, row), getRow(b, row)
		for x := 0; x < samples; x++ {
			for y := 0; y < samples; y++ {
				idx := x*samples + y
				penalty := 0.0
				for k := 0; k < 4; k++ {
					d := av[x][k] - bv[y][k]
					penalty += d * d
				}
				penalty += 0.5 * math.Pow(float64(x-y)/samples, 2)
				best := math.Inf(1)
				direction := byte(0)
				if x > 0 && y > 0 {
					best = cost[idx-samples-1]
				}
				if x > 0 && cost[idx-samples]+0.018 < best {
					best = cost[idx-samples] + 0.018
					direction = 1
				}
				if y > 0 && cost[idx-1]+0.018 < best {
					best = cost[idx-1] + 0.018
					direction = 2
				}
				if x == 0 && y == 0 {
					best = 0
				}
				cost[idx] = best + penalty
				parent[idx] = direction
			}
		}
		type node struct{ x, y int }
		path := make([]node, 0, 2*samples)
		for x, y := samples-1, samples-1; ; {
			path = append(path, node{x, y})
			if x == 0 && y == 0 {
				break
			}
			switch parent[x*samples+y] {
			case 0:
				x--
				y--
			case 1:
				x--
			case 2:
				y--
			}
		}
		// Resample by average position, so both textures share one mesh.
		cursor := len(path) - 1
		for col := 0; col <= columns; col++ {
			target := float64(col) * 2 * (samples - 1) / columns
			for cursor > 0 && float64(path[cursor-1].x+path[cursor-1].y) < target {
				cursor--
			}
			lo, hi := path[cursor], path[max(0, cursor-1)]
			f := 0.0
			if hi.x+hi.y > lo.x+lo.y {
				f = (target - float64(lo.x+lo.y)) / float64(hi.x+hi.y-lo.x-lo.y)
			}
			idx := row*(columns+1) + col
			result[0][idx] = point{float32((float64(lo.x) + f*float64(hi.x-lo.x)) * float64(a.Bounds().Dx()) / (samples - 1)), float32(row*a.Bounds().Dy()) / rows}
			result[1][idx] = point{float32((float64(lo.y) + f*float64(hi.y-lo.y)) * float64(b.Bounds().Dx()) / (samples - 1)), float32(row*b.Bounds().Dy()) / rows}
		}
	}
	// Vertical smoothing suppresses row-to-row noise without moving the top,
	// bottom or horizontal image boundaries.
	for side := range result {
		original := append([]point(nil), result[side]...)
		for row := 0; row <= rows; row++ {
			for col := 0; col <= columns; col++ {
				x := float32(0)
				for dy := -2; dy <= 2; dy++ {
					x += original[max(0, min(rows, row+dy))*(columns+1)+col].x / 5
				}
				// A small regular grid component guarantees nondegenerate triangles.
				result[side][row*(columns+1)+col].x = x*0.98 + float32(col*a.Bounds().Dx())/columns*0.02
			}
		}
	}
	return result
}
