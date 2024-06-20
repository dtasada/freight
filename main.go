package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	pos           rl.Vector3
	vel           rl.Vector3
	color         rl.Color
	radius        float32
	direction     float32
	rollFrame     float64
	movementSpeed float32
	onGround      bool
}

func newPlayer(radius float32, pos rl.Vector3, color rl.Color) Player {
	return Player{
		pos:           pos,
		vel:           rl.Vector3Zero(),
		color:         color,
		radius:        radius,
		direction:     0,
		rollFrame:     0,
		movementSpeed: 10,
		onGround:      false,
	}
}

func (self *Player) update() {
	dt := rl.GetFrameTime()

	// Gravity
	self.vel.Y -= gravity * dt
	self.pos.Y += self.vel.Y * dt
	camera.Target.Y += self.vel.Y * dt
	if self.pos.Y < 2 { // 2 is the "floor" and not 0 because eye level
		self.pos.Y = 2
		camera.Target.Y = 2
		self.vel.Y *= 0 // -0.8
		self.onGround = true
	}

	// Movement keys
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

		self.pos.X -= dz
		self.pos.Z += dx
		camera.Target.X -= dz
		camera.Target.Z += dx
	}

	if rl.IsKeyDown(rl.KeySpace) && self.onGround {
		self.vel.Y += 2 * dt
		self.onGround = false
	}
	if rl.IsKeyDown(rl.KeyC) {
		self.vel.Y -= 2 * dt
	}

}

const sensitivity float32 = 2.5

var gravity float32 = 9.81
var player Player = newPlayer(
	1,
	rl.NewVector3(0.0, 20.0, 0.0),
	rl.Yellow,
)
var camera rl.Camera3D = rl.Camera3D{
	Position:   player.pos,
	Target:     rl.NewVector3(0.0, 20.0, 10.0), // the 10.0 should be arbitrary?
	Up:         rl.NewVector3(0.0, 1.0, 0.0),   // Y is "up"
	Fovy:       100.0,
	Projection: rl.CameraPerspective,
}

func main() {
	// Init
	fmt.Println("Running game...")
	rl.InitWindow(1280, 720, "freight")
	rl.SetTargetFPS(60)
	rl.HideCursor()

	for !rl.WindowShouldClose() {
		/* Non-rendering logic here */
		player.update()

		/* Camera logic */
		camera.Position = player.pos
		camera.Position.Y = player.pos.Y + 2
		camera.Up.Z = 1 / 10 * float32(math.Sin(player.rollFrame)) // Camera roll wave
		mouseNormal := rl.Vector2Normalize(rl.GetMouseDelta())
		rl.UpdateCameraPro(
			&camera,
			rl.Vector3Zero(), // Movement vector
			rl.NewVector3( // Rotation vector
				sensitivity*mouseNormal.X,
				sensitivity*mouseNormal.Y,
				0, // Roll handled elsewhere
			),
			0.0, // No zoom
		)

		/* Rendering logic */
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode3D(camera)

		{ /* 3D drawing here */
			rl.DrawSphere(player.pos, player.radius, player.color)
			rl.DrawGrid(1000, 1)
			rl.DrawCube(rl.NewVector3(20, 20, 20), 10, 24, 10, rl.White)
		}

		rl.EndMode3D()

		/* 2D drawing here */
		rl.DrawFPS(16, 16)
		rl.DrawText(fmt.Sprintf("camera position: %v", camera.Position), 16, 20*2+8, 20, rl.White)
		rl.DrawText(fmt.Sprintf("player pos: %v", player.pos), 16, 20*3+8, 20, rl.White)
		rl.DrawText(fmt.Sprintf("camera target: %v", camera.Target), 16, 20*4+8, 20, rl.White)
		rl.DrawText(fmt.Sprintf("player vel: %v", player.vel), 16, 20*5+8, 20, rl.White)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
