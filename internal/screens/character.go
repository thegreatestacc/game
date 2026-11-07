package screens

import (
	"fmt"
	"image/color"
	"math"

	"game/internal/assets"
	"game/internal/domain"
	"game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CharacterScreen struct {
	p           ui.Painter
	a           *assets.Assets
	character   *domain.Character
	back        ui.Button
	minus, plus [domain.StatCount]ui.Button
}

func NewCharacter(p ui.Painter, a *assets.Assets, c *domain.Character) *CharacterScreen {
	s := &CharacterScreen{p: p, a: a, character: c, back: ui.Button{X: 48, Y: 36, W: 156, H: 48, Label: "← Назад"}}
	for i := range s.minus {
		y := 222 + i*45
		s.minus[i] = ui.Button{X: 966, Y: y, W: 34, H: 34, Label: "−"}
		s.plus[i] = ui.Button{X: 1186, Y: y, W: 34, H: 34, Label: "+"}
	}
	return s
}

func (s *CharacterScreen) Update(in ui.Input) Action {
	if s.back.Update(in) || in.Escape {
		return OpenMenu
	}
	for i := range s.minus {
		if s.minus[i].Update(in) {
			s.character.Stats.Adjust(domain.Stat(i), -1)
		}
		if s.plus[i].Update(in) {
			s.character.Stats.Adjust(domain.Stat(i), 1)
		}
	}
	return Stay
}

func (s *CharacterScreen) Draw(dst *ebiten.Image) {
	ui.Backdrop(dst)
	s.back.Draw(dst, s.p)
	s.p.Center(dst, "ПЕРСОНАЖ", 640, 44, 26, ui.Ink)
	ui.Frame(dst, 48, 120, 370, 510)
	ui.Sigil(dst, 233, 342, 140)
	s.p.Center(dst, "I  /  АРКАННЫЙ ОРДЕН", 233, 144, 16, ui.Gold)
	if s.a.Character != nil {
		b := s.a.Character.Bounds()
		scale := math.Min(300/float64(b.Dx()), 350/float64(b.Dy()))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(233-float64(b.Dx())*scale/2, 370-float64(b.Dy())*scale/2)
		dst.DrawImage(s.a.Character, op)
	} else {
		s.portrait(dst)
	}
	s.p.Center(dst, s.character.Name, 233, 553, 28, ui.Ink)
	s.p.Center(dst, fmt.Sprintf("%s • Уровень %d", s.character.Class, s.character.Level), 233, 593, 17, ui.Muted)
	s.p.Text(dst, s.character.Name, 466, 119, 36, ui.Ink)
	s.p.Text(dst, fmt.Sprintf("Класс: %s    /    Уровень: %d", s.character.Class, s.character.Level), 468, 170, 20, ui.Gold)
	rows := statRows(s.character.Stats)
	for i, row := range rows {
		y := float32(218 + i*45)
		if i%2 == 0 {
			ui.Box(dst, 466, y, 766, 43, ui.Panel)
		}
		c := ui.Ink
		if i == 0 {
			c = color.NRGBA{R: 148, G: 210, B: 163, A: 255}
		}
		if i == 1 {
			c = ui.Cyan
		}
		s.p.Text(dst, row[0], 484, float64(y+9), 21, ui.Muted)
		s.p.Center(dst, row[1], 1093, float64(y+9), 21, c)
		s.minus[i].Draw(dst, s.p)
		s.plus[i].Draw(dst, s.p)
	}
	s.p.Text(dst, "ESC  /  Вернуться в меню", 48, 680, 14, ui.Muted)
	s.p.Text(dst, "− / +  Шаг 1 · Изменения до закрытия игры", 882, 680, 14, ui.Muted)
}

// A deliberately stylized placeholder; no external art is required.
func (s *CharacterScreen) portrait(dst *ebiten.Image) {
	vector.DrawFilledCircle(dst, 233, 350, 105, color.NRGBA{R: 27, G: 43, B: 57, A: 255}, true)
	// Robe drawn as horizontal strips forms a tapered silhouette.
	for y := 330; y < 510; y++ {
		w := float32(45) + float32(y-330)*0.30
		ui.Box(dst, 233-w, float32(y), w*2, 1, color.NRGBA{R: 57, G: 76, B: 104, A: 255})
	}
	vector.DrawFilledCircle(dst, 233, 310, 49, ui.Border, true)
	vector.DrawFilledCircle(dst, 233, 314, 32, ui.Background, true)
	vector.StrokeLine(dst, 216, 314, 225, 314, 3, ui.Cyan, true)
	vector.StrokeLine(dst, 241, 314, 250, 314, 3, ui.Cyan, true)
	vector.StrokeLine(dst, 233, 363, 233, 493, 3, ui.Gold, true)
	vector.StrokeLine(dst, 154, 357, 142, 510, 7, ui.Gold, true)
	vector.DrawFilledCircle(dst, 154, 342, 18, ui.Cyan, true)
	vector.StrokeCircle(dst, 154, 342, 25, 1, ui.Cyan, true)
	ui.Box(dst, 184, 421, 98, 7, ui.Gold)
}
