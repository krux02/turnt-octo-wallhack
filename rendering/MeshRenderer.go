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
	RenLoc  RenderLocations
	RenData map[*gamestate.Mesh]*RenderData
}

func NewMeshRenderer() (mr *MeshRenderer) {
	mr = new(MeshRenderer)
	mr.Program = helpers.MakeProgram("Mesh.vs", "Mesh.fs")
	mr.Program.Use()
	mr.RenData = make(map[*gamestate.Mesh]*RenderData)
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

func (this *MeshRenderer) CreateRenderData(mesh *gamestate.Mesh) (rd RenderData) {

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

func (this *MeshRenderer) Render(mesh *gamestate.Mesh, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f) {
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
	Loc.ClippingPlane_ws.Uniform4f(ClippingPlane_ws[0], ClippingPlane_ws[1], ClippingPlane_ws[2], ClippingPlane_ws[3])

	numverts := meshData.Numverts

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
