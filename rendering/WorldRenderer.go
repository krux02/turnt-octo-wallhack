package rendering

import "github.com/go-gl/gl"
import "github.com/krux02/mathgl"
import "github.com/krux02/turnt-octo-wallhack/helpers"
import "github.com/krux02/turnt-octo-wallhack/world"
import "math"

// import "fmt"

type Vertex struct {
	Vertex_ms, Normal_ms mathgl.Vec3f
}

type WorldRenderer struct {
	Program       gl.Program
	WorldRenLoc   WorldRenderLocations
	WorldVAO      gl.VertexArray
	Indices       gl.Buffer
	Vertices      gl.Buffer
	WorldNumverts int
}

type WorldRenderLocations struct {
	Vertex_ms, Normal_ms                                    gl.AttribLocation
	Matrix, Model, Time, SeaLevel, Highlight                gl.UniformLocation
	Min_h, Max_h, U_color, U_texture, U_slope, U_screenRect gl.UniformLocation
}

func Vertices(m *world.HeightMap) []Vertex {
	vertices := make([]Vertex, (m.W+1)*(m.H+1))

	i := 0

	for y := 0; y <= m.H; y++ {
		for x := 0; x <= m.W; x++ {
			h := m.Get(x, y)
			pos := mathgl.Vec3f{float32(x), float32(y), h}
			nor := m.Normal(x, y)
			vertices[i] = Vertex{pos, nor}
			i += 1
		}
	}

	return vertices
}

func NewWorldRenderer(heightMap *world.HeightMap) *WorldRenderer {
	vertices := Vertices(heightMap)
	indices := heightMap.Triangulate()
	min_h, max_h := heightMap.Bounds()

	prog := helpers.MakeProgram("World.vs", "World.fs")
	prog.Use()

	vao_A := gl.GenVertexArray()
	vao_A.Bind()

	Loc := WorldRenderLocations{}
	helpers.BindLocations(prog, &Loc)

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	verticesBuffer := gl.GenBuffer()
	verticesBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&Loc, &Vertex{}, true)

	Loc.U_color.Uniform1i(0)
	Loc.U_texture.Uniform1i(1)
	Loc.U_slope.Uniform1i(2)
	Loc.U_screenRect.Uniform1i(3)
	Loc.Min_h.Uniform1f(min_h)
	Loc.Max_h.Uniform1f(max_h)

	return &WorldRenderer{
		prog,
		Loc,
		vao_A,
		indexBuffer,
		verticesBuffer,
		len(indices),
	}
}

func (wr *WorldRenderer) Delete() {
	wr.Program.Delete()
	wr.WorldVAO.Delete()
	wr.Indices.Delete()
	wr.Vertices.Delete()
}

func (wr *WorldRenderer) Render(world *world.HeightMap, Proj mathgl.Mat4f, View mathgl.Mat4f, time float64, highlight int) {
	wr.Program.Use()
	wr.WorldVAO.Bind()

	Loc := wr.WorldRenLoc
	Loc.Time.Uniform1f(float32(time))
	Loc.SeaLevel.Uniform1f(float32(math.Sin(time*0.1)*10 - 5))
	Loc.Highlight.Uniform1f(float32(highlight))

	numverts := wr.WorldNumverts

	
	ProjView := Proj.Mul4(View)

	w := world.W
	h := world.H

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			Model := mathgl.Translate3D(float64(i*w), float64(j*h), 0)
			ProjViewModel := ProjView.Mul4(Model)

			wr.WorldRenLoc.Matrix.UniformMatrix4f(false, (*[16]float32)(&ProjViewModel))
			wr.WorldRenLoc.Model.UniformMatrix4f(false, (*[16]float32)(&Model))

			gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
		}
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)
}
