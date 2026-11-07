package screens

import (
	"game/internal/domain"
	"game/internal/ui"
	"testing"
)

func TestNavigation(t *testing.T) {
	p := ui.Painter{}
	menu := NewMenu(p)
	if got := menu.Update(ui.Input{X: 600, Y: 390}); got != Stay {
		t.Fatal("hover activated menu", got)
	}
	if got := menu.Update(ui.Input{X: 600, Y: 450, Click: true}); got != OpenCharacter {
		t.Fatal("character click", got)
	}
	c := domain.NewArcanist()
	character := NewCharacter(p, nil, &c)
	if got := character.Update(ui.Input{X: 600, Y: 390}); got != Stay {
		t.Fatal("transition input leaked", got)
	}
	if got := character.Update(ui.Input{Escape: true}); got != OpenMenu {
		t.Fatal("escape", got)
	}
	if got := character.Update(ui.Input{X: 100, Y: 60, Click: true}); got != OpenMenu {
		t.Fatal("back click", got)
	}
	if got := menu.Update(ui.Input{X: 600, Y: 530, Click: true}); got != Exit {
		t.Fatal("exit click", got)
	}
	if got := menu.Update(ui.Input{X: 10, Y: 10, Click: true}); got != Stay {
		t.Fatal("outside click", got)
	}
}

func TestStatEditingSurvivesScreenRecreation(t *testing.T) {
	c := domain.NewArcanist()
	s := NewCharacter(ui.Painter{}, nil, &c)
	// A click increments once; subsequent held/hover ticks carry no Click edge.
	s.Update(ui.Input{X: 1200, Y: 320, Click: true})
	for i := 0; i < 60; i++ {
		s.Update(ui.Input{X: 1200, Y: 320})
	}
	if c.Stats.Strength != 6 {
		t.Fatalf("single click: %d", c.Stats.Strength)
	}
	if s.Update(ui.Input{Escape: true}) != OpenMenu {
		t.Fatal("escape")
	}
	s = NewCharacter(ui.Painter{}, nil, &c)
	if s.character.Stats.Strength != 6 {
		t.Fatal("edit lost on reopening")
	}
	s.Update(ui.Input{X: 980, Y: 320, Click: true})
	if c.Stats.Strength != 5 {
		t.Fatal("minus did not change shared character")
	}
}

func TestPlayNavigationAndCurrentCharacter(t *testing.T) {
	p := ui.Painter{}
	c := domain.NewArcanist()
	editor := NewCharacter(p, nil, &c)
	editor.Update(ui.Input{X: 1200, Y: 320, Click: true})
	menu := NewMenu(p)
	if got := menu.Update(ui.Input{X: 600, Y: 370}); got != Stay {
		t.Fatal("hover opened play")
	}
	if got := menu.Update(ui.Input{X: 600, Y: 370, Click: true}); got != OpenPlay {
		t.Fatalf("play click: %v", got)
	}
	play := NewPlay(p, nil, &c)
	if got := statRows(play.character.Stats)[2][1]; got != "6" {
		t.Fatalf("stale strength: %s", got)
	}
	for _, item := range play.character.Inventory {
		if item != nil {
			t.Fatal("new inventory should be empty")
		}
	}
	if got := play.Update(ui.Input{X: 600, Y: 370}); got != Stay {
		t.Fatal("transition input leaked")
	}
	if got := play.Update(ui.Input{Escape: true}); got != OpenMenu {
		t.Fatal("escape from play")
	}
	if got := play.Update(ui.Input{X: 100, Y: 60, Click: true}); got != OpenMenu {
		t.Fatal("back from play")
	}
	// Reopening the screen must retain the character's inventory as well as stats.
	c.Inventory[0] = &domain.Item{Name: "Посох", Quantity: 1}
	play = NewPlay(p, nil, &c)
	if play.character.Inventory[0] != c.Inventory[0] {
		t.Fatal("inventory lost on reopening")
	}
}

func TestPortraitRotationLoopsWithoutChangingCharacter(t *testing.T) {
	c := domain.NewArcanist()
	original := c
	s := NewPlay(ui.Painter{}, nil, &c)
	for i := 0; i < rotationPeriod/4; i++ {
		s.Update(ui.Input{})
	}
	if s.rotationTick != rotationPeriod/4 {
		t.Fatal("rotation did not advance")
	}
	for i := rotationPeriod / 4; i < rotationPeriod; i++ {
		s.Update(ui.Input{})
	}
	if s.rotationTick != 0 {
		t.Fatal("rotation did not loop")
	}
	if c != original {
		t.Fatal("animation modified character data")
	}
	if s.Update(ui.Input{Escape: true}) != OpenMenu {
		t.Fatal("cannot leave animated screen")
	}
}
