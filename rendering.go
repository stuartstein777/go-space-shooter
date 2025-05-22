package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/stuartstein777/go-space-shooter/resources"
	"golang.org/x/image/font/basicfont"
)

var whiteImg *ebiten.Image

func FillScreen(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the screen with black
}

func DrawShip(g *Game, screen *ebiten.Image, isBlack bool) {
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

	shipColour := color.RGBA{255, 255, 255, 255}

	if isBlack {
		shipColour = color.RGBA{0, 0, 0, 255} // black
	}

	if g.hasShield {
		// Flash for last 2 seconds (120 frames)
		if g.shieldTimer <= 120 {
			// Alternate every 10 frames between white and cyan
			if (g.shieldTimer/10)%2 == 0 {
				shipColour = color.RGBA{0, 255, 255, 220} // cyan
			} else {
				shipColour = color.RGBA{255, 255, 255, 255} // white
			}
		} else {
			shipColour = color.RGBA{0, 255, 255, 180} // normal cyan
		}
	}

	vector.StrokeLine(screen, float32(topX), float32(topY), float32(rightX), float32(rightY), 2, shipColour, true)
	vector.StrokeLine(screen, float32(rightX), float32(rightY), float32(bottomX), float32(bottomY), 2, shipColour, true)
	vector.StrokeLine(screen, float32(bottomX), float32(bottomY), float32(leftX), float32(leftY), 2, shipColour, true)
	vector.StrokeLine(screen, float32(leftX), float32(leftY), float32(topX), float32(topY), 2, shipColour, true)
}

func DrawEnemies(g *Game, screen *ebiten.Image) {
	for _, e := range g.enemies {
		col := color.RGBA{255, 0, 0, 255}

		// flash them red if they are hit
		if e.HitTimer == 0 {
			col = color.RGBA{255, 255, 0, 255}
		}

		if e.IsInvincible {
			vector.DrawFilledCircle(screen, float32(e.X), float32(e.Y), float32(e.Radius), color.RGBA{255, 0, 0, 255}, false)
			vector.StrokeCircle(screen, float32(e.X), float32(e.Y), float32(e.Radius), 2, col, false)
		} else {
			vector.StrokeCircle(screen, float32(e.X), float32(e.Y), float32(e.Radius), 2, col, false)
		}
	}
}

func DrawBullets(g *Game, screen *ebiten.Image) {

	bulletColor := color.RGBA{0, 255, 0, 100}

	if g.invincibleBulletsTimer > 0 {
		bulletColor = color.RGBA{255, 0, 255, 100}
	}
	for _, b := range g.bullets {
		if b.Active {
			vector.DrawFilledCircle(screen, float32(b.X), float32(b.Y), 4, bulletColor, false)
		}
	}
}

func DrawScore(g *Game, screen *ebiten.Image) {
	// Draw the score at the top left corner
	scoreText := "Score: " + strconv.Itoa(g.score)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.White)

	// Draw the bomb count below the score
	bombText := "Bombs: " + strconv.Itoa(g.bombs)
	text.Draw(screen, bombText, basicfont.Face7x13, 10, 40, color.RGBA{255, 200, 0, 255})
}

func DrawSplashScreen(screen *ebiten.Image) {
	bounds := screen.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	msg := "SPACE SHOOTER\n\nControls:\n\nA/D - Rotate\nW/S - Accelerate/Decelerate\nSPACE - Shoot\nB - Bomb\n\nPress SPACE to start"
	lines := strings.Split(msg, "\n")
	y := h/2 - len(lines)*12
	for i, line := range lines {
		bounds := text.BoundString(basicfont.Face7x13, line)
		x := (w - bounds.Dx()) / 2
		text.Draw(screen, line, basicfont.Face7x13, x, y+20*i, color.White)
	}
}

func DrawPowerups(g *Game, screen *ebiten.Image) {
	for _, p := range g.powerups {
		if !p.Active {
			continue
		}

		if p.Type == "shield" {

			shieldRect := image.Rect(0, 0, 32, 32) // x0, y0, x1, y1 in pixels
			shieldSprite := resources.TilesImage.SubImage(shieldRect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(shieldSprite, op)
		}

		if p.Type == "bomb" {

			bombRect := image.Rect(32, 0, 64, 32) // x0, y0, x1, y1 in pixels
			bombSprite := resources.TilesImage.SubImage(bombRect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(bombSprite, op)
		}

		fmt.Println("Powerup type:", p.Type)
		if p.Type == "invincibleBullets" {
			bulletRect := image.Rect(64, 0, 96, 32) // x0, y0, x1, y1 in pixels
			bulletSprite := resources.TilesImage.SubImage(bulletRect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(bulletSprite, op)
		}
	}
}
