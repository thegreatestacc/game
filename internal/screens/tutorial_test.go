package screens

import (
	"game/internal/domain"
	"game/internal/ui"
	"testing"
)

func TestTutorialProximityAndDialogue(t *testing.T) {
	c := domain.NewArcanist()
	s := NewLocation(ui.Painter{}, nil, &c)
	s.Update(ui.Input{Interact: true})
	if s.dialogue || s.nearWizard() {
		t.Fatal("spawn should require approaching NPC")
	}
	s.x = wizardX
	s.y = wizardY + 80
	if !s.nearWizard() {
		t.Fatal("nearby prompt missing")
	}
	s.Update(ui.Input{Interact: true})
	if !s.dialogue {
		t.Fatal("E did not open dialogue")
	}
	x, y := s.x, s.y
	s.Update(ui.Input{MoveX: 1})
	if s.x != x || s.y != y || !s.dialogue {
		t.Fatal("dialogue must pause movement and persist")
	}
	if s.Update(ui.Input{Escape: true}) != Stay || s.dialogue {
		t.Fatal("Escape should close dialogue, not leave map")
	}
	s.Update(ui.Input{Interact: true})
	s.Update(ui.Input{Interact: true})
	if s.dialogue {
		t.Fatal("second E should close dialogue")
	}
	s.y = wizardY + talkRadius + 1
	if s.nearWizard() {
		t.Fatal("prompt remains out of range")
	}
}
