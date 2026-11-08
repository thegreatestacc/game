package animation

import (
	"math"
	"testing"
)

func TestIdleBreathAndWalking(t *testing.T) {
	a := Idle{Weight: 1, Tick: BreathPeriod / 2}
	dx, dy := a.offset(0.64, 0.42)
	if dx <= 0 || dy >= 0 {
		t.Fatal("chest should expand and rise")
	}
	dx, dy = a.offset(0.6, 0.9)
	if math.Abs(dx)+math.Abs(dy) > 0.00001 {
		t.Fatal("feet moved")
	}
	for i := 0; i < 13; i++ {
		a.Update(true)
	}
	if a.Weight != 0 {
		t.Fatal("breathing did not fade out during movement")
	}
	for i := 0; i < 13; i++ {
		a.Update(false)
	}
	if a.Weight != 1 {
		t.Fatal("breathing did not resume")
	}
	a.Tick = 0
	x0, y0 := a.offset(0.64, 0.42)
	for i := 0; i < BreathPeriod; i++ {
		a.Update(false)
	}
	x1, y1 := a.offset(0.64, 0.42)
	if x0 != x1 || y0 != y1 {
		t.Fatal("breathing loop has a seam")
	}
}
