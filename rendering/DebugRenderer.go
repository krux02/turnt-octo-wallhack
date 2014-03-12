package rendering

import (
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
	helpers.BindLocations("line", renderer.Prog, &renderer.Loc)
	renderer.vao = gl.GenVertexArray()
	renderer.buffer = gl.GenBuffer()
	helpers.SetAttribPointers(&renderer.Loc, &debug.LineVertex{})
	return &renderer
}

func (this *LineRenderer) Render() {
	this.Prog.Use()
	this.vao.Bind()
	data := debug.ReadAndReset()
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(data), data, gl.STREAM_DRAW)
	gl.DrawArrays(gl.LINES, 0, len(data))
}
