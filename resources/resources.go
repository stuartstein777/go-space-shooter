package resources

import (
	_ "embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	TilesImage *ebiten.Image
	//go:embed sprites.png
	Tiles_png []byte
)
