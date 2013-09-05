package main

import (
	"github.com/go-gl/gl"
	"github.com/krux02/mathgl"
	"math/rand"
	"sort"
	"unsafe"
)
//ipmort "fmt"

type PalmTrees struct {
	positions  []mathgl.Vec3f
	sortedX    []int
	sortedY    []int
	sortedXInv []int
	sortedYInv []int
}

func (pt *PalmTrees) CreateBuffers() (gl.VertexArray, gl.Buffer, gl.Buffer) {
	vao := gl.GenVertexArray()
	vao.Bind()
	vertices := gl.GenBuffer()
	vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(PalmTrees{}))*len(pt.positions), pt.positions, gl.STATIC_DRAW)
	indices := gl.GenBuffer()
	indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	size := int(unsafe.Sizeof(int(0))) * len(pt.positions)
	gl.BufferData(gl.ARRAY_BUFFER, 4*size, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0*size, size, pt.sortedX)
	gl.BufferSubData(gl.ARRAY_BUFFER, 1*size, size, pt.sortedY)
	gl.BufferSubData(gl.ARRAY_BUFFER, 2*size, size, pt.sortedXInv)
	gl.BufferSubData(gl.ARRAY_BUFFER, 3*size, size, pt.sortedYInv)
	return vao, vertices, indices
}

type TreeRenderLocatins struct {
	vertex_os, texCoord, position_cs gl.AttribLocation
	Proj, image                      gl.UniformLocation
}

func (pt *PalmTrees) CreateProgram() {
	prog := MakeProgram("Sprite.vs", "Sprite.fs")
	locations := TreeRenderLocatins{}
	BindLocations(prog, &locations)
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

func NewPalmTrees(word *HeightMap, count int) *PalmTrees {

	pt := PalmTrees{
		make([]mathgl.Vec3f, count),
		make([]int, count),
		make([]int, count),
		make([]int, count),
		make([]int, count),
	}

	for i := 0; i < count; i++ {
		x := rand.Float32() * float32(word.W)
		y := rand.Float32() * float32(word.H)
		z := word.Get2f(x, y)
		pt.positions[i] = mathgl.Vec3f{x, y, z}
		pt.sortedX[i] = i
		pt.sortedY[i] = i
		pt.sortedXInv[i] = i
		pt.sortedYInv[i] = i
	}

	sorterX := &TreeSorter{
		pt.sortedX,
		func(a, b int) bool {
			return pt.positions[a][0] < pt.positions[b][0]
		},
	}
	sort.Sort(sorterX)
	sorterX.indices = pt.sortedXInv
	sort.Sort(sort.Reverse(sorterX))

	sorterY := &TreeSorter{
		pt.sortedY,
		func(a, b int) bool {
			return pt.positions[a][1] < pt.positions[b][1]
		},
	}

	sort.Sort(sorterY)
	sorterY.indices = pt.sortedYInv
	sort.Sort(sort.Reverse(sorterY))

	return &pt
}
