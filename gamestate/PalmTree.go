package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

// instance data for each tree
type PalmTree struct {
	InstancePosition_ws mgl.Vec4f
}

type TreeVertex struct {
	Vertex_ms mgl.Vec4f
	TexCoord  mgl.Vec2f
}

// forest
type Forest struct {
	renderstuff.AbstractMesh
	Positions []PalmTree
	Model     mgl.Mat4f
}

func (this *Forest) GetModel() mgl.Mat4f {
	return this.Model
}

func (this *Forest) SetModel(model mgl.Mat4f) {
	this.Model = model
}

func (this *Forest) GetMesh() renderstuff.IMesh {
	return this
}

func (this *Forest) Vertices() interface{} {
	return []TreeVertex{
		TreeVertex{mgl.Vec4f{0, 1, 2, 1}, mgl.Vec2f{1, 0}},
		TreeVertex{mgl.Vec4f{0, 1, 0, 1}, mgl.Vec2f{1, 1}},
		TreeVertex{mgl.Vec4f{0, -1, 0, 1}, mgl.Vec2f{0, 1}},
		TreeVertex{mgl.Vec4f{0, -1, 2, 1}, mgl.Vec2f{0, 0}},
	}
}

func (this *Forest) Mode() renderstuff.Mode {
	return renderstuff.TriangleFan
}

func (this *Forest) InstanceData() interface{} {
	return this.Positions
}
