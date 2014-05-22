package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type PalmRenderer struct {
	Prog    gl.Program
	Loc     RenderLocations
	RenData RenderData
}

// global information for all trees
type PalmShape struct {
	Vertex_os mgl.Vec4f
	TexCoord  mgl.Vec2f
}

type PalmTreeFullVertex struct {
	Vertex_ws mgl.Vec4f
	TexCoord  mgl.Vec2f
}

func (renderer *PalmRenderer) CreateRenderData(pt *gamestate.PalmTreesInstanceData) {
	renderer.RenData.VAO = gl.GenVertexArray()
	renderer.RenData.VAO.Bind()

	vertices, indices, numverts := CreateVertexDataBuffer()
	helpers.SetAttribPointers(&renderer.Loc, &PalmShape{})

	instanceDataBuffer := CreateInstanceDataBuffer(pt)
	helpers.SetAttribPointers(&renderer.Loc, &gamestate.PalmTree{})
	renderer.Loc.InstancePosition_ws.AttribDivisor(1)

	renderer.RenData.InstanceDataBuffer = instanceDataBuffer
	renderer.RenData.NumInstances = len(pt.Positions)
	renderer.RenData.Vertices = vertices
	renderer.RenData.Indices = indices
	renderer.RenData.Numverts = numverts
}

func CreateVertexDataBuffer() (vertices, indices gl.Buffer, numverts int) {
	fmt.Println("CreateVertexDataBuffer:")

	palmShape := []PalmShape{
		PalmShape{mgl.Vec4f{0, 1, 2, 1}, mgl.Vec2f{1, 0}},
		PalmShape{mgl.Vec4f{0, 1, 0, 1}, mgl.Vec2f{1, 1}},
		PalmShape{mgl.Vec4f{0, -1, 0, 1}, mgl.Vec2f{0, 1}},
		PalmShape{mgl.Vec4f{0, -1, 2, 1}, mgl.Vec2f{0, 0}},
	}

	vertices = gl.GenBuffer()
	vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(palmShape), palmShape, gl.STATIC_DRAW)

	indices = gl.GenBuffer()
	indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 8, &[4]uint16{0, 1, 2, 3}, gl.STATIC_DRAW)

	numverts = 4

	return
}

func CreateInstanceDataBuffer(pt *gamestate.PalmTreesInstanceData) gl.Buffer {
	fmt.Println("CreateInstanceDataBuffer:")
	vertices := gl.GenBuffer()
	vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(pt.Positions), pt.Positions, gl.STATIC_DRAW)

	// fmt.Println(pt.positions)
	return vertices
}

func (this *PalmRenderer) Delete() {
	this.Prog.Delete()
	this.RenData.Delete()
	*this = PalmRenderer{}
}

func NewPalmRenderer(pt *gamestate.PalmTreesInstanceData) *PalmRenderer {
	renderer := new(PalmRenderer)
	renderer.Prog = helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	renderer.Prog.Use()
	helpers.BindLocations("palm sprite", renderer.Prog, &renderer.Loc)
	renderer.Loc.PalmTree.Uniform1i(5)

	renderer.CreateRenderData(pt)

	return renderer
}

func (pt *PalmRenderer) Render(Proj, View mgl.Mat4f, Rot2D mgl.Mat3f, clippingPlane mgl.Vec4f) {
	pt.Prog.Use()
	pt.RenData.VAO.Bind()
	pt.Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	pt.Loc.View.UniformMatrix4f(false, glMat4(&View))
	pt.Loc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
	pt.Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])
	gl.DrawArraysInstanced(gl.TRIANGLE_FAN, 0, pt.RenData.Numverts, pt.RenData.NumInstances)
}
