package main

import (
	_ "fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	pos           rl.Vector3
	vel           rl.Vector3
	color         rl.Color
	radius        float32
	direction     float32
	movementSpeed float32
	jumpVel       float32
	rollFrame     float64
	rollStep      float64
	rollAmp       float64
	rollPeriod    float64
	onGround      bool
}

func newPlayer(radius float32, pos rl.Vector3, color rl.Color) Player {
	return Player{
		pos:           pos,
		vel:           rl.Vector3Zero(),
		color:         color,
		radius:        radius,
		direction:     0,
		movementSpeed: 10,
		jumpVel:       350,
		rollFrame:     0,
		rollStep:      0.2,
		rollAmp:       1.0 / 8.0,
		rollPeriod:    2 * math.Pi,
		onGround:      false,
	}
}

func (self *Player) update() {
	dt := rl.GetFrameTime()

	/* Gravity */
	self.vel.Y -= gravity * dt
	self.pos.Y += self.vel.Y * dt
	camera.Target.Y += self.vel.Y * dt
	if self.pos.Y < 2 { // 2 is the "floor" and not 0 because eye level
		self.pos.Y = 2
		self.vel.Y *= 0 // -0.8
		self.onGround = true
	}

	/* Movement keys */
	// meth
	cameraDirection := rl.Vector3Subtract(camera.Target, camera.Position)
	angle := float64(
		rl.Vector2Angle(
			rl.NewVector2(0, 1),
			rl.NewVector2(cameraDirection.X, cameraDirection.Z),
		),
	) + math.Pi/2
	dz := float32(math.Sin(angle)) * self.movementSpeed * dt
	dx := float32(math.Cos(angle)) * self.movementSpeed * dt
	if rl.IsKeyDown(rl.KeyW) {
		self.rollFrame += self.rollStep

		self.pos.Z += dz
		self.pos.X += dx
		camera.Target.Z += dz
		camera.Target.X += dx
	} else {
		if self.rollFrame > 0 {
			self.rollFrame -= self.rollStep
		} else {
			self.rollFrame += self.rollStep
		}

		if self.rollFrame > -self.rollStep || self.rollFrame < self.rollStep {
			self.rollFrame = 0
		}
	}

	sideRoll := float64(4)
	if rl.IsKeyDown(rl.KeyA) {
		if self.rollFrame < sideRoll {
			self.rollFrame += self.rollStep
		} else {
			self.rollFrame = sideRoll
		}

		// meth
		self.pos.X += dz
		self.pos.Z -= dx
		camera.Target.X += dz
		camera.Target.Z -= dx
	}
	if rl.IsKeyDown(rl.KeyS) {
		self.pos.Z -= dz
		self.pos.X -= dx
		camera.Target.Z -= dz
		camera.Target.X -= dx
	}
	if rl.IsKeyDown(rl.KeyD) {
		if self.rollFrame > -sideRoll {
			self.rollFrame -= self.rollStep
		} else {
			self.rollFrame = -sideRoll
		}

		// meth
		self.pos.X -= dz
		self.pos.Z += dx
		camera.Target.X -= dz
		camera.Target.Z += dx
	}

	if self.onGround {
		if rl.IsKeyDown(rl.KeySpace) {
			self.vel.Y += self.jumpVel * dt
			self.onGround = false
		}
		if rl.IsKeyDown(rl.KeyC) {
			self.vel.Y -= self.jumpVel * dt
		}
	}

	if self.rollFrame > self.rollPeriod {
		self.rollFrame = self.rollFrame - math.Floor(self.rollFrame)
	}

}
