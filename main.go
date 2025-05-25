package main

import (
	"bytes"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"math/rand"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/stuartstein777/go-space-shooter/resources"
)

var bigFont font.Face

func (g *Game) Draw(screen *ebiten.Image) {

	if resources.BackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0) // Top-left corner
		screen.DrawImage(resources.BackgroundImage, op)
	}

	if g.showSplash {
		DrawSplashScreen(g, screen)
		return
	}

	if g.flashTimer > 0 {
		screen.Fill(color.White)
		DrawShip(g, screen, true)
		return
	}

	if g.Anomaly.Incoming > 0 && g.Anomaly.Incoming%5 != 0 {
		msg := "ANOMALY INCOMING!"
		bounds := text.BoundString(bigFont, msg)
		x := (1280 - bounds.Dx()) / 2
		y := 120 // Near the top

		text.Draw(screen, msg, bigFont, x, y, color.RGBA{255, 80, 80, 255})
	}

	g.Anomaly.DrawAnomaly(screen)
	DrawShip(g, screen, false)
	DrawEnemies(g, screen)
	DrawBullets(g, screen)
	DrawPowerups(g, screen)
	DrawScore(g, screen)

}

func loadResources() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(resources.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	resources.TilesImage = ebiten.NewImageFromImage(img)

	// Load the background image
	img, _, err = image.Decode(bytes.NewReader(resources.BackgroundPNG))
	if err != nil {
		log.Fatal(err)
	}
	resources.BackgroundImage = ebiten.NewImageFromImage(img)

	fontBytes, err := ioutil.ReadFile("resources/Roboto-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}
	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	bigFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    48, // Make this as big as you want
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Reset() {
	g.playerLocation = Point{X: 640, Y: 480}
	g.velocity = 0
	g.maxSpeed = 20 // adjust as desired
	//g.score = 0
	g.enemies = make([]*Enemy, 0)
	g.bullets = make([]*Bullet, 0)
	g.shootCooldown = bulletCooldown
	g.showSplash = true
	g.hasShield = false
	g.powerups = make([]*Powerup, 0)
	g.invincibleBulletsTimer = 0
	g.frozenEnemiesTimer = 0
	g.shieldTimer = 0
	g.shipAngle = 0
	g.flashTimer = 0
	g.bombs = 0
	g.score = 0
	whiteImg = ebiten.NewImage(1, 1)
	whiteImg.Fill(color.White)
	g.Anomaly.Deactivate()
}

func (g *Game) HandleKeyPresses() {

	if ebiten.IsKeyPressed(ebiten.KeyB) && g.bombs > 0 && g.flashTimer == 0 {
		g.bombs--
		g.flashTimer = 20 // flash for 20 frames (~1/3 second at 60fps)

		// Kill all enemies
		for _, e := range g.enemies {
			e.Active = false
			g.score += getScore(int(e.Radius))
		}

		// Immediately remove inactive enemies
		activeEnemies := g.enemies[:0]
		for _, e := range g.enemies {
			if e.Active {
				activeEnemies = append(activeEnemies, e)
			}
		}

		g.enemies = activeEnemies
	}

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

		cx := float64(g.playerLocation.X)
		cy := float64(g.playerLocation.Y)
		shipLength := 40.0
		tipX := cx + shipLength*math.Sin(g.shipAngle)
		tipY := cy - shipLength*math.Cos(g.shipAngle)

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

	g.Anomaly.Update()

	// if they have a shield, reduce the timer on it.
	if g.hasShield {
		g.shieldTimer--
		if g.shieldTimer <= 0 {
			g.hasShield = false
		}
	}

	if g.Anomaly.IsActive {
		g.Anomaly.fadeTimer--
		if g.Anomaly.fadeTimer <= 0 {
			g.Anomaly.IsActive = false
		}
	}

	if g.Anomaly.IsActive && g.Anomaly.fadeTimer == 1 { // On anomaly "strike"
		dx := float64(g.playerLocation.X) - g.Anomaly.SafeX
		dy := float64(g.playerLocation.Y) - g.Anomaly.SafeY
		if dx*dx+dy*dy > g.Anomaly.SafeRadius*g.Anomaly.SafeRadius {
			//	g.flashTimer = 20 // flash for 20 frames (~1/3 second at 60fps)
			g.previousScore = g.score
			g.showSplash = true
			g.score = 0
		}
	}

	if g.frozenEnemiesTimer > 0 {
		g.frozenEnemiesTimer--
	}

	if g.invincibleBulletsTimer > 0 {
		g.invincibleBulletsTimer--
	}

	if g.flashTimer > 0 {
		g.flashTimer--
	}

	// Handle pressing enter on the splash screen to start the game
	if g.showSplash {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.showSplash = false
		}
		return nil
	}

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
	spawnEnemies(g)
	handleEnemyBounces(g)
	deSpawnEnemies(g)
	collisionDetectionBulletsAndEnemies(g)
	handleShooting(g)
	collisionDetectionPlayerAndEnemies(g)
	handlePowerupCollection(g)
	return nil
}

func spawnEnemies(g *Game) {
	if g.frozenEnemiesTimer > 0 {
		return
	}

	// Randomly spawn an enemy every ~60 frames (1 second at 60fps)
	if rand.Float64() < 1.0/60.0 {
		screenWidth, screenHeight := g.Layout(0, 0)
		spawnX, spawnY, targetX, targetY := randomEdgeLocation(screenWidth, screenHeight)
		radius := 40.0

		// Calculate normalized velocity vector
		dx := float64(targetX - spawnX)
		dy := float64(targetY - spawnY)
		dist := math.Hypot(dx, dy)
		speed := 3.0 // pixels per frame
		vx := dx / dist * speed
		vy := dy / dist * speed

		isInvincible := rand.Float64() < 0.05 // 5% chance to be invincible

		enemy := &Enemy{
			X:            float64(spawnX),
			Y:            float64(spawnY),
			VX:           vx,
			VY:           vy,
			Radius:       radius,
			Active:       true,
			IsInvincible: isInvincible,
		}
		g.enemies = append(g.enemies, enemy)
	}
}

func deSpawnEnemies(g *Game) {
	screenWidth, screenHeight := g.Layout(0, 0)
	activeEnemies := g.enemies[:0]
	for _, e := range g.enemies {
		// move enemies in the direction they are travelling, assuming they arent frozen
		if g.frozenEnemiesTimer == 0 {
			e.X += e.VX
			e.Y += e.VY
		}

		if e.HitTimer > 0 {
			e.HitTimer--
			if e.HitTimer == 0 {
				e.Active = false // de-spawn after flash

				r := rand.Float64()

				if r > 0.0 && r < 0.10 { // 10% chance to drop a powerup
					r = rand.Float64()

					if r < 0.05 { // 5% chance to drop a shield
						powerup := &Powerup{
							X:      e.X,
							Y:      e.Y,
							Type:   powerupShield,
							Active: true,
						}
						g.powerups = append(g.powerups, powerup)
					} else if r < 0.1 { // 5% chance to drop a bomb
						powerup := &Powerup{
							X:      e.X,
							Y:      e.Y,
							Type:   powerupBomb,
							Active: true,
						}
						g.powerups = append(g.powerups, powerup)
					} else if r < 0.5 { // 5% chance to drop a freeze enemies
						powerup := &Powerup{
							X:      e.X,
							Y:      e.Y,
							Type:   powerupFreezeEnemies,
							Active: true,
						}
						g.powerups = append(g.powerups, powerup)
					} else {
						powerup := &Powerup{
							X:      e.X,
							Y:      e.Y,
							Type:   powerupInvincibleBullets,
							Active: true,
						}
						g.powerups = append(g.powerups, powerup)
					}
				}

				continue
			}
			activeEnemies = append(activeEnemies, e)
			continue
		}

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
	loadResources()
	resources.LoadBackground()

	game := &Game{}
	game.Reset()
	game.hasShield = true
	game.shieldTimer = 100000
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Space Shooter")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
