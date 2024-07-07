package src

import (
	_ "fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Pos           rl.Vector3
	Vel           rl.Vector3
	Color         rl.Color
	Radius        float32
	Direction     float32
	MovementSpeed float32
	JumpVel       float32
	Angle         float64
	RollFrame     float64
	RollStep      float64
	RollAmp       float64
	RollPeriod    float64
	OnGround      bool
}

func NewPlayer(radius float32, pos rl.Vector3, color rl.Color) Player {
	return Player{
		Pos:           pos,
		Vel:           rl.Vector3Zero(),
		Color:         color,
		Radius:        radius,
		Direction:     0,
		MovementSpeed: 0.25,
		JumpVel:       350,
		Angle:         0,
		RollFrame:     0,
		RollStep:      0.2,
		RollAmp:       1.0 / 8.0,
		RollPeriod:    2 * math.Pi,
		OnGround:      false,
	}
}

func move(direction Direction, distance float32) {
	var movement Vector3
	if direction == DirectionRight {
		movement = rl.GetCameraRight(&Camera)
	} else if direction == DirectionForward {
		movement = rl.GetCameraForward(&Camera)
	}
	movement.Y = 0

	movement = rl.Vector3Scale(movement, distance)

	Camera.Position = rl.Vector3Add(Camera.Position, movement)
	Camera.Target = rl.Vector3Add(Camera.Target, movement)
}

func (self *Player) Update() {
	dt := rl.GetFrameTime()

	/* Gravity */
	self.Vel.Y -= Gravity * dt
	self.Pos = rl.Vector3Add(self.Pos, Vector3MultiplyValue(self.Vel, dt))

	/* Movement keys */
	if rl.IsCursorHidden() {
		mouseDelta := rl.GetMouseDelta()
		rl.CameraYaw(&Camera, -mouseDelta.X*Sensitivity, 0)
		rl.CameraPitch(&Camera, -mouseDelta.Y*Sensitivity, 1, 0, 0)
	}

	self.Angle = math.Atan2(
		float64(Camera.Target.Z)-float64(Camera.Position.Z),
		float64(Camera.Target.X)-float64(Camera.Position.X),
	)

	var movementSpeed float32
	if rl.IsKeyDown(rl.KeyLeftShift) {
		movementSpeed = self.MovementSpeed * 2
	} else if rl.IsKeyDown(rl.KeyC) {
		movementSpeed = self.MovementSpeed / 2
	} else {
		movementSpeed = self.MovementSpeed
	}

	if rl.IsKeyDown(rl.KeyW) {
		self.RollFrame += self.RollStep

		move(DirectionForward, movementSpeed)
	} else {
		if self.RollFrame > 0 {
			self.RollFrame -= self.RollStep
		} else {
			self.RollFrame += self.RollStep
		}

		if self.RollFrame > -self.RollStep || self.RollFrame < self.RollStep {
			self.RollFrame = 0
		}
	}

	sideRoll := float64(4)
	if rl.IsKeyDown(rl.KeyA) {
		if self.RollFrame < sideRoll {
			self.RollFrame += self.RollStep
		} else {
			self.RollFrame = sideRoll
		}

		move(DirectionRight, -movementSpeed)
	}
	if rl.IsKeyDown(rl.KeyS) {
		move(DirectionForward, -movementSpeed)
	}
	if rl.IsKeyDown(rl.KeyD) {
		if self.RollFrame > -sideRoll {
			self.RollFrame -= self.RollStep
		} else {
			self.RollFrame = -sideRoll
		}

		move(DirectionRight, movementSpeed)
	}

	if self.OnGround {
		if rl.IsKeyDown(rl.KeySpace) {
			self.Vel.Y += self.JumpVel * dt
			self.OnGround = false
		}
	}

	if self.RollFrame > self.RollPeriod {
		self.RollFrame = self.RollFrame - math.Floor(self.RollFrame)
	}

	self.Pos = Camera.Position
}
