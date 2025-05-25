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
		return x, -1, rand.Intn(screenWidth), screenHeight
	case 1: // Bottom
		x := rand.Intn(screenWidth)
		return x, screenHeight, rand.Intn(screenWidth), -1
	case 2: // Left
		y := rand.Intn(screenHeight)
		return -1, y, screenWidth, rand.Intn(screenHeight)
	case 3: // Right
		y := rand.Intn(screenHeight)
		return screenWidth, y, -1, rand.Intn(screenHeight)
	}
	return 0, 0, 0, 0 // fallback, shouldn't happen
}

func collisionDetectionBulletsAndEnemies(g *Game) {
	for _, b := range g.bullets {
		if !b.Active {
			continue
		}

		for _, e := range g.enemies {
			// if the enemy isn't active, skip it.
			if !e.Active {
				continue
			}

			// calculate the distance between the bullet and the enemy
			dx := b.X - e.X
			dy := b.Y - e.Y
			distSq := dx*dx + dy*dy
			radius := e.Radius

			// it's a collision if the distance is less than the radius of the enemy
			if distSq < radius*radius {

				// if a bullet hits an enemy and the enemy is invincible, remove the bullet
				if e.IsInvincible {
					b.Active = false
					activeBullets := g.bullets[:0]
					for _, b := range g.bullets {
						if b.Active {
							activeBullets = append(activeBullets, b)
						}
					}
					g.bullets = activeBullets
					continue
				}

				// make the bullet inactive, so it can't hit more than one enemy
				if g.invincibleBulletsTimer == 0 {
					b.Active = false
					g.score += getScore(int(e.Radius))

					if g.score == 10 {
						g.Anomaly.Activate()
					}
				}

				// break enemy into smaller enemies
				// if the enemy is smaller than 10 radius, remove it
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

						if g.frozenEnemiesTimer > 0 {
							angle := rand.Float64() * 2 * math.Pi
							offset := rand.Float64() * 4 // up to 4 pixels
							newEnemy.X += math.Cos(angle) * offset
							newEnemy.Y += math.Sin(angle) * offset
						}

						g.enemies = append(g.enemies, newEnemy)
					}
				}
				e.HitTimer = 6 // flash before de-spawn
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
	return 0 // default case, shouldn't happen
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

func pointInPolygon(px, py float64, poly [][2]float64) bool {
	inside := false
	j := len(poly) - 1
	for i := 0; i < len(poly); i++ {
		xi, yi := poly[i][0], poly[i][1]
		xj, yj := poly[j][0], poly[j][1]
		if ((yi > py) != (yj > py)) &&
			(px < (xj-xi)*(py-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func distToSegmentSquared(px, py, x1, y1, x2, y2 float64) float64 {
	l2 := (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)
	if l2 == 0 {
		return (px-x1)*(px-x1) + (py-y1)*(py-y1)
	}
	t := ((px-x1)*(x2-x1) + (py-y1)*(y2-y1)) / l2
	if t < 0 {
		return (px-x1)*(px-x1) + (py-y1)*(py-y1)
	}
	if t > 1 {
		return (px-x2)*(px-x2) + (py-y2)*(py-y2)
	}
	projX := x1 + t*(x2-x1)
	projY := y1 + t*(y2-y1)
	return (px-projX)*(px-projX) + (py-projY)*(py-projY)
}

func polygonCircleCollision(poly [][2]float64, cx, cy, radius float64) bool {
	// 1. Check if circle center is inside polygon
	if pointInPolygon(cx, cy, poly) {
		return true
	}
	// 2. Check if any edge is close enough to the circle
	for i := 0; i < len(poly); i++ {
		j := (i + 1) % len(poly)
		if distToSegmentSquared(cx, cy, poly[i][0], poly[i][1], poly[j][0], poly[j][1]) < radius*radius {
			return true
		}
	}
	return false
}

func collisionDetectionPlayerAndEnemies(g *Game) {
	// Calculate ship polygon points (same as in DrawShip)
	if g.hasShield {
		return
	}
	cx := float64(g.playerLocation.X)
	cy := float64(g.playerLocation.Y)
	shipHeight := 75.0
	shipWidth := 30.0
	angle := g.shipAngle

	topX, topY := cx, cy-shipHeight/2
	rightX, rightY := cx+shipWidth/2, cy
	bottomX, bottomY := cx, cy+shipHeight/4
	leftX, leftY := cx-shipWidth/2, cy

	topX, topY = RotatePoint(topX, topY, cx, cy, angle)
	rightX, rightY = RotatePoint(rightX, rightY, cx, cy, angle)
	bottomX, bottomY = RotatePoint(bottomX, bottomY, cx, cy, angle)
	leftX, leftY = RotatePoint(leftX, leftY, cx, cy, angle)

	shipPoly := [][2]float64{
		{topX, topY},
		{rightX, rightY},
		{bottomX, bottomY},
		{leftX, leftY},
	}

	for _, e := range g.enemies {
		if !e.Active {
			continue
		}
		if polygonCircleCollision(shipPoly, e.X, e.Y, e.Radius) {
			g.previousScore = g.score
			g.Reset()
			return
		}
	}
}

func handleEnemyBounces(g *Game) {
	for i := 0; i < len(g.enemies); i++ {
		e1 := g.enemies[i]
		if !e1.Active {
			continue
		}
		for j := i + 1; j < len(g.enemies); j++ {
			e2 := g.enemies[j]
			if !e2.Active {
				continue
			}
			dx := e2.X - e1.X
			dy := e2.Y - e1.Y
			distSq := dx*dx + dy*dy
			rSum := e1.Radius + e2.Radius
			if distSq < rSum*rSum {
				dist := math.Sqrt(distSq)
				if dist == 0 {
					// Prevent division by zero
					dist = 0.1
				}

				// Move them apart so they're just touching
				overlap := 0.5 * (rSum - dist)
				nx := dx / dist
				ny := dy / dist
				e1.X -= nx * overlap
				e1.Y -= ny * overlap
				e2.X += nx * overlap
				e2.Y += ny * overlap

				// Swap velocities (simple elastic collision)
				e1.VX, e2.VX = e2.VX, e1.VX
				e1.VY, e2.VY = e2.VY, e1.VY
			}
		}
	}
}

func handlePowerupCollection(g *Game) {
	cx := float64(g.playerLocation.X)
	cy := float64(g.playerLocation.Y)
	playerRadius := 20.0 // or whatever fits your ship

	for _, p := range g.powerups {
		if !p.Active {
			continue
		}
		dx := cx - p.X
		dy := cy - p.Y
		if dx*dx+dy*dy < (playerRadius+12)*(playerRadius+12) {
			p.Active = false
			if p.Type == powerupShield {
				g.ActivateShield()
			} else if p.Type == powerupBomb {
				if g.bombs < 2 {
					g.bombs++
				}
			} else if p.Type == powerupInvincibleBullets {
				g.invincibleBulletsTimer = 300 // 5 seconds @ 60fps
			} else if p.Type == powerupFreezeEnemies {
				g.frozenEnemiesTimer = 300 // 5 seconds @ 60fps
			}

		}
	}
}

func (g *Game) ActivateShield() {
	g.hasShield = true
	g.shieldTimer = 300 // shield lasts for 300 frames (5 seconds at 60fps)
}
