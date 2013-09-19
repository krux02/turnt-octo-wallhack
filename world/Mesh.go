package world

import (
	"fmt"
	ai "github.com/krux02/assimp"
	mgl "github.com/krux02/mathgl"
)

type MeshIndex uint16

type MeshVertex struct {
	Position mgl.Vec4f
	Normal   mgl.Vec4f
}

type Mesh struct {
	Vertices []MeshVertex
	Indices  []MeshIndex
}

func LoadMesh(filename string) (mesh *Mesh) {
	scene := ai.ImportFile(filename, 0)
	meshes := scene.Meshes()
	if len(meshes) != 1 {
		panic("not a single mesh")
	}
	aimesh := meshes[0]

	meshVertices := make([]MeshVertex, aimesh.NumVertices())
	for i, pos := range aimesh.Vertices() {
		v := pos.Values()
		meshVertices[i].Position = mgl.Vec4f([4]float32{v[0], v[1], v[2], 1})
	}
	for i, norm := range aimesh.Normals() {
		n := norm.Values()
		meshVertices[i].Normal = mgl.Vec4f([4]float32{n[0], n[1], n[2], 0})
	}

	meshIndices := make([]MeshIndex, aimesh.NumFaces()*3)
	for i, face := range aimesh.Faces() {
		indices := face.CopyIndices()
		meshIndices[i*3+0] = MeshIndex(indices[0])
		meshIndices[i*3+1] = MeshIndex(indices[1])
		meshIndices[i*3+2] = MeshIndex(indices[2])
	}

	fmt.Println("loaded mesh with", aimesh.NumVertices(), "Vertices")ViewModel

	return &Mesh{meshVertices, meshIndices}
}
