package gamestate

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	ai "github.com/krux02/assimp"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"math"
)

type MeshIndex uint16

type MeshVertex struct {
	Vertex_ms mgl.Vec4f
	Normal_ms mgl.Vec4f
}

type Mesh struct {
	vertices []MeshVertex
	indices  []MeshIndex
}

type Mode int

const (
	Points Mode = iota
	LineStrip
	LineLoop
	Lines
	TriangleStrip
	TriangleFan
	Triangles
)

type IMesh interface {
	Vertices() interface{}
	Indices() interface{}
	InstanceData() interface{}
	Mode() Mode
}

func (this *Mesh) Vertices() interface{} {
	return this.vertices
}

func (this *Mesh) Indices() interface{} {
	return this.indices
}

func (this *Mesh) InstanceData() interface{} {
	return nil
}

func (this *Mesh) Mode() Mode {
	return Triangles
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

func PortalQuad() (mesh *Mesh) {
	mesh = &Mesh{
		[]MeshVertex{
			MeshVertex{mgl.Vec4f{-1, -1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{1, -1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{1, 1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{-1, 1, 0, 1}, mgl.Vec4f{0, 0, 1, 0}},

			MeshVertex{mgl.Vec4f{-1, -1, 0.5, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{1, -1, 0.5, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{1, 1, 0.5, 1}, mgl.Vec4f{0, 0, 1, 0}},
			MeshVertex{mgl.Vec4f{-1, 1, 0.5, 1}, mgl.Vec4f{0, 0, 1, 0}},
		},
		[]MeshIndex{
			0, 1, 2, 2, 3, 0,
			4, 5, 6, 6, 7, 4,
		},
	}
	return
}

func PortalRect() (mesh *Mesh) {
	mesh = &Mesh{
		[]MeshVertex{
			MeshVertex{mgl.Vec4f{-1, 0, -2, 1}, mgl.Vec4f{0, 1, 0, 0}},
			MeshVertex{mgl.Vec4f{-1, 0, 2, 1}, mgl.Vec4f{0, 1, 0, 0}},
			MeshVertex{mgl.Vec4f{1, 0, 2, 1}, mgl.Vec4f{0, 1, 0, 0}},
			MeshVertex{mgl.Vec4f{1, 0, -2, 1}, mgl.Vec4f{0, 1, 0, 0}},
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

type MyLogStream int

func (mls MyLogStream) Log(msg string) {
	fmt.Println(msg)
}

func LoadMesh(filename string) (mesh *Mesh) {
	scene := ai.ImportFile(filename, 0)
	if scene == nil {
		panic(ai.GetErrorString())
	}
	scene = scene.ApplyPostProcessing(ai.Process_Triangulate)
	meshes := scene.Meshes()
	if len(meshes) != 1 {
		fmt.Println("Cameras", len(scene.Cameras()))
		fmt.Println("Animations", len(scene.Animations()))
		panic(fmt.Sprintf("not a single mesh in %s, found %d meshes", filename, len(meshes)))
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

func LoadMeshManaged(filename string) (mesh *Mesh) {
	mesh = LoadMesh(filename)
	helpers.Manage(mesh, filename)
	return mesh
}

func (this *Mesh) Update(filename string) {
	*this = *LoadMesh(filename)
}

func (m *Mesh) BoundingBox() (min mgl.Vec4f, max mgl.Vec4f) {
	min = mgl.Vec4f{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max = mgl.Vec4f{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}

	for _, v := range m.vertices {
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
