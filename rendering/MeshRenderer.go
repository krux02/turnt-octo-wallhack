package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"reflect"
)

type MeshRenderer struct {
	Program gl.Program
	RenLoc  RenderLocations
}

func NewMeshRenderer() (mr *MeshRenderer) {
	mr = new(MeshRenderer)
	mr.Program = helpers.MakeProgram("Mesh.vs", "Mesh.fs")
	mr.Program.Use()
	helpers.BindLocations("mesh", mr.Program, &mr.RenLoc)
	return
}

func (this *MeshRenderer) Delete() {
	this.Program.Delete()
	*this = MeshRenderer{}
}

func CreateMeshRenderData(mesh gamestate.IMesh, renLoc *RenderLocations) (rd RenderData) {
	vertices, indices := mesh.CreateVertexArray()
	verticesType := reflect.TypeOf(vertices)
	if verticesType.Kind() != reflect.Slice {
		panic("vertices is not a slice")
	}
	indicesType := reflect.TypeOf(indices)
	if indicesType.Kind() != reflect.Slice {
		panic("indices is not a slice")
	}

	rd.VAO = gl.GenVertexArray()
	rd.VAO.Bind()

	rd.Indices = gl.GenBuffer()
	rd.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	rd.Vertices = gl.GenBuffer()
	rd.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)

	verticesValue := reflect.ValueOf(vertices)
	indicesValue := reflect.ValueOf(indices)
	rd.Numverts = indicesValue.Len()
	vertex := verticesValue.Index(0).Addr().Interface()
	helpers.SetAttribPointers(renLoc, vertex)

	return
}

func (this *MeshRenderer) Render(meshData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f) {
	this.Program.Use()

	meshData.VAO.Bind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	Loc := this.RenLoc
	Loc.Model.UniformMatrix4f(false, glMat4(&Model))
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	Loc.ClippingPlane_ws.Uniform4f(ClippingPlane_ws[0], ClippingPlane_ws[1], ClippingPlane_ws[2], ClippingPlane_ws[3])

	numverts := meshData.Numverts

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_SHORT, uintptr(0))
}

func (this *MeshRenderer) UseProgram() {
	this.Program.Use()
}

func (this *MeshRenderer) RenderLocations() *RenderLocations {
	return &this.RenLoc
}

func (this *MeshRenderer) Update(entiy gamestate.IRenderEntity) {}
