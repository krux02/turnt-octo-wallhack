package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type Skybox struct {
	renderstuff.AbstractMesh
	vertices []SkyboxVertex
	indices  []uint16
}

type SkyboxVertex struct {
	InTexCoord mgl.Vec3f
}

func (this *Skybox) GetModel() mgl.Mat4f {
	return mgl.Ident4f()
}

func (this *Skybox) GetMesh() renderstuff.IMesh {
	return this
}

func (this *Skybox) Vertices() interface{} {
	return this.vertices
}

func (this *Skybox) Indices() interface{} {
	return this.indices
}

func (this *Skybox) Mode() renderstuff.Mode {
	return renderstuff.Triangles
}

func (this *Skybox) Init() *Skybox {
	this.vertices = []SkyboxVertex{
		SkyboxVertex{mgl.Vec3f{-1, -1, -1}},
		SkyboxVertex{mgl.Vec3f{1, -1, -1}},
		SkyboxVertex{mgl.Vec3f{-1, 1, -1}},
		SkyboxVertex{mgl.Vec3f{1, 1, -1}},
		SkyboxVertex{mgl.Vec3f{-1, -1, 1}},
		SkyboxVertex{mgl.Vec3f{1, -1, 1}},
		SkyboxVertex{mgl.Vec3f{-1, 1, 1}},
		SkyboxVertex{mgl.Vec3f{1, 1, 1}},
	}
	this.indices = []uint16{
		0, 2, 1, 1, 2, 3,
		0, 4, 6, 0, 6, 2,
		0, 5, 4, 0, 1, 5,
		4, 5, 6, 6, 5, 7,
		5, 1, 3, 5, 3, 7,
		6, 7, 2, 2, 7, 3,
	}
	return this
}
