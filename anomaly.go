package main

import "math/rand"

func (a *Anomaly) Update() error {

	if a.IsActive {
		// if we aren't at max alpha for anomaly, increase it
		// if it is at max, we start the flash timer
		if !a.flashing {
			if a.Alpha < 150 && a.fadeTimer%5 == 0 {
				a.Alpha += 3
				a.fadeTimer--
			} else if a.Alpha == 150 {
				a.fadeFlashTimer = 60
				a.flashing = true
			}
		}

		// if the fadeFlash is true, we flash the red and black
		if a.flashing {
			if a.fadeFlashTimer > 0 {

				if a.fadeFlashTimer%15 == 0 {
					a.Alpha = 10
				} else {
					a.Alpha = 150
				}

				a.fadeFlashTimer--
				if a.fadeFlashTimer == 0 {
					a.Deactivate()
				}
			}
		}
	}

	return nil
}

func (a *Anomaly) Activate() {
	a.IsActive = true
	a.fadeTimer = 300
	a.fadeFlashTimer = 0
	a.flashing = false
	a.Alpha = 1
	a.SafeRadius = 150

	a.SafeX = rand.Float64() * float64(1280)
	a.SafeY = rand.Float64() * float64(960)

	// clamp safex and safey to be within the screen bounds
	// so we get the whole safe circle on the screen
	if a.SafeX < a.SafeRadius {
		a.SafeX = a.SafeRadius
	}
	if a.SafeX > 1280-a.SafeRadius {
		a.SafeX = 1280 - a.SafeRadius
	}
	if a.SafeY < a.SafeRadius {
		a.SafeY = a.SafeRadius
	}
	if a.SafeY > 960-a.SafeRadius {
		a.SafeY = 960 - a.SafeRadius
	}
}

func (a *Anomaly) Deactivate() {
	a.IsActive = false
	a.fadeTimer = 0
	a.fadeFlashTimer = 0
	a.flashing = false
	a.Alpha = 0
	a.SafeX = 0
	a.SafeY = 0
	a.SafeRadius = 0
}
