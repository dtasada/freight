package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var sensitivity float32 = 2.5
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
		rl.ClearBackground(rl.SkyBlue)
		rl.BeginMode3D(camera)

		{ /* 3D drawing here */
			rl.DrawSphere(player.pos, player.radius, player.color)
			rl.DrawGrid(1000, 1)
			rl.DrawCube(rl.NewVector3(20, 20, 20), 10, 24, 10, rl.White)
			rl.DrawCube(rl.NewVector3(0, 0, 0), 1000, 0.1, 1000, rl.Green)
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
