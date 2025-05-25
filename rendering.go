package main

import (
	"image"
	"image/color"
	"math"
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

// draws the score in the top left corner
func DrawScore(g *Game, screen *ebiten.Image) {
	// Draw the score at the top left corner
	scoreText := "Score: " + strconv.Itoa(g.score)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.White)

	// Draw the bomb count below the score
	bombText := "Bombs: " + strconv.Itoa(g.bombs)
	text.Draw(screen, bombText, basicfont.Face7x13, 10, 40, color.RGBA{255, 200, 0, 255})
}

func DrawSplashScreen(g *Game, screen *ebiten.Image) {
	bounds := screen.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	msg := "SPACE SHOOTER\n\nControls:\n\nA/D - Rotate\nW/S - Accelerate/Decelerate\nSPACE - Shoot\nB - Bomb\n\nPress ENTER to start"
	lines := strings.Split(msg, "\n")
	y := 0

	if g.previousScore > 0 {
		msg = "GAME OVER\n\nPress ENTER to start again"
		lines = strings.Split(msg, "\n")

		scoreText := "Score: " + strconv.Itoa(g.previousScore)

		scale := 3.0 // 3x bigger
		face := basicfont.Face7x13
		b := text.BoundString(face, scoreText)
		scoreWidth := float64(b.Max.X-b.Min.X) * scale

		x := int(float64(w)/2 - (scoreWidth / 2))
		yScore := int(float64(h)/2 - float64(b.Max.Y-b.Min.Y)/2*scale)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(yScore)/scale)
		op.ColorScale.Scale(1, 1, 1, 1)
		text.DrawWithOptions(screen, scoreText, face, op)
	}
	y = h/2 - len(lines)*12
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

		if p.Type == powerupShield {

			shieldRect := image.Rect(0, 0, 32, 32) // x0, y0, x1, y1 in pixels
			shieldSprite := resources.TilesImage.SubImage(shieldRect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(shieldSprite, op)
		}

		if p.Type == powerupBomb {

			bombRect := image.Rect(32, 0, 64, 32) // x0, y0, x1, y1 in pixels
			bombSprite := resources.TilesImage.SubImage(bombRect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(bombSprite, op)
		}

		if p.Type == powerupInvincibleBullets {
			bulletRect := image.Rect(64, 0, 96, 32) // x0, y0, x1, y1 in pixels
			bulletSprite := resources.TilesImage.SubImage(bulletRect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(bulletSprite, op)
		}

		if p.Type == powerupFreezeEnemies {
			rect := image.Rect(0, 32, 32, 64) // x0, y0, x1, y1 in pixels
			sprite := resources.TilesImage.SubImage(rect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16) // Center the sprite
			screen.DrawImage(sprite, op)
		}
	}
}

func (a *Anomaly) DrawAnomaly(screen *ebiten.Image) {
	if !a.IsActive {
		return
	}
	w, h := screen.Size()

	// 1. Draw red overlay
	overlay := ebiten.NewImage(w, h)
	alpha := uint8(255 * a.Alpha)
	r := alpha * 255
	overlay.Fill(color.RGBA{r, 0, 0, a.Alpha})
	screen.DrawImage(overlay, nil)

	// 2. Draw fully opaque black circle to subtract from overlay
	// Use a smaller image as the mask
	diameter := int(a.SafeRadius * 2)
	mask := ebiten.NewImage(diameter, diameter)

	drawFilledCircle(mask, a.SafeRadius, a.SafeRadius, a.SafeRadius, color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(a.SafeX-a.SafeRadius, a.SafeY-a.SafeRadius)
	op.Blend = ebiten.BlendDestinationOut
	screen.DrawImage(mask, op)
}

func colorToFloats(c color.Color) (r, g, b, a float32) {
	cr, cg, cb, ca := c.RGBA()
	return float32(cr) / 65535, float32(cg) / 65535, float32(cb) / 65535, float32(ca) / 65535
}

func drawFilledCircle(dst *ebiten.Image, cx, cy, r float64, clr color.Color) {
	const steps = 64
	vertices := make([]ebiten.Vertex, 0, steps+1)
	indices := make([]uint16, 0, steps*3)

	cr, cg, cb, ca := colorToFloats(clr)

	// Center vertex
	vertices = append(vertices, ebiten.Vertex{
		DstX:   float32(cx),
		DstY:   float32(cy),
		ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca,
	})

	// Outer ring
	for i := 0; i <= steps; i++ {
		angle := 2 * math.Pi * float64(i) / float64(steps)
		x := float32(cx + r*math.Cos(angle))
		y := float32(cy + r*math.Sin(angle))
		vertices = append(vertices, ebiten.Vertex{
			DstX:   x,
			DstY:   y,
			ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca,
		})
	}

	// Indices for triangle fan
	for i := 1; i <= steps; i++ {
		indices = append(indices, 0, uint16(i), uint16(i+1))
	}

	// Dummy image (required by DrawTriangles)
	src := ebiten.NewImage(1, 1)
	src.Fill(color.White)

	dst.DrawTriangles(vertices, indices, src, nil)
}
