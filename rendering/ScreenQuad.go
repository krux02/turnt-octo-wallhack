package rendering

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type ScreenQuadVertex struct {
	Vertex_ndc mgl.Vec4
}

var ScreenQuadMesh = renderstuff.Mesh{
	Vertices: []ScreenQuadVertex{
		{mgl.Vec4{-1, -1, 0, 1}},
		{mgl.Vec4{3, -1, 0, 1}},
		{mgl.Vec4{-1, 3, 0, 1}},
	},
	Mode: renderstuff.Triangles,
}

type ScreenQuad struct{}

func (this *ScreenQuad) Mesh() *renderstuff.Mesh {
	return &ScreenQuadMesh
}

func (this *ScreenQuad) Model() mgl.Mat4 {
	return mgl.Ident4()
}
