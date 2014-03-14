package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/debug"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type LineRenderLocatins struct {
	View, Proj       gl.UniformLocation
	Vertex_ws, Color gl.AttribLocation
}

type LineRenderer struct {
	Prog   gl.Program
	Loc    LineRenderLocatins
	vao    gl.VertexArray
	buffer gl.Buffer
}

func NewLineRenderer() *LineRenderer {
	renderer := LineRenderer{}
	renderer.Prog = helpers.MakeProgram("Line.vs", "Line.fs")
	renderer.Prog.Use()
	renderer.vao = gl.GenVertexArray()
	renderer.vao.Bind()
	helpers.BindLocations("line", renderer.Prog, &renderer.Loc)
	renderer.buffer = gl.GenBuffer()
	helpers.SetAttribPointers(&renderer.Loc, &debug.LineVertex{})
	return &renderer
}

func (this *LineRenderer) Render(Proj, View mgl.Mat4f) {
	this.Prog.Use()
	this.vao.Bind()
	data := debug.Read()
	if len(data) > 0 {
		fmt.Println(data)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(data), data, gl.STREAM_DRAW)
		this.Loc.Proj.UniformMatrix4f(false, glMat(&Proj))
		this.Loc.View.UniformMatrix4f(false, glMat(&View))
		gl.DrawArrays(gl.LINES, 0, len(data))
	}
}

func (this *LineRenderer) Delete() {
	this.Prog.Delete()
	this.vao.Delete()
	this.buffer.Delete()
	*this = LineRenderer{}
}
