package screens

import (
	"game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Action int

const (
	Stay Action = iota
	OpenCharacter
	OpenPlay
	OpenMenu
	Exit
)

type Screen interface {
	Update(ui.Input) Action
	Draw(*ebiten.Image)
}
