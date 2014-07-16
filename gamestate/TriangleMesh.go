package gamestate

import (
	"fmt"
	ai "github.com/krux02/assimp"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
	"math"
)

type MeshIndex uint16

type MeshVertex struct {
	Vertex_ms mgl.Vec4
	Normal_ms mgl.Vec4
}

type TriangleMesh renderstuff.Mesh

func QuadMesh() (mesh *renderstuff.Mesh) {
	mesh = new(renderstuff.Mesh)
	mesh.Vertices = []MeshVertex{
		MeshVertex{mgl.Vec4{-1, -1, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{1, -1, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{1, 1, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{-1, 1, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
	}
	mesh.Indices = []MeshIndex{0, 1, 2, 2, 3, 0}
	mesh.Mode = renderstuff.Triangles
	return mesh
}

const (
	a = 1
	b = 0.8
	c = 0.0
)

func PortalQuad() (mesh *renderstuff.Mesh) {
	mesh = new(renderstuff.Mesh)
	mesh.Vertices = []MeshVertex{
		MeshVertex{mgl.Vec4{-a, -a, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{a, -a, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{a, a, 0, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{-a, a, 0, 1}, mgl.Vec4{0, 0, 1, 0}},

		MeshVertex{mgl.Vec4{-b, -b, c, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{b, -b, c, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{b, b, c, 1}, mgl.Vec4{0, 0, 1, 0}},
		MeshVertex{mgl.Vec4{-b, b, c, 1}, mgl.Vec4{0, 0, 1, 0}},
	}
	mesh.Indices = []MeshIndex{
		0, 4, 3, 3, 7, 2, 2, 6, 1, 1, 5, 0,
		7, 3, 4, 6, 2, 7, 5, 1, 6, 4, 0, 5,
		4, 5, 7, 6, 7, 5,
	}
	mesh.Mode = renderstuff.Triangles
	return
}

func PortalRect() (mesh *renderstuff.Mesh) {
	mesh = new(renderstuff.Mesh)
	mesh.Vertices = []MeshVertex{
		MeshVertex{mgl.Vec4{-1, 0, -2, 1}, mgl.Vec4{0, 1, 0, 0}},
		MeshVertex{mgl.Vec4{-1, 0, 2, 1}, mgl.Vec4{0, 1, 0, 0}},
		MeshVertex{mgl.Vec4{1, 0, 2, 1}, mgl.Vec4{0, 1, 0, 0}},
		MeshVertex{mgl.Vec4{1, 0, -2, 1}, mgl.Vec4{0, 1, 0, 0}},
	}
	mesh.Indices = []MeshIndex{0, 1, 2, 2, 3, 0}
	mesh.Mode = renderstuff.Triangles
	return mesh
}

func Min(v1, v2 mgl.Vec4) (min mgl.Vec4) {
	for i := 0; i < 4; i++ {
		if v1[i] < v2[i] {
			min[i] = v1[i]
		} else {
			min[i] = v2[i]
		}
	}
	return
}

func Max(v1, v2 mgl.Vec4) (min mgl.Vec4) {
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

func LoadMesh(filename string) (mesh *renderstuff.Mesh) {
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

	mesh = new(renderstuff.Mesh)

	meshvertices := make([]MeshVertex, aimesh.NumVertices())
	for i, pos := range aimesh.Vertices() {
		v := pos.Values()
		meshvertices[i].Vertex_ms = mgl.Vec4([4]float32{v[0], v[1], v[2], 1})
	}
	for i, norm := range aimesh.Normals() {
		n := norm.Values()
		meshvertices[i].Normal_ms = mgl.Vec4([4]float32{n[0], n[1], n[2], 0})
	}
	mesh.Vertices = meshvertices

	meshindices := make([]MeshIndex, aimesh.NumFaces()*3)
	for i, face := range aimesh.Faces() {
		indices := face.CopyIndices()
		meshindices[i*3+0] = MeshIndex(indices[0])
		meshindices[i*3+1] = MeshIndex(indices[1])
		meshindices[i*3+2] = MeshIndex(indices[2])
	}
	mesh.Indices = meshindices

	mesh.Mode = renderstuff.Triangles
	fmt.Println("loaded mesh with", aimesh.NumVertices(), "Vertices")
	return mesh
}

func LoadMeshManaged(filename string) (mesh *renderstuff.Mesh) {
	mesh = LoadMesh(filename)
	helpers.Manage((*TriangleMesh)(mesh), filename)
	return mesh
}

func (this *TriangleMesh) Update(filename string) {
	*this = TriangleMesh(*LoadMesh(filename))
}

func (m *TriangleMesh) BoundingBox() (min mgl.Vec4, max mgl.Vec4) {
	min = mgl.Vec4{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max = mgl.Vec4{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}

	vertices := m.Vertices.([]MeshVertex)
	for _, v := range vertices {
		min = Min(min, v.Vertex_ms)
		max = Max(max, v.Vertex_ms)
	}
	return
}

// returns the 8 vertices of the box that is defined by two if it's vertices
func (m *TriangleMesh) MakeBoxVertices() (verts [8]mgl.Vec4) {
	bottomLeft, topRight := m.BoundingBox()

	var i int
	bounds := [2]mgl.Vec4{bottomLeft, topRight}
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
