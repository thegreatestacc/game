package assets

import "testing"

func TestEmbeddedTurntableLoads(t *testing.T) {
	a, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(a.CharacterTurns) != 8 {
		t.Fatalf("embedded atlas did not load: got %d frames", len(a.CharacterTurns))
	}
	if a.Character != a.CharacterTurns[0] {
		t.Fatal("profile should use the front-facing frame")
	}
	for i, frame := range a.CharacterTurns {
		b := frame.Bounds()
		if b.Min.X != 0 || b.Min.Y != 0 || b.Dx() < 400 || b.Dy() < 400 {
			t.Fatalf("frame %d has invalid bounds: %v", i, b)
		}
	}
}
