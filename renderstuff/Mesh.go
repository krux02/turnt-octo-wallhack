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

type Mesh struct {
	Vertices     interface{}
	Indices      interface{}
	InstanceData interface{}
	Mode         Mode
}
