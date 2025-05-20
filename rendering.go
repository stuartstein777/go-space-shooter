package main

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

var whiteImg *ebiten.Image

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

	shipColour := color.RGBA{255, 255, 255, 255}

	if g.hasShield {
		shipColour = color.RGBA{0, 255, 255, 180}
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

func DrawSplashScreen(screen *ebiten.Image) {
	w, h := screen.Size()
	msg := "SPACE SHOOTER\n\nControls:\nA/D - Rotate\nW/S - Accelerate/Decelerate\nSPACE - Shoot\n\nPress SPACE to start"
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
			cx, cy := float32(p.X), float32(p.Y)
			size := float32(16)

			points := [][2]float32{
				{cx, cy - size},
				{cx - size*0.7, cy + size*0.3},
				{cx, cy + size},
				{cx + size*0.7, cy + size*0.3},
			}

			verts := []ebiten.Vertex{
				{DstX: points[0][0], DstY: points[0][1], ColorR: 0, ColorG: 0.78, ColorB: 1, ColorA: 0.7},
				{DstX: points[1][0], DstY: points[1][1], ColorR: 0, ColorG: 0.78, ColorB: 1, ColorA: 0.7},
				{DstX: points[2][0], DstY: points[2][1], ColorR: 0, ColorG: 0.78, ColorB: 1, ColorA: 0.7},
				{DstX: points[3][0], DstY: points[3][1], ColorR: 0, ColorG: 0.78, ColorB: 1, ColorA: 0.7},
			}
			indices := []uint16{
				0, 1, 2,
				0, 2, 3,
			}
			screen.DrawTriangles(verts, indices, whiteImg, nil)

			outlineColor := color.RGBA{0, 255, 255, 255}
			for i := 0; i < len(points); i++ {
				j := (i + 1) % len(points)
				vector.StrokeLine(screen,
					points[i][0], points[i][1],
					points[j][0], points[j][1],
					2, outlineColor, false)
			}
		}
	}
}
