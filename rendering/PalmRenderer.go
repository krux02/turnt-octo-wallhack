package rendering

import (
	"github.com/go-gl/gl"
	mgl "github.com/krux02/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/world"
	"math/rand"
)

import "fmt"

type TreeRenderLocatins struct {
	Vertex_os, TexCoord, Position_ws gl.AttribLocation
	Proj, View, PalmTree, Rot2D      gl.UniformLocation
}

// global information for all trees
type PalmShape struct {
	Vertex_os mgl.Vec4f
	TexCoord  mgl.Vec2f
}

// instance data for each tree
type PalmTree struct {
	Position_ws mgl.Vec4f
}

type PalmTreeFullVertex struct {
	Vertex_ws mgl.Vec4f
	TexCoord  mgl.Vec2f
}

func CreateVertexDataBuffer() gl.Buffer {
	fmt.Println("CreateVertexDataBuffer:")

	palmShape := []PalmShape{
		PalmShape{mgl.Vec4f{0, 1, 2, 1}, mgl.Vec2f{1, 0}},
		PalmShape{mgl.Vec4f{0, 1, 0, 1}, mgl.Vec2f{1, 1}},
		PalmShape{mgl.Vec4f{0, -1, 0, 1}, mgl.Vec2f{0, 1}},
		PalmShape{mgl.Vec4f{0, -1, 2, 1}, mgl.Vec2f{0, 0}},
	}

	palmShapeBuffer := gl.GenBuffer()
	palmShapeBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(palmShape), palmShape, gl.STATIC_DRAW)

	return palmShapeBuffer
}

func (pt *PalmTreesInstanceData) CreateInstanceDataBuffer() gl.Buffer {
	fmt.Println("CreateInstanceDataBuffer:")
	vertices := gl.GenBuffer()
	vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(pt.positions), pt.positions, gl.STATIC_DRAW)

	// fmt.Println(pt.positions)
	return vertices
}

type TreeSorter struct {
	indices []int
	by      func(a, b int) bool
}

type PalmTreesInstanceData struct {
	positions []PalmTree
}

func NewPalmTreesInstanceData(world *world.HeightMap, count int) *PalmTreesInstanceData {

	pt := &PalmTreesInstanceData{
		make([]PalmTree, count),
	}

	for i := 0; i < count; i++ {

		var x, y float32
		for true {
			x = rand.Float32() * float32(world.W)
			y = rand.Float32() * float32(world.H)
			n := world.Normalf(x, y)
			if n.Dot(mgl.Vec3f{0, 0, 1}) > 0.65 {
				break
			}
		}

		z := world.Get2f(x, y)

		pt.positions[i] = PalmTree{mgl.Vec4f{x, y, z, 1}}
	}

	return pt
}

type PalmTreesBuffers struct {
	Vao                gl.VertexArray
	InstanceDataBuffer gl.Buffer
	VertexDataBuffer   gl.Buffer
}

type PalmTrees struct {
	Prog    gl.Program
	Loc     TreeRenderLocatins
	Buffers PalmTreesBuffers
	Count   int
}

func NewPalmTrees(world *world.HeightMap, count int) *PalmTrees {
	pt := NewPalmTreesInstanceData(world, count)

	Prog := helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	Prog.Use()

	vao := gl.GenVertexArray()
	vao.Bind()

	Loc := TreeRenderLocatins{}
	helpers.BindLocations(Prog, &Loc)

	fmt.Println(Loc)
	Loc.PalmTree.Uniform1i(3)

	vertexDataBuffer := CreateVertexDataBuffer()
	helpers.SetAttribPointers(&Loc, &PalmShape{}, true)

	instanceDataBuffer := pt.CreateInstanceDataBuffer()
	helpers.SetAttribPointers(&Loc, &PalmTree{}, true)
	Loc.Position_ws.AttribDivisor(1)

	buffers := PalmTreesBuffers{vao, instanceDataBuffer, vertexDataBuffer}

	return &PalmTrees{Prog, Loc, buffers, count}
}

func (pt *PalmTrees) Render(Proj, View mgl.Mat4f, Rot2D mgl.Mat3f) {
	gl.Disable(gl.BLEND)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	pt.Prog.Use()
	pt.Buffers.Vao.Bind()

	pt.Loc.Proj.UniformMatrix4f(false, (*[16]float32)(&Proj))
	pt.Loc.View.UniformMatrix4f(false, (*[16]float32)(&View))
	pt.Loc.Rot2D.UniformMatrix3f(false, (*[9]float32)(&Rot2D))

	gl.DrawArraysInstanced(gl.TRIANGLE_FAN, 0, 4, pt.Count)
}
