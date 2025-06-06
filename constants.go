package main

const (
	rotateSpeed = 0.06 // radians per frame
	accel       = 0.2  // acceleration per frame
	friction    = 0.01 // natural slow down

	bulletSpeed    = 10.0
	bulletRadius   = 4.0
	bulletCooldown = 15 // frames between shots
	enemySpeed     = 1.0
)

const (
	tileSize = 32
)

const (
	powerupShield            = 1
	powerupBomb              = 2
	powerupInvincibleBullets = 3
	powerupFreezeEnemies     = 4
	powerUpMystery           = 5
)
