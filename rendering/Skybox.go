package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type SkyboxRenderer struct {
	Program gl.Program
	RenLoc  SkyboxRenderLocations
	RenData RenderData
}

type SkyboxRenderLocations struct {
	Proj, View, Skybox gl.UniformLocation
	InTexCoord         gl.AttribLocation
}

func NewSkyboxRenderer() *SkyboxRenderer {
	renderer := new(SkyboxRenderer)
	renderer.Program = helpers.MakeProgram("Skybox.vs", "Skybox.fs")
	renderer.Program.Use()
	renderer.CreateRenderData()
	helpers.BindLocations("skybox", renderer.Program, &renderer.RenLoc)
	return renderer
}

func (this *SkyboxRenderer) Delete() {
	this.Program.Delete()
	this.RenData.Indices.Delete()
	this.RenData.VAO.Delete()
	*this = SkyboxRenderer{}
}

type SkyboxVertex struct {
	InTexCoord mgl.Vec3f
}

func (this *SkyboxRenderer) CreateRenderData() {
	this.RenData.VAO = gl.GenVertexArray()
	this.RenData.VAO.Bind()
	this.RenData.Indices = gl.GenBuffer()
	this.RenData.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	indices := []uint16{
		0, 2, 1, 1, 2, 3,
		0, 4, 6, 0, 6, 2,
		0, 5, 4, 0, 1, 5,
		4, 5, 6, 6, 5, 7,
		5, 1, 3, 5, 3, 7,
		6, 7, 2, 2, 7, 3,
	}
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	vertices := []mgl.Vec3f{
		mgl.Vec3f{-1, -1, -1},
		mgl.Vec3f{1, -1, -1},
		mgl.Vec3f{-1, 1, -1},
		mgl.Vec3f{1, 1, -1},
		mgl.Vec3f{-1, -1, 1},
		mgl.Vec3f{1, -1, 1},
		mgl.Vec3f{-1, 1, 1},
		mgl.Vec3f{1, 1, 1},
	}

	this.RenData.Vertices = gl.GenBuffer()
	this.RenData.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)
	this.RenData.Numverts = 36

	helpers.SetAttribPointers(&this.RenLoc, &SkyboxVertex{})
	return
}

func (this *SkyboxRenderer) Render(Proj mgl.Mat4f, View mgl.Mat4f, textureUnit int) {
	this.Program.Use()
	this.RenData.VAO.Bind()

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	Loc.Skybox.Uniform1i(textureUnit)

	gl.DrawElements(gl.TRIANGLES, this.RenData.Numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
