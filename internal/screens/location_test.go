package screens

import (
	"game/internal/domain"
	"game/internal/ui"
	"testing"
)

func TestLocationFacesMovementAndKeepsDirectionAtRest(t *testing.T) {
	for _, tc := range []struct {
		name   string
		dx, dy float64
		want   int
	}{
		{"south", 0, 1, 0}, {"southeast", 1, 1, 1}, {"east", 1, 0, 2}, {"northeast", 1, -1, 3},
		{"north", 0, -1, 4}, {"northwest", -1, -1, 5}, {"west", -1, 0, 6}, {"southwest", -1, 1, 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := domain.NewArcanist()
			s := NewLocation(ui.Painter{}, nil, &c)
			s.Update(ui.Input{MoveX: tc.dx, MoveY: tc.dy})
			if s.facing != tc.want {
				t.Fatalf("facing %d, want %d", s.facing, tc.want)
			}
			for i := 0; i < 60; i++ {
				s.Update(ui.Input{})
			}
			if s.facing != tc.want {
				t.Fatal("direction reset while idle")
			}
		})
	}
}

func TestFacingAtBorderAndZeroSpeed(t *testing.T) {
	c := domain.NewArcanist()
	s := NewLocation(ui.Painter{}, nil, &c)
	s.x = worldWidth - 260
	s.Update(ui.Input{MoveX: 1, MoveY: -1})
	if s.facing != 4 {
		t.Fatal("should face north when sliding along east border")
	}
	s.Update(ui.Input{MoveX: 1})
	if s.facing != 4 {
		t.Fatal("blocked movement should preserve last facing")
	}
	c.Stats.MoveSpeedPercent = 0
	s.Update(ui.Input{MoveX: -1})
	if s.facing != 4 {
		t.Fatal("zero speed should not change facing")
	}
}
