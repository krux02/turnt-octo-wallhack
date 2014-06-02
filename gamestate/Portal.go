package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type Portal struct {
	Entity
	Normal mgl.Vec4f
	Target *Portal
	mesh   *renderstuff.Mesh
}

func (this *Portal) Mesh() *renderstuff.Mesh {
	return this.mesh
}

func (this *Portal) SetMesh(mesh *renderstuff.Mesh) {
	this.mesh = mesh
}

func (this *Portal) ClippingPlane(front bool) mgl.Vec4f {
	var sgn float32
	if front {
		sgn = 1
	} else {
		sgn = -1
	}

	clippingPlane := this.Model().Mul4x1(this.Normal.Mul(sgn))
	p := this.Position
	clippingPlane[3] = -clippingPlane.Dot(mgl.Vec4f{p[0], p[1], p[2], 0})
	return clippingPlane
}

func (this *Portal) Transform() mgl.Mat4f {
	return this.Target.Model().Mul4(this.View())
}

func (this *Portal) Dimension(dim int) float32 {
	return this.Position[dim]
}
