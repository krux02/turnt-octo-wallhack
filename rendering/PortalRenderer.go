package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type PortalRenderer struct {
	Program gl.Program
	RenLoc  RenderLocations
	RenData RenderData
}

func NewPortalRenderer() (mr *PortalRenderer) {
	mr = new(PortalRenderer)
	mr.Program = helpers.MakeProgram("Portal.vs", "Portal.fs")
	mr.Program.Use()
	helpers.BindLocations("portal", mr.Program, &mr.RenLoc)
	mr.RenData = CreatePortalRenderData()
	helpers.SetAttribPointers(&mr.RenLoc, &gamestate.MeshVertex{}, false)
	return
}

func (this *PortalRenderer) Delete() {
	this.Program.Delete()
	this.RenData.Delete()
	*this = PortalRenderer{}
}

func CreatePortalRenderData() (md RenderData) {
	mesh := gamestate.QuadMesh()
	md.VAO = gl.GenVertexArray()
	md.VAO.Bind()

	md.Indices = gl.GenBuffer()
	md.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Indices()), mesh.Indices(), gl.STATIC_DRAW)

	md.Vertices = gl.GenBuffer()
	md.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(mesh.Vertices()), mesh.Vertices(), gl.STATIC_DRAW)

	md.Numverts = len((mesh.Indices()).([]gamestate.MeshIndex))

	return
}

func (this *PortalRenderer) Render(Portal *gamestate.Portal, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f, TextureUnit int) {
	this.Program.Use()
	this.RenData.VAO.Bind()

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Model.UniformMatrix4f(false, glMat4(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	Loc.Image.Uniform1i(TextureUnit)
	Loc.ClippingPlane_ws.Uniform4f(ClippingPlane_ws[0], ClippingPlane_ws[1], ClippingPlane_ws[2], ClippingPlane_ws[3])

	numverts := this.RenData.Numverts

	gl.Enable(gl.DEPTH_CLAMP)
	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}
