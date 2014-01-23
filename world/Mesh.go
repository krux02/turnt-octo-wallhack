package world

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	ai "github.com/krux02/assimp"
	"math"
)

type MeshIndex uint16

type MeshVertex struct {
	Vertex_ms mgl.Vec4f
	Normal_ms mgl.Vec4f
}

type Mesh struct {
	Vertices []MeshVertex
	Indices  []MeshIndex
}

func QuadMesh() (mesh *Mesh) {
	mesh = &Mesh{
		[]MeshVertex{
			MeshVertex{mgl.Vec4f{-1, -1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{1, -1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{1, 1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{-1, 1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
		},
		[]MeshIndex{0, 1, 2, 2, 3, 0},
	}
	return
}

func Min(v1, v2 mgl.Vec4f) (min mgl.Vec4f) {
	for i := 0; i < 4; i++ {
		if v1[i] < v2[i] {
			min[i] = v1[i]
		} else {
			min[i] = v2[i]
		}
	}
	return
}

func Max(v1, v2 mgl.Vec4f) (min mgl.Vec4f) {
	for i := 0; i < 4; i++ {
		if v1[i] < v2[i] {
			min[i] = v2[i]
		} else {
			min[i] = v1[i]
		}
	}
	return
}

func LoadMesh(filename string) (mesh *Mesh) {
	scene := ai.ImportFile(filename, 0)
	scene.ApplyPostProcessing(ai.Process_Triangulate)
	meshes := scene.Meshes()
	if len(meshes) != 1 {
		panic(fmt.Sprintf("not a single mesh, found %d meshes", len(meshes)))
	}
	aimesh := meshes[0]

	meshVertices := make([]MeshVertex, aimesh.NumVertices())
	for i, pos := range aimesh.Vertices() {
		v := pos.Values()
		meshVertices[i].Vertex_ms = mgl.Vec4f([4]float32{v[0], v[1], v[2], 1})
	}
	for i, norm := range aimesh.Normals() {
		n := norm.Values()
		meshVertices[i].Normal_ms = mgl.Vec4f([4]float32{n[0], n[1], n[2], 0})
	}

	meshIndices := make([]MeshIndex, aimesh.NumFaces()*3)
	for i, face := range aimesh.Faces() {
		indices := face.CopyIndices()
		meshIndices[i*3+0] = MeshIndex(indices[0])
		meshIndices[i*3+1] = MeshIndex(indices[1])
		meshIndices[i*3+2] = MeshIndex(indices[2])
	}

	fmt.Println("loaded mesh with", aimesh.NumVertices(), "Vertices")

	return &Mesh{meshVertices, meshIndices}
}

func (m *Mesh) BoundingBox() (min mgl.Vec4f, max mgl.Vec4f) {
	min = mgl.Vec4f{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max = mgl.Vec4f{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}

	for _, v := range m.Vertices {
		min = Min(min, v.Vertex_ms)
		max = Max(max, v.Vertex_ms)
	}
	return
}

// returns the 8 vertices of the box that is defined by two if it's vertices
func (m *Mesh) MakeBoxVertices() (verts [8]mgl.Vec4f) {
	bottomLeft, topRight := m.BoundingBox()

	var i int
	bounds := [2]mgl.Vec4f{bottomLeft, topRight}
	for x := 0; x <= 1; x++ {
		for y := 0; y <= 1; y++ {
			for z := 0; z <= 1; z++ {
				verts[i][0] = bounds[x][0]
				verts[i][1] = bounds[y][1]
				verts[i][2] = bounds[z][2]
				verts[i][3] = 1
				i++
			}
		}
	}
	return
}
