package renderstuff

import (
// "fmt"
)

type Mode int

const (
	Points Mode = iota
	LineStrip
	LineLoop
	Lines
	TriangleStrip
	TriangleFan
	Triangles
	Undef
)

type IMesh interface {
	Vertices() interface{}
	Indices() interface{}
	InstanceData() interface{}
	Mode() Mode
}

// embed this for default values
type AbstractMesh struct {
	vertics             interface{}
	verticesChanged     bool
	indicesChanged      bool
	instanceDataChanged bool
}

func (this *AbstractMesh) Indices() interface{} {
	return nil
}

func (this *AbstractMesh) InstanceData() interface{} {
	return nil
}

func (this *AbstractMesh) VerticesChanged() bool {
	return this.verticesChanged
}

func (this *AbstractMesh) SetVerticesChanged(changed bool) {
	this.verticesChanged = changed
}

func (this *AbstractMesh) IndicesChanged() bool {
	return this.indicesChanged
}

func (this *AbstractMesh) SetIndicesChanged(changed bool) {
	this.indicesChanged = changed
}

func (this *AbstractMesh) InstanceDataChanged() bool {
	return this.instanceDataChanged
}

func (this *AbstractMesh) SetInstanceDataChanged(changed bool) {
	this.instanceDataChanged = changed
}
