package rendering

import (
	"github.com/go-gl/gl"
	"github.com/krux02/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"unsafe"
)

var progB gl.Program
var vao_B gl.VertexArray

func InitScreenQuad() {
	progB = helpers.MakeProgram("screenQuad.vs", "ScreenQuad.fs")

	progB.Use()

	vao_B = gl.GenVertexArray()

	vao_B.Bind()

	a_positionLoc := progB.GetAttribLocation("a_position")
	a_positionLoc.EnableArray()
	a_positionBuffer := gl.GenBuffer()
	a_positionBuffer.Bind(gl.ARRAY_BUFFER)

	arr := []mathgl.Vec4f{
		mathgl.Vec4f{-1, -1, 0, 1},
		mathgl.Vec4f{3, -1, 0, 1},
		mathgl.Vec4f{-1, 3, 0, 1},
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(arr)*int(unsafe.Sizeof(mathgl.Vec4f{})), arr, gl.STATIC_DRAW)
	a_positionLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))

	progB.GetUniformLocation("u_screenRect").Uniform1i(3)
}

func RenderScreenQuad() {
	progB.Use()
	vao_B.Bind()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func FreeScreenQuad() {
	progB.Delete()
	vao_B.Delete()
}
