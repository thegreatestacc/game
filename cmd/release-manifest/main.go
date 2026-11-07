// release-manifest is a build utility, never shipped to players.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"game/internal/updater"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("usage: release-manifest VERSION GAME_EXE OUTPUT_JSON")
	}
	f, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha256.New()
	size, err := io.Copy(h, f)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(updater.Manifest{Version: os.Args[1], Asset: updater.GameAsset, SHA256: hex.EncodeToString(h.Sum(nil)), Size: size}, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(os.Args[3], data, 0644); err != nil {
		log.Fatal(err)
	}
}
