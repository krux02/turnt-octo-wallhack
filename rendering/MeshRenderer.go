package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/world"
)

type MeshRenderer struct {
	Program gl.Program
	RenLoc  MeshRenderLocations
}

type MeshRenderLocations struct {
	Vertex_ms, Normal_ms gl.AttribLocation
	Proj, View, Model    gl.UniformLocation
}

type MeshRenderData struct {
	VAO      gl.VertexArray
	Indices  gl.Buffer
	Vertices gl.Buffer
	Numverts int
}

func NewMeshRenderer() (mr *MeshRenderer) {
	mr = new(MeshRenderer)
	mr.Program = helpers.MakeProgram("Mesh.vs", "Mesh.fs")
	mr.Program.Use()
	helpers.BindLocations(mr.Program, &mr.RenLoc)
	helpers.PrintLocations(&mr.RenLoc)
	return
}

func (this *MeshRenderer) Delete() {
	this.Program.Delete()
	this.Program = 0
}

func (this *MeshRenderer) CreateRenderData(mesh *world.Mesh) (md MeshRenderData) {

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

func (this *MeshRenderer) Render(meshData *MeshRenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f) {
	this.Program.Use()
	meshData.VAO.Bind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, (*[16]float32)(&View))
	Loc.Model.UniformMatrix4f(false, (*[16]float32)(&Model))
	Loc.Proj.UniformMatrix4f(false, (*[16]float32)(&Proj))

	numverts := meshData.Numverts

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
