package screens

import (
	"fmt"
	"math"

	"game/internal/assets"
	"game/internal/domain"
	"game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

// PlayScreen reads the same character that is edited in CharacterScreen.
type PlayScreen struct {
	p            ui.Painter
	character    *domain.Character
	back         ui.Button
	rotationTick int
	assets       *assets.Assets
}

// Eight seconds per revolution at the game's fixed 60 updates per second.
const rotationPeriod = 480

func NewPlay(p ui.Painter, a *assets.Assets, c *domain.Character) *PlayScreen {
	return &PlayScreen{p: p, assets: a, character: c, back: ui.Button{X: 48, Y: 36, W: 156, H: 48, Label: "← Назад"}}
}

func (s *PlayScreen) Update(in ui.Input) Action {
	if s.back.Update(in) || in.Escape {
		return OpenMenu
	}
	s.rotationTick = (s.rotationTick + 1) % rotationPeriod
	return Stay
}

func (s *PlayScreen) Draw(dst *ebiten.Image) {
	ui.Backdrop(dst)
	s.back.Draw(dst, s.p)
	s.p.Center(dst, "ПЕРСОНАЖ И ИНВЕНТАРЬ", 640, 44, 26, ui.Ink)
	ui.Frame(dst, 48, 120, 292, 510)
	ui.Frame(dst, 356, 120, 492, 510)
	ui.Frame(dst, 864, 120, 368, 510)
	s.p.Center(dst, s.character.Name, 194, 145, 29, ui.Ink)
	s.p.Center(dst, fmt.Sprintf("%s · Уровень %d", s.character.Class, s.character.Level), 194, 187, 18, ui.Gold)
	ui.Sigil(dst, 194, 370, 114)
	if s.assets != nil && s.assets.Turntable != nil {
		s.assets.Turntable.Draw(dst, 194, 550, 330, float64(s.rotationTick)/rotationPeriod)
	} else {
		ui.DrawTurningMage(dst, 194, 510, float64(s.rotationTick)*2*math.Pi/rotationPeriod)
	}
	s.p.Center(dst, "Обзор персонажа · 360°", 194, 576, 17, ui.Muted)
	s.p.Text(dst, "Характеристики", 376, 145, 26, ui.Ink)
	s.p.Text(dst, "Текущие значения", 376, 187, 17, ui.Gold)
	for i, row := range statRows(s.character.Stats) {
		y := float32(234 + i*41)
		if i%2 == 0 {
			ui.Box(dst, 368, y, 468, 39, ui.Background)
		}
		s.p.Text(dst, row[0], 380, float64(y+9), 18, ui.Muted)
		s.p.Center(dst, row[1], 772, float64(y+9), 18, ui.Ink)
	}
	s.p.Text(dst, "Инвентарь", 884, 145, 26, ui.Ink)
	occupied := 0
	for _, item := range s.character.Inventory {
		if item != nil {
			occupied++
		}
	}
	s.p.Text(dst, fmt.Sprintf("Занято ячеек: %d / %d", occupied, len(s.character.Inventory)), 884, 187, 17, ui.Gold)
	for i, item := range s.character.Inventory {
		x, y := float32(884+(i%3)*110), float32(234+(i/3)*80)
		ui.Frame(dst, x, y, 104, 72)
		s.p.Text(dst, fmt.Sprintf("%02d", i+1), float64(x+10), float64(y+8), 13, ui.Muted)
		if item == nil {
			s.p.Center(dst, "—", float64(x+52), float64(y+28), 24, ui.Border)
		} else {
			// Keep long future item names inside their slot.
			name := []rune(item.Name)
			if len(name) > 9 {
				name = append(name[:8], '…')
			}
			s.p.Center(dst, string(name), float64(x+52), float64(y+28), 14, ui.Ink)
			s.p.Center(dst, fmt.Sprintf("×%d", item.Quantity), float64(x+52), float64(y+49), 13, ui.Gold)
		}
	}
	if occupied == 0 {
		s.p.Center(dst, "Инвентарь пока пуст", 1048, 576, 18, ui.Muted)
	}
	s.p.Text(dst, "ESC  /  Вернуться в меню", 48, 680, 14, ui.Muted)
	s.p.Text(dst, "Характеристики меняются в меню «Персонаж»", 858, 680, 14, ui.Muted)
}
