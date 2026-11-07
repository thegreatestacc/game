package animation

import (
	"image"
	"image/color"
	"testing"
)

func TestCorrespondenceTracksMovingSilhouette(t *testing.T) {
	makeImage := func(offset int) *image.NRGBA {
		img := image.NewNRGBA(image.Rect(0, 0, 128, 128))
		for y := 0; y < 128; y++ {
			for x := offset; x < offset+32; x++ {
				img.SetNRGBA(x, y, color.NRGBA{255, 255, 255, 255})
			}
		}
		return img
	}
	matched := matchRows(makeImage(24), makeImage(64))
	a, b := matched[0], matched[1]
	center := (rows/2)*(columns+1) + columns/2
	if b[center].x-a[center].x < 25 {
		t.Fatal("mesh did not follow moving silhouette")
	}
	for _, mesh := range [][]point{a, b} {
		for y := 0; y <= rows; y++ {
			for x := 1; x <= columns; x++ {
				i := y*(columns+1) + x
				if mesh[i].x <= mesh[i-1].x {
					t.Fatal("folded mesh")
				}
			}
		}
	}
}
