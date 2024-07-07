/* package src: base engine functions structs and constants */
package src

import (
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)

type Vector3 = rl.Vector3
type f32 = float32
type f64 = float64

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionForward
)

var Gravity float32 = 1 / 250
var TargetFPS int32 = 60
var Sensitivity float32 = 0.0035
var MovementSpeed float32 = 1.0

var Caskaydia rl.Font

var Camera = rl.Camera3D{
	Position:   rl.NewVector3(0, 4, 0),
	Target:     rl.NewVector3(0, 4, 1),
	Up:         rl.NewVector3(0.0, 1.0, 0.0), // Asserts Y to be the vertical axis
	Fovy:       100.0,
	Projection: rl.CameraPerspective,
}

func Vector3DivideValue(vec Vector3, div float32) Vector3 {
	return rl.NewVector3(vec.X/div, vec.Y/div, vec.Z/div)
}

func Vector3MultiplyValue(vec Vector3, mult float32) Vector3 {
	return rl.NewVector3(vec.X*mult, vec.Y*mult, vec.Z*mult)
}

func RelativePlacement(origin, offset Vector3, angle float64) Vector3 {
	otAngle := float64(rl.Vector3Angle(origin, offset))
	dz := float32(math.Sin(otAngle))
	dx := float32(math.Cos(otAngle))

	return rl.NewVector3(
		origin.X+dx,
		origin.Y,
		origin.Z+dz,
	)
}

func DrawText(text string, x, y float32, fontSize float32, tint rl.Color) {
	rl.DrawTextEx(Caskaydia, text, rl.NewVector2(x, y), fontSize, 1, tint)
}
