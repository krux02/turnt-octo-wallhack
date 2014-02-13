package world

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type Portal struct {
	Position    mgl.Vec3f
	Orientation mgl.Quatf
	Mesh        *Mesh
	Target      *Portal
}

func (this *Portal) ModelMat4() (Model mgl.Mat4f) {
	pos := this.Position
	Model = mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(this.Orientation.Mat4())
	return
}

func (this *Portal) ClippingPlane(front bool) mgl.Vec4f {
	var sgn float32
	if front {
		sgn = 1
	} else {
		sgn = -1
	}

	clippingPlane := this.ModelMat4().Mul4x1(mgl.Vec4f{0, sgn, 0, 0})
	p := this.Position
	clippingPlane[3] = -clippingPlane.Dot(mgl.Vec4f{p[0], p[1], p[2], 0})
	return clippingPlane
}
