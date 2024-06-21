package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var sensitivity float32 = 3.5
var gravity float32 = 9.81
var fstr = fmt.Sprintf

var player Player = newPlayer(
	1,
	rl.NewVector3(0.0, 20.0, 0.0),
	rl.Yellow,
)
var camera rl.Camera3D = rl.Camera3D{
	Position:   player.pos,
	Target:     rl.NewVector3(0.0, 20.0, 10.0), // the Z:10.0 should be arbitrary?
	Up:         rl.NewVector3(0.0, 1.0, 0.0),   // Asserts Y to be the vertical axis
	Fovy:       100.0,
	Projection: rl.CameraPerspective,
}

func main() {
	/* Init */
	fmt.Println("Running game...")
	rl.InitWindow(1280, 720, "freight")
	rl.SetTargetFPS(60)
	rl.DisableCursor()
	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable 4x MSAA if available

	/* Generate textures */
	ground := rl.LoadModelFromMesh(rl.GenMeshPlane(1000, 1000, 3, 3))
	defer rl.UnloadModel(ground)

	cube := rl.LoadModelFromMesh(rl.GenMeshCube(10, 24, 10))
	defer rl.UnloadModel(cube)

	lightShader := rl.LoadShader("./resources/shaders/lighting.vs", "./resources/shaders/lighting.fs")
	defer rl.UnloadShader(lightShader)

	grass := rl.LoadTexture("./resources/images/grass.png")
	defer rl.UnloadTexture(grass)

	grassShader := rl.LoadShader("./resources/shaders/base.vs", "./resources/shaders/grayscale.fs")
	defer rl.UnloadShader(grassShader)

	*lightShader.Locs = rl.GetShaderLocation(lightShader, "viewPos")
	ambientLoc := rl.GetShaderLocation(lightShader, "ambient")
	shaderValue := []float32{0.1, 0.1, 0.1, 1.0}
	rl.SetShaderValue(lightShader, ambientLoc, shaderValue, rl.ShaderUniformVec4)

	ground.Materials.Shader = lightShader
	cube.Materials.Shader = lightShader

	lights := []Light{
		NewLight(LightTypePoint, rl.NewVector3(0, 100, 0), rl.NewVector3(00, 00, 00), rl.Yellow, lightShader),
	}

	for !rl.WindowShouldClose() {
		/* Non-rendering logic here */
		player.update()

		/* Camera logic */
		camera.Position = player.pos
		camera.Position.Y = player.pos.Y + 2
		camera.Up.Z = float32(math.Sin(player.rollFrame) * player.rollAmp) // Camera roll wave
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

		rl.SetShaderValue(lightShader,
			*lightShader.Locs,
			[]float32{camera.Position.X, camera.Position.Y, camera.Position.Z},
			rl.ShaderUniformVec3,
		)

		/* Rendering logic */
		rl.BeginDrawing()

		rl.ClearBackground(rl.SkyBlue)

		rl.BeginMode3D(camera)

		/* 3D drawing here */
		rl.BeginShaderMode(grassShader)
		// rl.DrawModel(ground, rl.Vector3Zero(), 1, rl.Green)
		rl.DrawTexture(grass, 0, 0, rl.White)
		rl.EndShaderMode()

		rl.DrawSphere(player.pos, player.radius, player.color)
		rl.DrawModel(cube, rl.NewVector3(6, 0, 6), 1, rl.DarkBlue)

		for _, light := range lights {
			light.UpdateValues()

			if light.Enabled == 1 {
				rl.DrawSphereEx(light.Position, 0.2, 8, 8, light.Color)
			} else {
				rl.DrawSphereWires(light.Position, 0.2, 8, 8, rl.Fade(light.Color, 0.3))
			}
		}

		rl.EndMode3D()

		/* 2D drawing here */
		rl.DrawFPS(16, 16)
		rl.DrawText(fstr("camera position: %v", camera.Position), 16, 20*2+8, 20, rl.White)
		rl.DrawText(fstr("player pos: %v", player.pos), 16, 20*3+8, 20, rl.White)
		rl.DrawText(fstr("camera target: %v", camera.Target), 16, 20*4+8, 20, rl.White)
		rl.DrawText(fstr("player vel: %v", player.vel), 16, 20*5+8, 20, rl.White)
		rl.DrawText(fstr("camera up: %v", camera.Up), 16, 20*6+8, 20, rl.White)
		rl.DrawText(fstr("player rollFrame: %v", player.rollFrame), 16, 20*7+8, 20, rl.White)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
