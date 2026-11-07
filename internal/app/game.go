package app

import (
	"game/internal/assets"
	"game/internal/domain"
	"game/internal/screens"
	"game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	screen    screens.Screen
	assets    *assets.Assets
	character domain.Character
	painter   ui.Painter
}

func New() (*Game, error) {
	a, err := assets.Load()
	if err != nil {
		return nil, err
	}
	p := ui.Painter{Font: a.Font}
	return &Game{screen: screens.NewMenu(p), assets: a, character: domain.NewArcanist(), painter: p}, nil
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	in := ui.Input{X: x, Y: y, Click: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft), Escape: inpututil.IsKeyJustPressed(ebiten.KeyEscape)}
	// Only the old screen receives this tick's input. A held mouse button
	// cannot trigger another action after switching screens.
	switch g.screen.Update(in) {
	case screens.OpenCharacter:
		g.screen = screens.NewCharacter(g.painter, g.assets, &g.character)
	case screens.OpenPlay:
		g.screen = screens.NewPlay(g.painter, g.assets, &g.character)
	case screens.OpenMenu:
		g.screen = screens.NewMenu(g.painter)
	case screens.Exit:
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(dst *ebiten.Image)     { g.screen.Draw(dst) }
func (g *Game) Layout(_, _ int) (int, int) { return 1280, 720 }
