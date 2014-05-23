package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"unsafe"
)

type ScreenQuadRenderer struct {
	Prog      gl.Program
	Vao       gl.VertexArray
	Buffer    gl.Buffer
	Locations RenderLocations
}

func NewScreenQuadRenderer() (this *ScreenQuadRenderer) {
	this = new(ScreenQuadRenderer)
	this.Prog = helpers.MakeProgram("ScreenQuad.vs", "ScreenQuad.fs")
	this.Prog.Use()

	this.Vao = gl.GenVertexArray()
	this.Vao.Bind()
	helpers.BindLocations("screen quad", this.Prog, &this.Locations)

	this.Locations.Vertex_ndc.EnableArray()
	this.Buffer = gl.GenBuffer()
	this.Buffer.Bind(gl.ARRAY_BUFFER)

	arr := []mgl.Vec4f{
		mgl.Vec4f{-1, -1, 0, 1},
		mgl.Vec4f{3, -1, 0, 1},
		mgl.Vec4f{-1, 3, 0, 1},
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(arr)*int(unsafe.Sizeof(mgl.Vec4f{})), arr, gl.STATIC_DRAW)
	this.Locations.Vertex_ndc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))

	return
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
