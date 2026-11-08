package animation

import "testing"

func TestWalkingAlternatesLegsAndArms(t *testing.T) {
	w := Walk{Tick: WalkPeriod / 4, Weight: 1}
	_, left := w.offset(0.42, 0.9)
	_, right := w.offset(0.58, 0.9)
	if left >= 0 || right <= 0 {
		t.Fatal("feet should alternate")
	}
	_, hand := w.offset(0.68, 0.52)
	if hand >= 0 {
		t.Fatal("arm should swing opposite the right foot")
	}
	for i := 0; i < 12; i++ {
		w.Update(false)
	}
	if w.Weight != 0 || w.Tick != 0 {
		t.Fatal("walking must settle after stopping")
	}
	dx, dy := w.offset(0.58, 0.9)
	if dx != 0 || dy != 0 {
		t.Fatal("stopped walk still deforms sprite")
	}
	w.Update(true)
	if w.Weight <= 0 || w.Weight >= 1 {
		t.Fatal("walking should fade in")
	}
}
