package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	opt "github.com/moznion/go-optional"

	"github.com/dtasada/freight/src"
)

var sensitivity float32 = 3.5
var gravity float32 = 9.81
var fstr = fmt.Sprintf

var player src.Player = src.NewPlayer(
	1,
	rl.NewVector3(0.0, 20.0, 0.0),
	rl.Yellow,
)

func main() {
	/* Init raylib */
	fmt.Println("Running game...")
	rl.InitWindow(1280, 720, "paper")
	rl.SetExitKey(0)
	rl.DisableCursor()
	rl.SetTargetFPS(src.TargetFPS)

	/* Raylib flags */
	rl.SetTraceLogLevel(rl.LogWarning)
	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable 4x MSAA if available
	rl.SetConfigFlags(rl.FlagWindowResizable)

	/* Generate textures */
	lightShader := rl.LoadShader("./resources/shaders/lighting.vs", "./resources/shaders/lighting.fs")
	defer rl.UnloadShader(lightShader)

	*lightShader.Locs = rl.GetShaderLocation(lightShader, "viewPos")
	ambientLoc := rl.GetShaderLocation(lightShader, "ambient")
	shaderValue := []float32{0.1, 0.1, 0.1, 1.0}
	rl.SetShaderValue(lightShader, ambientLoc, shaderValue, rl.ShaderUniformVec4)

	ground := src.NewObjectFromMesh(rl.GenMeshPlane(1000, 1000, 3, 3), nil, opt.Some(lightShader))
	defer rl.UnloadModel(ground.Model)

	cube := src.NewObjectFromMesh(rl.GenMeshCube(10, 24, 10), nil, opt.Some(lightShader))
	defer rl.UnloadModel(cube.Model)

	gun := src.NewObjectFromFile("./resources/models/revolver.gltf", opt.Some("./resources/images/revolver.png"), opt.Some(lightShader))
	defer rl.UnloadModel(gun.Model)

	light := src.NewLight(src.LightTypePoint, player.Pos, rl.Vector3Zero(), rl.Yellow, 1, lightShader)

	for !rl.WindowShouldClose() {
		{ /* Pre-render logic here */
			if rl.IsCursorHidden() {
				light.Position = src.Camera.Target
				light.Target = src.Camera.Target
			}

			rl.SetShaderValue(
				lightShader,
				*lightShader.Locs,
				[]float32{src.Camera.Position.X, src.Camera.Position.Y, src.Camera.Position.Z},
				rl.ShaderUniformVec3,
			)

			if rl.IsKeyPressed(rl.KeyEscape) {
				if rl.IsCursorHidden() {
					rl.EnableCursor()
				} else {
					rl.DisableCursor()
				}
			}

			player.Update()
		}

		{ /* Drawing */
			rl.BeginDrawing()

			rl.ClearBackground(rl.SkyBlue)

			{ /* 3D drawing here */
				rl.BeginMode3D(src.Camera)

				rl.DrawModel(ground.Model, rl.Vector3Zero(), 1, rl.Green)
				rl.DrawModel(cube.Model, rl.NewVector3(6, 0, 6), 1, rl.DarkBlue)

				rl.DrawModel(gun.Model, src.RelativePlacement(player.Pos, rl.NewVector3(0.3, 0, 0.3), player.Angle), 0.05, rl.White)

				light.Update()

				rl.EndMode3D()
			}

			{ /* 2D drawing here */
				rl.DrawFPS(16, 16)
				rl.DrawText(fstr("camera position: %v", src.Camera.Position), 16, 20*2+8, 20, rl.White)
				rl.DrawText(fstr("camera target: %v", src.Camera.Target), 16, 20*3+8, 20, rl.White)
				rl.DrawText(fstr("camera up: %v", src.Camera.Up), 16, 20*4+8, 20, rl.White)
				rl.DrawText(fstr("player pos: %v", player.Pos), 16, 20*5+8, 20, rl.White)
				rl.DrawText(fstr("player vel: %v", player.Vel), 16, 20*6+8, 20, rl.White)
				rl.DrawText(fstr("player angle: %v", player.Angle), 16, 20*7+8, 20, rl.White)
				rl.DrawText(fstr("player rollFrame: %v", player.RollFrame), 16, 20*8+8, 20, rl.White)
			}

			rl.EndDrawing()
		}
	}

	rl.CloseWindow()
}
