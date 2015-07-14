package rendering

import (
	//"fmt"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type SkyboxVertex struct {
	InTexCoord mgl.Vec3
}

var SkyboxMesh = renderstuff.Mesh{
	Vertices: []SkyboxVertex{
		{mgl.Vec3{-1, -1, -1}},
		{mgl.Vec3{1, -1, -1}},
		{mgl.Vec3{-1, 1, -1}},
		{mgl.Vec3{1, 1, -1}},
		{mgl.Vec3{-1, -1, 1}},
		{mgl.Vec3{1, -1, 1}},
		{mgl.Vec3{-1, 1, 1}},
		{mgl.Vec3{1, 1, 1}},
	},
	Indices: []uint16{
		0, 2, 1, 1, 2, 3,
		0, 4, 6, 0, 6, 2,
		0, 5, 4, 0, 1, 5,
		4, 5, 6, 6, 5, 7,
		5, 1, 3, 5, 3, 7,
		6, 7, 2, 2, 7, 3,
	},
	Mode: renderstuff.Triangles,
}

type Skybox struct{}

func (this *Skybox) Model() mgl.Mat4 {
	return mgl.Ident4()
}

func (this *Skybox) Mesh() *renderstuff.Mesh {
	return &SkyboxMesh
}
