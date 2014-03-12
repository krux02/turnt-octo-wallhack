package debug

import (
	mgl "github.com/Jragonmiris/mathgl"
	"unsafe"
)

const LINE_VERTEX_STRIDE = int(unsafe.Sizeof(&LineVertex{}))

type LineVertex struct {
	Vertex_ws mgl.Vec4f
	Color     mgl.Vec4f
}

var color mgl.Vec4f
var data = make([]LineVertex, 0, 64)

func Line(pos1, pos2 mgl.Vec4f) {
	v1 := LineVertex{pos1, color}
	v2 := LineVertex{pos2, color}
	data = append(data, v1, v2)
}

func Color(_color mgl.Vec4f) {
	color = _color
}

func ReadAndReset() []LineVertex {
	d := data
	data = data[0:0]
	return d
}
