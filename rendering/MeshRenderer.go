package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type MeshRenderer struct {
	Program  gl.Program
	RenLoc   MeshRenderLocations
	MeshData map[*gamestate.Mesh]*MeshRenderData
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

func (this *MeshRenderData) Delete() {
	this.VAO.Delete()
	this.Indices.Delete()
	this.Vertices.Delete()
	*this = MeshRenderData{}
}

func NewMeshRenderer() (mr *MeshRenderer) {
	mr = new(MeshRenderer)
	mr.Program = helpers.MakeProgram("Mesh.vs", "Mesh.fs")
	mr.Program.Use()
	mr.MeshData = make(map[*gamestate.Mesh]*MeshRenderData)
	helpers.BindLocations("mesh", mr.Program, &mr.RenLoc)
	return
}

func (this *MeshRenderer) Delete() {
	this.Program.Delete()
	for _, rd := range this.MeshData {
		rd.Delete()
	}
	*this = MeshRenderer{}
}

func (this *MeshRenderer) createRenderData(mesh *gamestate.Mesh) (md MeshRenderData) {

	md.VAO = gl.GenVertexArray()
	md.VAO.Bind()

	md.Indices = gl.GenBuffer()
	md.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Indices), mesh.Indices, gl.STATIC_DRAW)

	md.Vertices = gl.GenBuffer()
	md.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Vertices), mesh.Vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&this.RenLoc, &gamestate.MeshVertex{})

	md.Numverts = len(mesh.Indices)

	return
}

func (this *MeshRenderer) Render(mesh *gamestate.Mesh, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f) {
	this.Program.Use()

	meshData := this.MeshData[mesh]
	if meshData == nil {
		md := this.createRenderData(mesh)
		meshData = &md
		this.MeshData[mesh] = &md
	}

	meshData.VAO.Bind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat(&View))
	Loc.Model.UniformMatrix4f(false, glMat(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat(&Proj))

	numverts := meshData.Numverts

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
