package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type Entity struct {
	Position    mgl.Vec4f
	Orientation mgl.Quatf
}

func (e Entity) Model() mgl.Mat4f {
	pos := e.Position
	return mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(e.Orientation.Mat4())
}

func (e Entity) View() mgl.Mat4f {
	return e.Model().Inv()
}
