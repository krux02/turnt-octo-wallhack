package gamestate

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

// instance data for each tree
type PalmTree struct {
	InstancePosition_ws mgl.Vec4
}

type TreeVertex struct {
	Vertex_ms mgl.Vec4
	TexCoord  mgl.Vec2
}

// forest
type Forest struct {
	mesh      *renderstuff.Mesh
	Positions []PalmTree
	model     mgl.Mat4
}

func (this *Forest) Model() mgl.Mat4 {
	return this.model
}

func (this *Forest) SetModel(model mgl.Mat4) {
	this.model = model
}

func (this *Forest) Mesh() *renderstuff.Mesh {
	if this.mesh == nil {
		this.mesh = &renderstuff.Mesh{
			Vertices: []TreeVertex{
				TreeVertex{mgl.Vec4{0, 1, 2, 1}, mgl.Vec2{1, 0}},
				TreeVertex{mgl.Vec4{0, 1, 0, 1}, mgl.Vec2{1, 1}},
				TreeVertex{mgl.Vec4{0, -1, 0, 1}, mgl.Vec2{0, 1}},
				TreeVertex{mgl.Vec4{0, -1, 2, 1}, mgl.Vec2{0, 0}},
			},
			InstanceData: this.Positions,
			Mode:         renderstuff.TriangleFan,
		}
	}
	return this.mesh
}
