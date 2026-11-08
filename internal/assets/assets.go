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
//go:embed fonts/* sprites/* locations/*
var Files embed.FS

type Assets struct {
	Wizard         *ebiten.Image
	Clearing       *ebiten.Image
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
	wizardData, err := Files.ReadFile("sprites/tutorial-wizard.png")
	if err != nil {
		return nil, err
	}
	wizardImage, _, err := image.Decode(bytes.NewReader(wizardData))
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить NPC: %w", err)
	}
	a.Wizard = ebiten.NewImageFromImage(wizardImage)
	ground, err := Files.ReadFile("locations/clearing.png")
	if err != nil {
		return nil, err
	}
	groundImage, _, err := image.Decode(bytes.NewReader(ground))
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить локацию: %w", err)
	}
	a.Clearing = ebiten.NewImageFromImage(groundImage)
	if data, err := Files.ReadFile("sprites/starter-pixel-turntable.png"); err == nil {
		if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
			a.Turntable = animation.NewTurntable(img)
			a.CharacterTurns = a.Turntable.Frames
			a.Character = a.CharacterTurns[0]
		}
	}
	// The starter has no equipment. All screens share the same starter sprite;
	// old armored artwork must not override this appearance.
	return a, nil
}
