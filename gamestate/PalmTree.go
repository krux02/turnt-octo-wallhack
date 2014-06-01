package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
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
	AbstractMesh
	Positions []PalmTree
	Model     mgl.Mat4f
}

func (this *Forest) GetModel() mgl.Mat4f {
	return this.Model
}

func (this *Forest) SetModel(model mgl.Mat4f) {
	this.Model = model
}

func (this *Forest) GetMesh() IMesh {
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

func (this *Forest) Mode() Mode {
	return TriangleFan
}

func (this *Forest) InstanceData() interface{} {
	return this.Positions
}
