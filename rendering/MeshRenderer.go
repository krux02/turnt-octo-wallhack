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
	rd.VAO = gl.GenVertexArray()
	rd.VAO.Bind()

	{
		vertices := mesh.Vertices()
		verticesType := reflect.TypeOf(vertices)
		if verticesType.Kind() != reflect.Slice {
			panic("Vertices is not a slice")
		}
		rd.Vertices = gl.GenBuffer()
		rd.Vertices.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)
		rd.Numverts = reflect.ValueOf(vertices).Len()
		helpers.SetAttribPointers(renLoc, reflect.ValueOf(vertices).Index(0).Addr().Interface(), false)
	}

	if indices := mesh.Indices(); indices != nil {
		indicesType := reflect.TypeOf(indices)
		if indicesType.Kind() != reflect.Slice {
			panic("Indices is not a slice")
		}
		rd.Indices = gl.GenBuffer()
		rd.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)
		rd.Numverts = reflect.ValueOf(indices).Len()
	}

	if instanceData := mesh.InstanceData(); instanceData != nil {
		Type := reflect.TypeOf(instanceData)
		if Type.Kind() != reflect.Slice {
			panic("InstanceData is not a slice")
		}
		rd.InstanceDataBuffer = gl.GenBuffer()
		rd.InstanceDataBuffer.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(instanceData), instanceData, gl.STATIC_DRAW)
		helpers.SetAttribPointers(renLoc, reflect.ValueOf(instanceData).Index(0).Addr().Interface(), true)

		rd.NumInstances = reflect.ValueOf(instanceData).Len()
	}

	return
}

func (this *MeshRenderer) Render(meshData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f, additionalUniforms map[string]int) {
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
