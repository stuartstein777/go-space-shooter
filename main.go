package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	highDPIImageCh = make(chan *ebiten.Image)
)

type Point struct {
	X int
	Y int
}

type Game struct {
	playerLocation  Point
	playerDirection int     // 0 == up, 180 = down, 90 = right, 270 = left
	shipAngle       float64 // in radians
	velocity        float64
	maxSpeed        float64
}

func rotatePoint(x, y, cx, cy, angle float64) (float64, float64) {
	sin, cos := math.Sin(float64(angle)), math.Cos(float64(angle))
	dx, dy := float64(x-cx), float64(y-cy)
	return cx + dx*cos - dy*sin, cy + dx*sin + dy*cos
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the player
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the screen with black

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
	topX, topY = rotatePoint(topX, topY, cx, cy, angle)
	rightX, rightY = rotatePoint(rightX, rightY, cx, cy, angle)
	bottomX, bottomY = rotatePoint(bottomX, bottomY, cx, cy, angle)
	leftX, leftY = rotatePoint(leftX, leftY, cx, cy, angle)

	// Draw the ship
	vector.StrokeLine(screen, float32(topX), float32(topY), float32(rightX), float32(rightY), 2, color.White, true)
	vector.StrokeLine(screen, float32(rightX), float32(rightY), float32(bottomX), float32(bottomY), 2, color.White, true)
	vector.StrokeLine(screen, float32(bottomX), float32(bottomY), float32(leftX), float32(leftY), 2, color.White, true)
	vector.StrokeLine(screen, float32(leftX), float32(leftY), float32(topX), float32(topY), 2, color.White, true)
}

func (g *Game) Reset() {
	g.playerLocation = Point{X: 640, Y: 480}
	g.velocity = 0
	g.maxSpeed = 20 // adjust as desired
}

func (g *Game) Update() error {
	const rotateSpeed = 0.08 // radians per frame
	const accel = 0.2        // acceleration per frame
	const friction = 0.01    // natural slow down
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.shipAngle -= rotateSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.shipAngle += rotateSpeed
	}

	// Acceleration/Deceleration
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.velocity += accel
		if g.velocity > g.maxSpeed {
			g.velocity = g.maxSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.velocity -= accel
		if g.velocity < 0 {
			g.velocity = 0
		}
	}

	// Apply friction if not accelerating
	if !ebiten.IsKeyPressed(ebiten.KeyW) && !ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.velocity > 0 {
			g.velocity -= friction
			if g.velocity < 0 {
				g.velocity = 0
			}
		}
	}

	// Move ship forward in the direction it's facing
	g.playerLocation.X += int(g.velocity * math.Sin(g.shipAngle))
	g.playerLocation.Y -= int(g.velocity * math.Cos(g.shipAngle))

	// Screen wrapping
	screenWidth, screenHeight := g.Layout(0, 0)
	if g.playerLocation.X < 0 {
		g.playerLocation.X = screenWidth - 1
	}
	if g.playerLocation.X >= screenWidth {
		g.playerLocation.X = 0
	}
	if g.playerLocation.Y < 0 {
		g.playerLocation.Y = screenHeight - 1
	}
	if g.playerLocation.Y >= screenHeight {
		g.playerLocation.Y = 0
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 960
}

func (g *Game) MovePlayer(direction string) {

}

func main() {
	game := &Game{}
	game.Reset()
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
