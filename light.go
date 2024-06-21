package main

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightType int32

const (
	LightTypeDirectional LightType = iota
	LightTypePoint
)

type Light struct {
	Shader    rl.Shader
	LightType LightType
	Position  rl.Vector3
	Target    rl.Vector3
	Color     rl.Color
	Enabled   int32
	// shader locations
	EnabledLoc int32
	TypeLoc    int32
	PosLoc     int32
	TargetLoc  int32
	ColorLoc   int32
}

const maxLightsCount = 4

var lightCount = 0

func NewLight(
	lightType LightType,
	position, target rl.Vector3,
	color rl.Color,
	shader rl.Shader,
) Light {
	light := Light{
		Shader: shader,
	}
	if lightCount < maxLightsCount {
		light.Enabled = 1
		light.LightType = lightType
		light.Position = position
		light.Target = target
		light.Color = color
		light.EnabledLoc = rl.GetShaderLocation(shader, fstr("lights[%d].enabled", lightCount))
		light.TypeLoc = rl.GetShaderLocation(shader, fstr("lights[%d].type", lightCount))
		light.PosLoc = rl.GetShaderLocation(shader, fstr("lights[%d].position", lightCount))
		light.TargetLoc = rl.GetShaderLocation(shader, fstr("lights[%d].target", lightCount))
		light.ColorLoc = rl.GetShaderLocation(shader, fstr("lights[%d].color", lightCount))
		light.UpdateValues()
		lightCount++
	}
	return light
}

func (lt *Light) UpdateValues() {
	// Send to shader light enabled state and type
	rl.SetShaderValue(lt.Shader, lt.EnabledLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.Enabled)), 4), rl.ShaderUniformInt)
	rl.SetShaderValue(lt.Shader, lt.TypeLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.LightType)), 4), rl.ShaderUniformInt)

	// Send to shader light position values
	rl.SetShaderValue(lt.Shader, lt.PosLoc, []float32{lt.Position.X, lt.Position.Y, lt.Position.Z}, rl.ShaderUniformVec3)

	// Send to shader light target target values
	rl.SetShaderValue(lt.Shader, lt.TargetLoc, []float32{lt.Target.X, lt.Target.Y, lt.Target.Z}, rl.ShaderUniformVec3)

	// Send to shader light color values
	rl.SetShaderValue(
		lt.Shader, lt.ColorLoc,
		[]float32{
			float32(lt.Color.R) / 255,
			float32(lt.Color.G) / 255,
			float32(lt.Color.B) / 255,
			float32(lt.Color.A) / 255,
		},
		rl.ShaderUniformVec4,
	)
}
