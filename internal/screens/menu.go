package screens

import (
	"game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScreen struct {
	p                     ui.Painter
	play, character, exit ui.Button
}

func NewMenu(p ui.Painter) *MenuScreen {
	return &MenuScreen{p: p,
		play:      ui.Button{X: 480, Y: 345, W: 320, H: 60, Label: "Играть"},
		character: ui.Button{X: 480, Y: 425, W: 320, H: 60, Label: "Персонаж"},
		exit:      ui.Button{X: 480, Y: 505, W: 320, H: 60, Label: "Выход"}}
}

func (s *MenuScreen) Update(in ui.Input) Action {
	if s.play.Update(in) {
		return OpenPlay
	}
	if s.character.Update(in) {
		return OpenCharacter
	}
	if s.exit.Update(in) {
		return Exit
	}
	return Stay
}

func (s *MenuScreen) Draw(dst *ebiten.Image) {
	ui.Backdrop(dst)
	ui.Sigil(dst, 640, 185, 80)
	s.p.Center(dst, "S", 640, 140, 70, ui.Gold)
	s.p.Center(dst, "SpellGame", 640, 275, 58, ui.Ink)
	s.play.Draw(dst, s.p)
	s.character.Draw(dst, s.p)
	s.exit.Draw(dst, s.p)
	s.p.Center(dst, "Искусство начинается с искры", 640, 595, 19, ui.Muted)
	s.p.Text(dst, "SPELLGAME  /  MVP 0.1", 48, 680, 14, ui.Muted)
	s.p.Text(dst, "АРХИВ АРКАНИСТА", 1030, 680, 14, ui.Muted)
}
