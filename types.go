package main

type Point struct {
	X int
	Y int
}

type Enemy struct {
	X, Y   float64
	VX, VY float64
	Radius float64
	Active bool
	Size   int
}

type Bullet struct {
	X, Y   float64
	VX, VY float64
	Active bool
}

type Game struct {
	playerLocation  Point
	playerDirection int     // 0 == up, 180 = down, 90 = right, 270 = left
	shipAngle       float64 // in radians
	velocity        float64
	maxSpeed        float64
	enemies         []*Enemy
	bullets         []*Bullet
	shootCooldown   int
}
