package resources

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadBackground() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(BackgroundPNG))
	if err != nil {
		panic(err)
	}
	BackgroundImage = ebiten.NewImageFromImage(img)
}

var (
	TilesImage *ebiten.Image
	//go:embed sprites.png
	Tiles_png []byte
	//go:embed starfield.png
	BackgroundPNG   []byte
	BackgroundImage *ebiten.Image
)
