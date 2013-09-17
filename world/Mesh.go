package world

import (
	// "fmt"
	ai "github.com/krux02/assimp"
	mgl "github.com/krux02/mathgl"
)

type MeshIndex uint16

type MeshVertex struct {
	Position mgl.Vec4f
	Normal mgl.Vec4f
}

type Mesh struct {
	vertices []MeshVertex
	indices []MeshIndex
}


func LoadMesh(filename string) *Mesh {
	scene := ai.ImportFile(filename, 0)
	meshes := scene.Meshes()
	if len(meshes) != 1 {
		panic("not a single mesh")
	}
	mesh := meshes[0]

	meshVertices := make([]MeshVertex, mesh.NumVertices())
	for i,pos := range mesh.Vertices() {
		v := pos.Values()
		 meshVertices[i].Position = mgl.Vec4f([4]float32{v[0],v[1],v[2],1})
	}
	for i,norm := range mesh.Normals() {
		n := norm.Values()
		meshVertices[i].Normal = mgl.Vec4f([4]float32{n[0],n[1],n[2],0})
	}

	meshIndices := make([]MeshIndex, mesh.NumFaces()*3)
	for i,face := range mesh.Faces() {
		indices := face.CopyIndices()
		meshIndices[i*3+0] = MeshIndex(indices[0])
		meshIndices[i*3+1] = MeshIndex(indices[1])
		meshIndices[i*3+2] = MeshIndex(indices[2])
	}

	return &Mesh{ meshVertices, meshIndices }
}
