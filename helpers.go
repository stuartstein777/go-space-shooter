package main

import (
	"math"
	"math/rand"
)

func RotatePoint(x, y, cx, cy, angle float64) (float64, float64) {
	sin, cos := math.Sin(float64(angle)), math.Cos(float64(angle))
	dx, dy := float64(x-cx), float64(y-cy)
	return cx + dx*cos - dy*sin, cy + dx*sin + dy*cos
}

func randomEdgeLocation(screenWidth, screenHeight int) (int, int, int, int) {
	edge := rand.Intn(4) // 0=top, 1=bottom, 2=left, 3=right
	switch edge {
	case 0: // Top
		x := rand.Intn(screenWidth)
		return x, -1, x, screenHeight - 1
	case 1: // Bottom
		x := rand.Intn(screenWidth)
		return x, screenHeight, x, 0
	case 2: // Left
		y := rand.Intn(screenHeight)
		return -1, y, screenWidth - 1, y
	case 3: // Right
		y := rand.Intn(screenHeight)
		return screenWidth, y, 0, y
	}
	return 0, 0, 0, 0 // fallback, shouldn't happen
}

func collisionDetectionBulletsAndEnemies(g *Game) {
	for _, b := range g.bullets {
		if !b.Active {
			continue
		}
		for _, e := range g.enemies {
			if !e.Active {
				continue
			}
			dx := b.X - e.X
			dy := b.Y - e.Y
			distSq := dx*dx + dy*dy
			radius := e.Radius
			if distSq < radius*radius {
				b.Active = false
				g.score += getScore(int(e.Radius))
				if e.Radius > 10 {
					newRadius := e.Radius / 2
					// if the enemy is larger than 10 radius, split it into two smaller enemies
					// spawn them in random directions
					for i := 0; i < 2; i++ {
						angle := rand.Float64() * 2 * math.Pi
						speed := 3.0
						vx := math.Cos(angle) * speed
						vy := math.Sin(angle) * speed
						newEnemy := &Enemy{
							X:      e.X,
							Y:      e.Y,
							VX:     vx,
							VY:     vy,
							Radius: newRadius,
							Active: true,
						}
						g.enemies = append(g.enemies, newEnemy)
					}
				}
				e.HitTimer = 6 // flash before despawn
				break
			}
		}
	}
}

func getScore(radius int) int {
	switch {
	case radius == 40:
		return 10
	case radius == 20:
		return 20
	case radius == 10:
		return 40
	}
}

func handleShooting(g *Game) {
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
}
