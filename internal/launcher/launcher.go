package launcher

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"game/internal/assets"
	"game/internal/buildinfo"
	"game/internal/ui"
	"game/internal/updater"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type event struct {
	message, version string
	fraction         float64
	done, failed     bool
}

type Launcher struct {
	p                     ui.Painter
	client                *updater.Client
	play, update          ui.Button
	ctx                   context.Context
	cancel                context.CancelFunc
	events                chan event
	status, version       string
	progress              float64
	busy, failed, started bool
}

func New() (*Launcher, error) {
	data, err := assets.Files.ReadFile("fonts/Go-Regular.ttf")
	if err != nil {
		return nil, err
	}
	font, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	root, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("не удалось найти папку данных пользователя: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Launcher{p: ui.Painter{Font: font}, client: updater.New(buildinfo.Repository, os.Getenv("SPELLGAME_GITHUB_TOKEN"), filepath.Join(root, "SpellGame")),
		ctx: ctx, cancel: cancel, events: make(chan event, 32), status: "Подготовка проверки версии…",
		play: ui.Button{X: 392, Y: 467, W: 232, H: 60, Label: "Играть"}, update: ui.Button{X: 648, Y: 467, W: 260, H: 60, Label: "Обновить игру"}}, nil
}

func (l *Launcher) Close() { l.cancel() }

func (l *Launcher) start(action string) {
	if l.busy {
		return
	}
	l.busy = true
	l.failed = false
	l.progress = 0
	l.status = "Проверка последней версии на GitHub…"
	go func() {
		send := func(e event) {
			select {
			case l.events <- e:
			case <-l.ctx.Done():
			}
		}
		fail := func(err error) { send(event{message: err.Error(), done: true, failed: true}) }
		available, err := l.client.Latest(l.ctx)
		if err != nil {
			fail(err)
			return
		}
		version := available.Manifest.Version
		send(event{message: "Доступна версия: " + version, version: version})
		if action == "check" {
			send(event{message: "Проверка успешна. Нажмите «Играть» или «Обновить игру».", done: true})
			return
		}
		path, err := l.client.Install(l.ctx, available, func(message string, fraction float64) {
			select {
			case l.events <- event{message: message, fraction: fraction}:
			default:
			}
		})
		if err != nil {
			fail(err)
			return
		}
		if action == "update" {
			send(event{message: "Игра обновлена до " + version + ". Можно играть.", fraction: 1, done: true})
			return
		}
		if l.ctx.Err() != nil {
			return
		}
		// An absolute path and working directory avoid dependence on where the
		// launcher was opened. No shell, arguments or executable names come from GitHub.
		cmd := exec.Command(path)
		cmd.Dir = filepath.Dir(path)
		// The game itself has no reason to inherit the launcher's GitHub token.
		for _, entry := range os.Environ() {
			if !strings.HasPrefix(strings.ToUpper(entry), "SPELLGAME_GITHUB_TOKEN=") {
				cmd.Env = append(cmd.Env, entry)
			}
		}
		if err := cmd.Start(); err != nil {
			fail(fmt.Errorf("не удалось запустить игру; повторите обновление и проверьте доступ к файлу"))
			return
		}
		send(event{message: "Игра запущена. Лаунчер ожидает её закрытия.", fraction: 1})
		err = cmd.Wait()
		if err != nil {
			fail(fmt.Errorf("игра завершилась с ошибкой; попробуйте обновить её"))
			return
		}
		send(event{message: "Игра закрыта. Перед следующим запуском версия будет проверена снова.", fraction: 1, done: true})
	}()
}

func (l *Launcher) Update() error {
	if !l.started {
		l.started = true
		l.start("check")
	}
	for {
		select {
		case e := <-l.events:
			l.status = e.message
			l.progress = e.fraction
			l.failed = e.failed
			if e.version != "" {
				l.version = e.version
			}
			if e.done {
				l.busy = false
			}
		default:
			goto received
		}
	}
received:
	x, y := ebiten.CursorPosition()
	in := ui.Input{X: x, Y: y, Click: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)}
	if l.busy {
		in.Click = false
		in.X = -1
		in.Y = -1
	}
	if l.play.Update(in) {
		l.start("play")
	}
	if l.update.Update(in) {
		l.start("update")
	}
	return nil
}

func (l *Launcher) Draw(dst *ebiten.Image) {
	ui.Backdrop(dst)
	l.p.Center(dst, "SpellGame", 640, 105, 62, ui.Ink)
	l.p.Center(dst, "ЛАУНЧЕР", 640, 185, 18, ui.Gold)
	ui.Frame(dst, 272, 250, 736, 176)
	title := "Проверка и обновление игры"
	if l.version != "" {
		title = "Последний релиз: " + l.version
	}
	l.p.Text(dst, title, 300, 272, 23, ui.Ink)
	l.wrapped(dst, l.status, 300, 315, 675, 20)
	ui.Box(dst, 300, 399, 680, 4, ui.Border)
	if l.progress > 0 {
		ui.Box(dst, 300, 399, float32(680*l.progress), 4, ui.Gold)
	}
	l.play.Draw(dst, l.p)
	l.update.Draw(dst, l.p)
	if l.busy {
		l.p.Center(dst, "Подождите завершения текущей операции", 640, 552, 17, ui.Muted)
	} else {
		l.p.Center(dst, "Перед каждым запуском проверяется версия и целостность игры", 640, 552, 17, ui.Muted)
	}
	l.p.Center(dst, "Приватный GitHub · "+buildinfo.Repository, 640, 604, 15, ui.Muted)
	l.p.Text(dst, "ЛАУНЧЕР "+buildinfo.Version, 48, 680, 14, ui.Muted)
}

func (l *Launcher) wrapped(dst *ebiten.Image, message string, x, y, width, size float64) {
	face := &text.GoTextFace{Source: l.p.Font, Size: size}
	line := ""
	lines := 0
	for _, word := range strings.Fields(message) {
		candidate := strings.TrimSpace(line + " " + word)
		w, _ := text.Measure(candidate, face, size)
		if w > width && line != "" {
			l.p.Text(dst, line, x, y, size, ui.Muted)
			line = word
			y += 27
			lines++
			if lines >= 2 {
				break
			}
		} else {
			line = candidate
		}
	}
	c := ui.Muted
	if l.failed {
		c = ui.Gold
	}
	l.p.Text(dst, line, x, y, size, c)
}

func (l *Launcher) Layout(int, int) (int, int) { return 1280, 720 }
