package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
	"github.com/krux02/mathgl"
	"math/rand"
	"sort"
	"unsafe"
)

import "fmt"

type TreeRenderLocatins struct {
	Vertex_os, TexCoord, Position_ws gl.AttribLocation
	Proj, View, PalmTree             gl.UniformLocation
}

// global information for all trees
type PalmShape struct {
	Vertex_os mathgl.Vec4f
	TexCoord  mathgl.Vec2f
}

// instance data for each tree
type PalmTree struct {
	Position_ws mathgl.Vec3f
}

func CreateVertexDataBuffer() gl.Buffer {

	palmShape := []PalmShape{
		PalmShape{mathgl.Vec4f{-1, 0, 0, 1}, mathgl.Vec2f{0, 0}},
		PalmShape{mathgl.Vec4f{-1, 2, 0, 1}, mathgl.Vec2f{0, 1}},
		PalmShape{mathgl.Vec4f{1, 2, 0, 1}, mathgl.Vec2f{1, 1}},
		PalmShape{mathgl.Vec4f{1, 0, 0, 1}, mathgl.Vec2f{1, 0}},
	}

	palmShapeBuffer := gl.GenBuffer()
	palmShapeBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(palmShape)*int(unsafe.Sizeof(PalmShape{})), palmShape, gl.STATIC_DRAW)

	return palmShapeBuffer
}

func (pt *PalmTreesInstanceData) CreateInstanceDataBuffer() gl.Buffer {
	vertices := gl.GenBuffer()
	vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(PalmTree{}))*len(pt.positions), pt.positions, gl.STATIC_DRAW)
	return vertices
}

func (pt *PalmTreesInstanceData) CreateIndexDataBuffer() gl.Buffer {
	indices := gl.GenBuffer()
	indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	size := int(unsafe.Sizeof(int(0))) * len(pt.positions)
	gl.BufferData(gl.ARRAY_BUFFER, 4*size, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0*size, size, pt.sortedX)
	gl.BufferSubData(gl.ARRAY_BUFFER, 1*size, size, pt.sortedY)
	gl.BufferSubData(gl.ARRAY_BUFFER, 2*size, size, pt.sortedXInv)
	gl.BufferSubData(gl.ARRAY_BUFFER, 3*size, size, pt.sortedYInv)
	return indices
}

type TreeSorter struct {
	indices []int
	by      func(a, b int) bool
}

func (ts *TreeSorter) Len() int {
	return len(ts.indices)
}

func (ts *TreeSorter) Less(i, j int) bool {
	return ts.by(ts.indices[i], ts.indices[j])
}

func (ts *TreeSorter) Swap(i, j int) {
	ts.indices[i], ts.indices[j] = ts.indices[j], ts.indices[i]
}

type PalmTreesInstanceData struct {
	positions  []PalmTree
	sortedX    []int
	sortedY    []int
	sortedXInv []int
	sortedYInv []int
}

func NewPalmTreesInstanceData(world *HeightMap, count int) *PalmTreesInstanceData {

	pt := &PalmTreesInstanceData{
		make([]PalmTree, count),
		make([]int, count),
		make([]int, count),
		make([]int, count),
		make([]int, count),
	}

	for i := 0; i < count; i++ {
		x := rand.Float32() * float32(world.W)
		y := rand.Float32() * float32(world.H)
		z := world.Get2f(x, y)
		pt.positions[i] = PalmTree{mathgl.Vec3f{x, y, z}}
		pt.sortedX[i] = i
		pt.sortedY[i] = i
		pt.sortedXInv[i] = i
		pt.sortedYInv[i] = i
	}

	sorterX := &TreeSorter{
		pt.sortedX,
		func(a, b int) bool {
			return pt.positions[a].Position_ws[0] < pt.positions[b].Position_ws[0]
		},
	}
	sort.Sort(sorterX)
	sorterX.indices = pt.sortedXInv
	sort.Sort(sort.Reverse(sorterX))

	sorterY := &TreeSorter{
		pt.sortedY,
		func(a, b int) bool {
			return pt.positions[a].Position_ws[1] < pt.positions[b].Position_ws[1]
		},
	}

	sort.Sort(sorterY)
	sorterY.indices = pt.sortedYInv
	sort.Sort(sort.Reverse(sorterY))

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
	Count 	int
}

func NewPalmTrees(world *HeightMap, count int) *PalmTrees {
	pt := NewPalmTreesInstanceData(world, count)
	Prog := glh.NewProgram(
		glh.Shader{gl.VERTEX_SHADER, ReadShaderFile("Sprite.vs")},
		glh.Shader{gl.FRAGMENT_SHADER, ReadShaderFile("Sprite.fs")},
	)

	Prog.Use()

	vao := gl.GenVertexArray()
	vao.Bind()

	Loc := TreeRenderLocatins{}
	BindLocations(Prog, &Loc)

	fmt.Println(Loc)
	Loc.PalmTree.Uniform1i(3)

	vertexDataBuffer := CreateVertexDataBuffer()
	SetAttribPointers(&Loc, &PalmShape{})

	instanceDataBuffer := pt.CreateInstanceDataBuffer()
	SetAttribPointers(&Loc, &PalmTree{})
	Loc.Vertex_os.AttribDivisor(1)
	Loc.TexCoord.AttribDivisor(1)

	return &PalmTrees{Prog, Loc, PalmTreesBuffers{vao, instanceDataBuffer, vertexDataBuffer}, count}
}

func (pt* PalmTrees) Render(Proj,View mathgl.Mat4f) {

	pt.Prog.Use()

	pt.Loc.Proj.UniformMatrix4f(false, (*[16]float32)(&Proj))
	pt.Loc.View.UniformMatrix4f(false, (*[16]float32)(&View))

	gl.DrawArraysInstanced(gl.TRIANGLE_FAN, 0, 4, pt.Count)
}