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
