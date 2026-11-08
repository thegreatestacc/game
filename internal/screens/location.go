package screens

import (
	"fmt"
	"game/internal/animation"
	"game/internal/assets"
	"game/internal/domain"
	"game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const worldWidth, worldHeight = 1920.0, 1280.0

const wizardX, wizardY, talkRadius = 960.0, 640.0, 105.0

type LocationScreen struct {
	p        ui.Painter
	a        *assets.Assets
	c        *domain.Character
	x, y     float64
	back     ui.Button
	idle     animation.Idle
	facing   int // Atlas order: south, southeast, east, northeast, north, northwest, west, southwest.
	dialogue bool
}

func NewLocation(p ui.Painter, a *assets.Assets, c *domain.Character) *LocationScreen {
	return &LocationScreen{p: p, a: a, c: c, x: worldWidth / 2, y: worldHeight/2 + 160, back: ui.Button{X: 1080, Y: 24, W: 170, H: 44, Label: "В меню"}}
}

func (s *LocationScreen) Update(in ui.Input) Action {
	if s.dialogue {
		s.idle.Update(false)
		if in.Interact || in.Escape {
			s.dialogue = false
		}
		return Stay
	}
	if in.Escape || s.back.Update(in) {
		return OpenMenu
	}
	if in.Interact && s.nearWizard() {
		s.dialogue = true
		s.idle.Update(false)
		return Stay
	}
	dx, dy := in.MoveX, in.MoveY
	previousX, previousY := s.x, s.y
	if length := math.Hypot(dx, dy); length > 0 {
		// Fixed 60 Hz update; diagonal walking has the same speed.
		speed := 220.0 / 60 * math.Max(0, s.c.Stats.MoveSpeedPercent) / 100
		s.x += dx / length * speed
		s.y += dy / length * speed
	}
	// The dense border is decorative forest; walking stays in the clearing.
	s.x = math.Max(260, math.Min(worldWidth-260, s.x))
	s.y = math.Max(220, math.Min(worldHeight-220, s.y))
	if distance := math.Hypot(s.x-wizardX, s.y-wizardY); distance < 28 {
		if distance > 0 {
			s.x = wizardX + (s.x-wizardX)*28/distance
			s.y = wizardY + (s.y-wizardY)*28/distance
		} else {
			s.x = previousX
			s.y = previousY
		}
	}
	actualX, actualY := s.x-previousX, s.y-previousY
	moving := math.Hypot(actualX, actualY) > 0.0001
	if moving {
		// Screen Y grows downward. Quantize actual movement after collision
		// clamping, so sliding along a border also faces the correct direction.
		s.facing = (int(math.Round(math.Atan2(actualX, actualY)/(math.Pi/4))) + 8) % 8
	}
	s.idle.Update(moving)
	return Stay
}

func (s *LocationScreen) Draw(dst *ebiten.Image) {
	cx := math.Max(0, math.Min(worldWidth-1280, s.x-640))
	cy := math.Max(0, math.Min(worldHeight-720, s.y-360))
	op := &ebiten.DrawImageOptions{}
	b := s.a.Clearing.Bounds()
	op.GeoM.Scale(worldWidth/float64(b.Dx()), worldHeight/float64(b.Dy()))
	op.GeoM.Translate(-cx, -cy)
	dst.DrawImage(s.a.Clearing, op)
	x, y := s.x-cx, s.y-cy
	if s.y >= wizardY {
		s.drawWizard(dst, cx, cy)
	}
	vector.DrawFilledCircle(dst, float32(x), float32(y-3), 13, color.NRGBA{R: 12, G: 23, B: 16, A: 125}, true)
	character := s.a.Character
	if len(s.a.CharacterTurns) == 8 && s.a.CharacterTurns[s.facing] != nil {
		character = s.a.CharacterTurns[s.facing]
	}
	if character != nil {
		s.idle.Draw(dst, character, x, y, 86)
	} else {
		vector.DrawFilledCircle(dst, float32(x), float32(y-20), 16, ui.Gold, false)
	}
	if s.y < wizardY {
		s.drawWizard(dst, cx, cy)
	}
	if s.nearWizard() && !s.dialogue {
		x, y := wizardX-cx, wizardY-cy-119
		ui.Frame(dst, float32(x-86), float32(y), 172, 34)
		s.p.Center(dst, "нажмите E", x, y+6, 20, ui.Gold)
	}
	ui.Frame(dst, 24, 24, 370, 86)
	s.p.Text(dst, "ЛЕСНАЯ ПОЛЯНА", 42, 35, 22, ui.Gold)
	s.p.Text(dst, fmt.Sprintf("%s   HP %d/%d   MP %d/%d", s.c.Name, s.c.Stats.Health, s.c.Stats.MaxHealth, s.c.Stats.Mana, s.c.Stats.MaxMana), 42, 73, 16, ui.Ink)
	s.back.Draw(dst, s.p)
	ui.Frame(dst, 24, 660, 650, 40)
	s.p.Text(dst, "W A S D — движение    E — действие    Esc — главное меню", 40, 670, 16, ui.Ink)
	if s.dialogue {
		ui.Frame(dst, 200, 480, 880, 160)
		s.p.Text(dst, "Загадочный волшебник", 228, 497, 24, ui.Gold)
		s.p.Text(dst, "Приветствую, путник. Чтобы передвигаться, используй WASD.", 228, 541, 21, ui.Ink)
		s.p.Text(dst, "Для действия нажимай E, когда находишься рядом.", 228, 575, 21, ui.Ink)
		s.p.Text(dst, "E или Esc — закончить разговор", 228, 612, 16, ui.Muted)
	}
}

func (s *LocationScreen) nearWizard() bool { return math.Hypot(s.x-wizardX, s.y-wizardY) <= talkRadius }

func (s *LocationScreen) drawWizard(dst *ebiten.Image, cx, cy float64) {
	x, y := wizardX-cx, wizardY-cy
	vector.DrawFilledCircle(dst, float32(x), float32(y-3), 17, color.NRGBA{R: 12, G: 23, B: 16, A: 125}, true)
	if s.a.Wizard == nil {
		return
	}
	b := s.a.Wizard.Bounds()
	scale := 110.0 / float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x-float64(b.Dx())*scale/2, y-102)
	dst.DrawImage(s.a.Wizard, op)
}
