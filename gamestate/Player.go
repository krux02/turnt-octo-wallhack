package gamestate

import (
	// "fmt"
	mgl "github.com/krux02/mathgl/mgl32"
	// "math"
)

type PlayerInput struct {
	Move   mgl.Vec4
	Rotate mgl.Vec4
}

type Player struct {
	Camera   Camera
	Input    PlayerInput
	Velocity mgl.Vec4
}

func (p *Player) GetCamera() *Camera {
	return &p.Camera
}

func (p *Player) Position() mgl.Vec4 {
	return p.Camera.Position
}
