package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

type ScreenQuadVertex struct {
	Vertex_ndc mgl.Vec4f
}

type ScreenQuad struct {
	gamestate.AbstractMesh
	vertices []ScreenQuadVertex
	indices  []uint16
}

func (this *ScreenQuad) Init() *ScreenQuad {
	this.vertices = []ScreenQuadVertex{
		ScreenQuadVertex{mgl.Vec4f{-1, -1, 0, 1}},
		ScreenQuadVertex{mgl.Vec4f{3, -1, 0, 1}},
		ScreenQuadVertex{mgl.Vec4f{-1, 3, 0, 1}},
	}
	return this
}

func (this *ScreenQuad) GetModel() mgl.Mat4f {
	return mgl.Ident4f()
}

func (this *ScreenQuad) GetMesh() gamestate.IMesh {
	return this
}

func (this *ScreenQuad) Vertices() interface{} {
	return this.vertices
}

func (this *ScreenQuad) Indices() interface{} {
	return nil
}

func (this *ScreenQuad) Mode() gamestate.Mode {
	return gamestate.Triangles
}
