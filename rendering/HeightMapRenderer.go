package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

// import "fmt"

type WorldVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

type HeightMapRenderer struct {
	Program gl.Program
	RenLoc  WorldRenderLocations
	Data    WorldRenderData
}

type WorldRenderData struct {
	VAO      gl.VertexArray
	Indices  gl.Buffer
	Vertices gl.Buffer
	Numverts int
}

type WorldRenderLocations struct {
	Vertex_ms, Normal_ms                     gl.AttribLocation
	Matrix, Model                            gl.UniformLocation
	U_HeightMap, U_color, U_texture, U_slope gl.UniformLocation
	U_clippingPlane, Min_h, Max_h            gl.UniformLocation
}

func Vertices(m *gamestate.HeightMap) []WorldVertex {
	vertices := make([]WorldVertex, (m.W+1)*(m.H+1))

	i := 0

	for y := 0; y <= m.H; y++ {
		for x := 0; x <= m.W; x++ {
			h := m.Get(x, y)
			pos := mgl.Vec3f{float32(x), float32(y), h}
			nor := m.Normal(x, y)
			vertices[i] = WorldVertex{pos, nor}
			i += 1
		}
	}

	return vertices
}

func NewHeightMapRenderer(heightMap *gamestate.HeightMap) (this *HeightMapRenderer) {
	vertices := Vertices(heightMap)
	indices := TriangulationIndices(heightMap.W, heightMap.H)
	min_h, max_h := heightMap.Bounds()

	this = new(HeightMapRenderer)

	this.Program = helpers.MakeProgram("HeightMap.vs", "HeightMap.fs")
	this.Program.Use()

	helpers.BindLocations(this.Program, &this.RenLoc)

	this.Data.VAO = gl.GenVertexArray()
	this.Data.VAO.Bind()

	this.Data.Indices = gl.GenBuffer()
	this.Data.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	this.Data.Vertices = gl.GenBuffer()
	this.Data.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&this.RenLoc, &WorldVertex{})

	this.Data.Numverts = len(indices)

	this.RenLoc.U_color.Uniform1i(3)
	this.RenLoc.U_texture.Uniform1i(1)
	this.RenLoc.U_slope.Uniform1i(2)
	this.RenLoc.Min_h.Uniform1f(min_h)
	this.RenLoc.Max_h.Uniform1f(max_h)

	return
}

func (wr *HeightMapRenderer) Delete() {
	wr.Program.Delete()
	wr.Data.VAO.Delete()
	wr.Data.Indices.Delete()
	wr.Data.Vertices.Delete()
}

func (wr *HeightMapRenderer) Render(Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, clippingPlane mgl.Vec4f) {
	wr.Program.Use()
	wr.Data.VAO.Bind()

	Loc := wr.RenLoc
	Loc.U_clippingPlane.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	numverts := wr.Data.Numverts

	ProjView := Proj.Mul4(View)
	ProjViewModel := ProjView.Mul4(Model)

	wr.RenLoc.Matrix.UniformMatrix4f(false, (*[16]float32)(&ProjViewModel))
	wr.RenLoc.Model.UniformMatrix4f(false, (*[16]float32)(&Model))

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
}
