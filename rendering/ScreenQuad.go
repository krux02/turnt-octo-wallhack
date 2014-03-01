package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"unsafe"
)

type ScreenQuadLocations struct {
	Image        gl.UniformLocation
	Position_ndc gl.AttribLocation
}

type ScreenQuadRenderer struct {
	Prog      gl.Program
	Vao       gl.VertexArray
	Buffer    gl.Buffer
	Locations ScreenQuadLocations
}

func NewScreenQuadRenderer() *ScreenQuadRenderer {

	prog := helpers.MakeProgram("ScreenQuad.vs", "ScreenQuad.fs")
	prog.Use()

	vao := gl.GenVertexArray()
	vao.Bind()

	locations := ScreenQuadLocations{}
	helpers.BindLocations("screen quad", prog, &locations)

	locations.Position_ndc.EnableArray()
	a_positionBuffer := gl.GenBuffer()
	a_positionBuffer.Bind(gl.ARRAY_BUFFER)

	arr := []mgl.Vec4f{
		mgl.Vec4f{-1, -1, 0, 1},
		mgl.Vec4f{3, -1, 0, 1},
		mgl.Vec4f{-1, 3, 0, 1},
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(arr)*int(unsafe.Sizeof(mgl.Vec4f{})), arr, gl.STATIC_DRAW)
	locations.Position_ndc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))

	return &ScreenQuadRenderer{prog, vao, a_positionBuffer, locations}
}

func (this *ScreenQuadRenderer) Delete() {
	this.Prog.Delete()
	this.Vao.Delete()
	this.Buffer.Delete()
	*this = ScreenQuadRenderer{}
}

func (this *ScreenQuadRenderer) Render(textureUnit int) {
	this.Prog.Use()
	this.Vao.Bind()
	this.Locations.Image.Uniform1i(textureUnit)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
