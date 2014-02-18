package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type Portal struct {
	Entity
	Mesh   *Mesh
	Target *Portal
}

func (this *Portal) ClippingPlane(front bool) mgl.Vec4f {
	var sgn float32
	if front {
		sgn = 1
	} else {
		sgn = -1
	}

	clippingPlane := this.Model().Mul4x1(mgl.Vec4f{0, sgn, 0, 0})
	p := this.Position
	clippingPlane[3] = -clippingPlane.Dot(mgl.Vec4f{p[0], p[1], p[2], 0})
	return clippingPlane
}

func (this *Portal) Transform() mgl.Mat4f {
	Mat1 := this.View()
	Mat2 := this.Target.Model()
	return Mat2.Mul4(Mat1)
}
