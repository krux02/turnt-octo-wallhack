package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type Entity struct {
	Position    mgl.Vec3f
	Orientation mgl.Quatf
}

func (e Entity) Model() mgl.Mat4f {
	pos := e.Position
	return mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(e.Orientation.Mat4())
}

func (e Entity) View() mgl.Mat4f {
	Tx := e.Position[0]
	Ty := e.Position[1]
	Tz := e.Position[2]
	return e.Orientation.Inverse().Mat4().Mul4(mgl.Translate3D(-Tx, -Ty, -Tz))
}
