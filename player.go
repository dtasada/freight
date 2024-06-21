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
		camera.Target.Y = 2
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
		self.rollFrame += 0.2

		self.pos.Z += dz
		self.pos.X += dx
		camera.Target.Z += dz
		camera.Target.X += dx
	} else {
		self.rollFrame = 0
	}

	if rl.IsKeyDown(rl.KeyA) {
		self.rollFrame = 1

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
		self.rollFrame = -1

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

}
