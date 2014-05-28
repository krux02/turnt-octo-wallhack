package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/debug"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type LineRenderer struct {
	Prog   gl.Program
	RenLoc RenderLocations
	vao    gl.VertexArray
	buffer gl.Buffer
}

func NewLineRenderer() *LineRenderer {
	renderer := LineRenderer{}
	renderer.Prog = helpers.MakeProgram("Line.vs", "Line.fs")
	renderer.Prog.Use()
	renderer.vao = gl.GenVertexArray()
	renderer.vao.Bind()
	helpers.BindLocations("line", renderer.Prog, &renderer.RenLoc)
	renderer.buffer = gl.GenBuffer()
	renderer.buffer.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&renderer.RenLoc, &debug.LineVertex{}, false)

	fmt.Println("Line render location ", renderer.RenLoc)
	return &renderer
}

func (this *LineRenderer) Render(Proj, View mgl.Mat4f) {

	data := debug.Read()
	if len(data) > 0 {
		this.Prog.Use()
		this.vao.Bind()
		this.buffer.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(data), data, gl.STREAM_DRAW)
		this.RenLoc.Proj.UniformMatrix4f(false, glMat4(&Proj))
		this.RenLoc.View.UniformMatrix4f(false, glMat4(&View))
		gl.DrawArrays(gl.LINES, 0, len(data))
	}
}

func (this *LineRenderer) Delete() {
	this.Prog.Delete()
	this.vao.Delete()
	this.buffer.Delete()
	*this = LineRenderer{}
}
