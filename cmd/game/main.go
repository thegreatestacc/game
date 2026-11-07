package main

import (
	"flag"
	"fmt"
	"os"

	"game/internal/app"
	"game/internal/assets"
	"game/internal/buildinfo"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	license := flag.Bool("font-license", false, "print the embedded font license")
	flag.Parse()
	if *license {
		data, _ := assets.Files.ReadFile("fonts/LICENSE.txt")
		fmt.Print(string(data))
		return
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "SpellGame:", err)
		showError(err.Error())
		os.Exit(1)
	}
}

func run() error {
	g, err := app.New()
	if err != nil {
		return err
	}
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("SpellGame — " + buildinfo.Version)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	return ebiten.RunGame(g)
}
