package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type ScreenQuadVertex struct {
	Vertex_ndc mgl.Vec4f
}

var ScreenQuadMesh = renderstuff.Mesh{
	Vertices: []ScreenQuadVertex{
		ScreenQuadVertex{mgl.Vec4f{-1, -1, 0, 1}},
		ScreenQuadVertex{mgl.Vec4f{3, -1, 0, 1}},
		ScreenQuadVertex{mgl.Vec4f{-1, 3, 0, 1}},
	},
	Mode: renderstuff.Triangles,
}

type ScreenQuad struct{}

func (this *ScreenQuad) Mesh() *renderstuff.Mesh {
	return &ScreenQuadMesh
}

func (this *ScreenQuad) Model() mgl.Mat4f {
	return mgl.Ident4f()
}
