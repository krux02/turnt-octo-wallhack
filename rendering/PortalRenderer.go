package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/world"
)

type PortalRenderer struct {
	Program gl.Program
	RenLoc  PortalRenderLocations
}

type PortalRenderLocations struct {
	Vertex_ms, Normal_ms       gl.AttribLocation
	Proj, View, Model, U_Image gl.UniformLocation
}

type PortalRenderData struct {
	VAO      gl.VertexArray
	Indices  gl.Buffer
	Vertices gl.Buffer
	Numverts int
}

func NewPortalRenderer() (mr *PortalRenderer) {
	mr = new(PortalRenderer)
	mr.Program = helpers.MakeProgram("Portal.vs", "Portal.fs")
	mr.Program.Use()
	helpers.BindLocations(mr.Program, &mr.RenLoc)
	helpers.PrintLocations(&mr.RenLoc)
	return
}

func (this *PortalRenderer) Delete() {
	this.Program.Delete()
	this.Program = 0
}

func (this *PortalRenderer) CreateRenderData(mesh *world.Mesh) (md PortalRenderData) {

	md.VAO = gl.GenVertexArray()
	md.VAO.Bind()

	md.Indices = gl.GenBuffer()
	md.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Indices), mesh.Indices, gl.STATIC_DRAW)

	md.Vertices = gl.GenBuffer()
	md.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Vertices), mesh.Vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&this.RenLoc, &world.MeshVertex{}, true)

	md.Numverts = len(mesh.Indices)

	return
}

func glMat(mat *mgl.Mat4f) *[16]float32 {
	return (*[16]float32)(mat)
}

func (this *PortalRenderer) Render(meshData *PortalRenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, textureUnit int) {
	this.Program.Use()
	meshData.VAO.Bind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat(&View))
	Loc.Model.UniformMatrix4f(false, glMat(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat(&Proj))
	Loc.U_Image.Uniform1i(textureUnit)

	numverts := meshData.Numverts

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
