package animation

import "math"

// Walk is a front-facing gait. Its phase is independent of movement direction.
type Walk struct {
	Tick   int
	Weight float64
}

const WalkPeriod = 42 // One left/right stride per 0.7 seconds at 60 Hz.

func (w *Walk) Update(moving bool) {
	if moving {
		w.Weight = math.Min(1, w.Weight+1.0/8)
	} else {
		w.Weight = math.Max(0, w.Weight-1.0/10)
	}
	if moving || w.Weight > 0 {
		w.Tick = (w.Tick + 1) % WalkPeriod
	} else {
		w.Tick = 0
	}
}

// The leg and hand masks deform separate parts of the same pixel sprite.
// There is no global image scaling or change to the character's world position.
func (w *Walk) offset(x, y float64) (float64, float64) {
	phase := float64(w.Tick) * 2 * math.Pi / WalkPeriod
	step := math.Sin(phase) * w.Weight
	side := math.Tanh((x - 0.5) * 35)
	legs := smooth(0.62, 0.84, y)
	arms := smooth(0.08, 0.17, math.Abs(x-0.5)) * smooth(0.30, 0.48, y) * (1 - smooth(0.58, 0.68, y))
	body := 1 - smooth(0.64, 0.88, y)
	// Feet alternate up/down; arms swing opposite to the foot on that side.
	dy := side * step * (0.047*legs - 0.034*arms)
	dy -= 0.009 * (1 - math.Cos(2*phase)) * w.Weight * body
	dx := 0.007*step*body + side*0.008*(1-math.Cos(phase))*w.Weight*legs
	return dx, dy
}
