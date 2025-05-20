package main

type Point struct {
	X int
	Y int
}

type Enemy struct {
	X, Y         float64
	VX, VY       float64
	Radius       float64
	Active       bool
	Size         int
	HitTimer     int
	IsInvincible bool
}

type Bullet struct {
	X, Y   float64
	VX, VY float64
	Active bool
}

// todo: refactor this a bit to use more entities to collectively store data
// e.g. maxSpeed, shipAngle, velocity, etc. could be in a separate player struct
type Game struct {
	playerLocation Point
	shipAngle      float64 // in radians
	velocity       float64
	maxSpeed       float64
	enemies        []*Enemy
	bullets        []*Bullet
	shootCooldown  int
	score          int
	showSplash     bool
	powerups       []*Powerup
	hasShield      bool
	shieldTimer    int
	bombs          int
	flashTimer     int
}

type Powerup struct {
	X, Y   float64
	Type   string // e.g. "shield"
	Active bool
}
