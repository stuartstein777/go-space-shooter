package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func FillScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the screen with black
}

func DrawShip(g *Game, screen *ebiten.Image) {
	// player is an elongated diamond shape
	cx := float64(g.playerLocation.X)
	cy := float64(g.playerLocation.Y)
	shipHeight := float64(75.0)
	shipWidth := float64(30.0)

	// Define the four points of the diamond
	topX, topY := cx, cy-shipHeight/2
	rightX, rightY := cx+shipWidth/2, cy
	bottomX, bottomY := cx, cy+shipHeight/4
	leftX, leftY := cx-shipWidth/2, cy

	angle := g.shipAngle
	topX, topY = RotatePoint(topX, topY, cx, cy, angle)
	rightX, rightY = RotatePoint(rightX, rightY, cx, cy, angle)
	bottomX, bottomY = RotatePoint(bottomX, bottomY, cx, cy, angle)
	leftX, leftY = RotatePoint(leftX, leftY, cx, cy, angle)

	vector.StrokeLine(screen, float32(topX), float32(topY), float32(rightX), float32(rightY), 2, color.White, true)
	vector.StrokeLine(screen, float32(rightX), float32(rightY), float32(bottomX), float32(bottomY), 2, color.White, true)
	vector.StrokeLine(screen, float32(bottomX), float32(bottomY), float32(leftX), float32(leftY), 2, color.White, true)
	vector.StrokeLine(screen, float32(leftX), float32(leftY), float32(topX), float32(topY), 2, color.White, true)
}

func DrawEnemies(g *Game, screen *ebiten.Image) {
	for _, e := range g.enemies {
		col := color.RGBA{255, 0, 0, 255}

		// flash them red if they are hit
		if e.HitTimer == 0 {
			col = color.RGBA{255, 255, 0, 255}
		}

		vector.StrokeCircle(screen, float32(e.X), float32(e.Y), float32(e.Radius), 2, col, false)
	}
}

func DrawBullets(g *Game, screen *ebiten.Image) {
	for _, b := range g.bullets {
		vector.DrawFilledCircle(screen, float32(b.X), float32(b.Y), 4, color.RGBA{R: 0, G: 255, B: 0, A: 100}, false)
	}
}

func DrawScore(g *Game, screen *ebiten.Image) {
	// Draw the score at the top left corner
	scoreText := "Score: " + strconv.Itoa(g.score)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.White)
}
