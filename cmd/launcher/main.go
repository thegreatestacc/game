package main

import (
	"fmt"
	"game/internal/launcher"
	"github.com/hajimehoshi/ebiten/v2"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		showError(err.Error())
		os.Exit(1)
	}
}

func run() error {
	l, err := launcher.New()
	if err != nil {
		return err
	}
	defer l.Close()
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("SpellGame — Лаунчер")
	ebiten.SetTPS(60)
	ebiten.SetRunnableOnUnfocused(true)
	return ebiten.RunGame(l)
}
