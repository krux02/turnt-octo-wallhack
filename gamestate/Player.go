package gamestate

import (
	// "fmt"
	mgl "github.com/Jragonmiris/mathgl"
	// "math"
)

type PlayerInput struct {
	Move   mgl.Vec4f
	Rotate mgl.Vec4f
}

type Player struct {
	Camera   Camera
	Input    PlayerInput
	Velocity mgl.Vec4f
}

func (p *Player) GetCamera() *Camera {
	return &p.Camera
}

func (p *Player) Position() mgl.Vec4f {
	return p.Camera.Position
}
