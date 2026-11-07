package assets

import (
	"bytes"
	"embed"
	"fmt"
	"game/internal/animation"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Files includes the font license in the executable as well as the font.
//
//go:embed fonts/* sprites/*
var Files embed.FS

type Assets struct {
	Font           *text.GoTextFaceSource
	Character      *ebiten.Image
	CharacterTurns []*ebiten.Image
	Turntable      *animation.Turntable
}

func Load() (*Assets, error) {
	data, err := Files.ReadFile("fonts/Go-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить встроенный шрифт: %w", err)
	}
	font, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать встроенный шрифт: %w", err)
	}
	a := &Assets{Font: font}
	if data, err := Files.ReadFile("sprites/arcanist-turntable.png"); err == nil {
		if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
			a.Turntable = animation.NewTurntable(img)
			a.CharacterTurns = a.Turntable.Frames
			a.Character = a.CharacterTurns[0]
		}
	}
	// A missing or invalid optional sprite uses the procedural portrait.
	if data, err := Files.ReadFile("sprites/character.png"); err == nil {
		if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
			a.Character = ebiten.NewImageFromImage(img)
		}
	}
	return a, nil
}
