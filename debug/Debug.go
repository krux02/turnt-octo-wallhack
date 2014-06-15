package debug

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"unsafe"
)

const LINE_VERTEX_STRIDE = int(unsafe.Sizeof(&LineVertex{}))

type LineVertex struct {
	Vertex_ws mgl.Vec4
	Color     mgl.Vec4
}

var color mgl.Vec4
var data = make([]LineVertex, 0, 64)

func Line(pos1, pos2 mgl.Vec4) {
	v1 := LineVertex{pos1, color}
	v2 := LineVertex{pos2, color}
	data = append(data, v1, v2)
}

func Color(_color mgl.Vec4) {
	color = _color
}

func Read() []LineVertex {
	return data
}

func Reset() {
	data = data[0:0]
}
