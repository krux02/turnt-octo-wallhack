package main

import (
	// "fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	// "math"
)

type PlayerInput struct {
	move   mgl.Vec3f
	rotate mgl.Vec3f
}

type Player interface {
	SetInput(input PlayerInput)
	GetCamera() *rendering.Camera
	Position() mgl.Vec3f
	Update(gamestate *GameState)
}

type MyPlayer struct {
	Camera   rendering.Camera
	input    PlayerInput
	velocety mgl.Vec3f
}

func (p *MyPlayer) SetInput(input PlayerInput) {
	p.input = input
}

func (p *MyPlayer) GetCamera() *rendering.Camera {
	return &p.Camera
}

func (p *MyPlayer) Position() mgl.Vec3f {
	return p.Camera.Position
}

func (p *MyPlayer) Update(gamestate *GameState) {
	rot := p.input.rotate
	p.Camera.Yaw(rot[0])
	p.Camera.Pitch(rot[1])
	p.Camera.Roll(rot[2])

	move := p.input.move
	if move.Len() > 0 {
		move = move.Normalize()
	}
	move = p.Camera.Direction.Rotate(move)

	if gamestate.Options.DisablePlayerPhysics {
		move = move.Mul(0.1)
		p.velocety = move
		p.Camera.MoveAbsolute(move)
	} else {
		move = move.Mul(0.01)
		p.velocety = p.velocety.Add(move)
		p.Camera.MoveAbsolute(p.velocety)

		groundHeight := gamestate.World.HeightMap.Get2f(p.Position()[0], p.Position()[1])

		height := p.Camera.Position[2]
		minHeight := groundHeight + 1.5
		maxHeight := groundHeight + 20

		if height < minHeight {
			diff := minHeight - height
			p.velocety[2] += diff
			p.Camera.Position[2] += diff
		}

		/*
			w := float32(gamestate.World.HeightMap.W)
			h := float32(gamestate.World.HeightMap.H)
			if p.Camera.Position[0] < 0 {
				p.Camera.Position[0] += w
			} else if p.Camera.Position[0] >= w {
				p.Camera.Position[0] -= w
			}

			if p.Camera.Position[1] < 0 {
				p.Camera.Position[1] += h
			} else if p.Camera.Position[1] >= h {
				p.Camera.Position[1] -= h
			}
		*/

		p.velocety = p.velocety.Mul(0.95)

		if height > maxHeight && p.velocety[2] > -0.0 {
			p.velocety[2] -= 0.001
		}

		groundFactor := (height - groundHeight) / 20
		if groundFactor > 1 {
			groundFactor = 1
		}
		if groundFactor < 0 {
			groundFactor = 0
		}
		groundFactor = 1 - groundFactor
		groundFactor = groundFactor * groundFactor
	}
}
