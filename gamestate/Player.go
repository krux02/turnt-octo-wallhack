package gamestate

import (
	// "fmt"
	mgl "github.com/Jragonmiris/mathgl"
	// "math"
)

type PlayerInput struct {
	Move   mgl.Vec3f
	Rotate mgl.Vec3f
}

type Player struct {
	Camera   Camera
	Input    PlayerInput
	Velocity mgl.Vec3f
}

func (p *Player) GetCamera() *Camera {
	return &p.Camera
}

func (p *Player) Position() mgl.Vec3f {
	return p.Camera.Position
}
