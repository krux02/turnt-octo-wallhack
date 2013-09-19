package rendering

import "github.com/krux02/turnt-octo-wallhack/world"
import "github.com/krux02/turnt-octo-wallhack/helpers"
import "github.com/go-gl/gl"

type MeshRenderer struct {
	Program  gl.Program
	RenLoc   MeshRenderLocations
}

type MeshRenderLocations struct {
	Position,Normal gl.AttribLocation
}

type MeshRenderData struct {
	VAO      gl.VertexArray
	Indices  gl.Buffer
	Vertices gl.Buffer
	Numverts int
}

// creates and activates a new Program
func NewMeshRenderer() (mr MeshRenderer) {
	mr.Program = helpers.MakeProgram("World.vs", "World.fs")
	mr.Program.Use()
	helpers.BindLocations(mr.Program, &mr.RenLoc)
	return
}

func (this *MeshRenderer) CreateMeshRenderData(mesh *world.Mesh) (md MeshRenderData) {

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

func (this *MeshRenderer) Render(meshData *MeshRenderData) {
}
