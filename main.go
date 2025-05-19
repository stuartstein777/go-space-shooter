package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {

	FillScreen(screen)
	DrawShip(g, screen)
	DrawEnemies(g, screen)
	DrawBullets(g, screen)
}

func (g *Game) Reset() {
	g.playerLocation = Point{X: 640, Y: 480}
	g.velocity = 0
	g.maxSpeed = 20 // adjust as desired
}

func (g *Game) HandleKeyPresses() {

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

	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.shootCooldown == 0 {
		// Calculate tip of ship (longest section)
		cx := float64(g.playerLocation.X)
		cy := float64(g.playerLocation.Y)
		shipLength := 40.0 // Should match your ship's tip length
		tipX := cx + shipLength*math.Sin(g.shipAngle)
		tipY := cy - shipLength*math.Cos(g.shipAngle)

		bulletSpeed := 10.0
		bullet := &Bullet{
			X:      tipX,
			Y:      tipY,
			VX:     bulletSpeed * math.Sin(g.shipAngle),
			VY:     -bulletSpeed * math.Cos(g.shipAngle),
			Active: true,
		}
		g.bullets = append(g.bullets, bullet)
		g.shootCooldown = 10 // frames between shots
	}
}

func movePlayerShip(g *Game) {
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
}

func (g *Game) Update() error {

	g.HandleKeyPresses()

	// Apply friction if not accelerating
	if !ebiten.IsKeyPressed(ebiten.KeyW) && !ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.velocity > 0 {
			g.velocity -= friction
			if g.velocity < 0 {
				g.velocity = 0
			}
		}
	}

	movePlayerShip(g)
	SpawnEnemies(g)
	DeSpawnEnemies(g)

	if g.shootCooldown > 0 {
		g.shootCooldown--
	}
	// Move bullets and remove inactive/out-of-bounds ones
	screenWidth, screenHeight := g.Layout(0, 0)
	activeBullets := g.bullets[:0]
	for _, b := range g.bullets {
		b.X += b.VX
		b.Y += b.VY
		if b.X < 0 || b.X > float64(screenWidth) || b.Y < 0 || b.Y > float64(screenHeight) {
			continue
		}
		activeBullets = append(activeBullets, b)
	}
	g.bullets = activeBullets

	return nil
}

func SpawnEnemies(g *Game) {
	// Randomly spawn an enemy every ~60 frames (1 second at 60fps)
	if rand.Float64() < 1.0/60.0 {
		screenWidth, screenHeight := g.Layout(0, 0)
		spawnX, spawnY, targetX, targetY := randomEdgeLocation(screenWidth, screenHeight)
		radius := 20.0

		// Calculate normalized velocity vector
		dx := float64(targetX - spawnX)
		dy := float64(targetY - spawnY)
		dist := math.Hypot(dx, dy)
		speed := 3.0 // pixels per frame
		vx := dx / dist * speed
		vy := dy / dist * speed

		enemy := &Enemy{
			X:      float64(spawnX),
			Y:      float64(spawnY),
			VX:     vx,
			VY:     vy,
			Radius: radius,
			Active: true,
		}
		g.enemies = append(g.enemies, enemy)
	}
}

func DeSpawnEnemies(g *Game) {
	screenWidth, screenHeight := g.Layout(0, 0)
	activeEnemies := g.enemies[:0]
	for _, e := range g.enemies {
		e.X += e.VX
		e.Y += e.VY
		// Remove if out of bounds
		if e.X+e.Radius < 0 || e.X-e.Radius > float64(screenWidth) ||
			e.Y+e.Radius < 0 || e.Y-e.Radius > float64(screenHeight) {
			continue
		}
		activeEnemies = append(activeEnemies, e)
	}
	g.enemies = activeEnemies
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 960
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
