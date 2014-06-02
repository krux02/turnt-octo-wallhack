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
	mesh      *renderstuff.Mesh
	Positions []PalmTree
	model     mgl.Mat4f
}

func (this *Forest) Model() mgl.Mat4f {
	return this.model
}

func (this *Forest) SetModel(model mgl.Mat4f) {
	this.model = model
}

func (this *Forest) Mesh() *renderstuff.Mesh {
	if this.mesh == nil {
		this.mesh = &renderstuff.Mesh{
			Vertices: []TreeVertex{
				TreeVertex{mgl.Vec4f{0, 1, 2, 1}, mgl.Vec2f{1, 0}},
				TreeVertex{mgl.Vec4f{0, 1, 0, 1}, mgl.Vec2f{1, 1}},
				TreeVertex{mgl.Vec4f{0, -1, 0, 1}, mgl.Vec2f{0, 1}},
				TreeVertex{mgl.Vec4f{0, -1, 2, 1}, mgl.Vec2f{0, 0}},
			},
			InstanceData: this.Positions,
			Mode:         renderstuff.TriangleFan,
		}
	}
	return this.mesh
}
