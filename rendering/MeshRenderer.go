package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type MeshRenderer struct {
	Program gl.Program
	RenLoc  MeshRenderLocations
	RenData map[*gamestate.Mesh]*MeshRenderData
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
	mr.RenData = make(map[*gamestate.Mesh]*MeshRenderData)
	helpers.BindLocations("mesh", mr.Program, &mr.RenLoc)
	return
}

func (this *MeshRenderer) Delete() {
	this.Program.Delete()
	for _, rd := range this.RenData {
		rd.Delete()
	}
	*this = MeshRenderer{}
}

func (this *MeshRenderer) CreateRenderData(mesh *gamestate.Mesh) (rd MeshRenderData) {

	rd.VAO = gl.GenVertexArray()
	rd.VAO.Bind()

	rd.Indices = gl.GenBuffer()
	rd.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Indices), mesh.Indices, gl.STATIC_DRAW)

	rd.Vertices = gl.GenBuffer()
	rd.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Vertices), mesh.Vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&this.RenLoc, &gamestate.MeshVertex{})

	rd.Numverts = len(mesh.Indices)

	return
}

func (this *MeshRenderer) Render(mesh *gamestate.Mesh, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f) {
	this.Program.Use()

	meshData := this.RenData[mesh]
	if meshData == nil {
		md := this.CreateRenderData(mesh)
		meshData = &md
		this.RenData[mesh] = &md
	}

	meshData.VAO.Bind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Model.UniformMatrix4f(false, glMat4(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))

	numverts := meshData.Numverts

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
