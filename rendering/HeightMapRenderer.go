package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

// import "fmt"

type HeightMapVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

type HeightMapRenderer struct {
	Program gl.Program
	RenLoc  HeightMapRenderLocations
	Data    HeightMapRenderData
}

type HeightMapRenderData struct {
	VAO      gl.VertexArray
	Indices  gl.Buffer
	Vertices gl.Buffer
	Numverts int
}

type HeightMapRenderLocations struct {
	Vertex_ms, Normal_ms                     gl.AttribLocation
	Matrix, Model                            gl.UniformLocation
	HeightMap, Color, Texture, Slope         gl.UniformLocation
	ClippingPlane_ws, LowerBound, UpperBound gl.UniformLocation
}

func Vertices(m *gamestate.HeightMap) []HeightMapVertex {
	vertices := make([]HeightMapVertex, (m.W+1)*(m.H+1))

	i := 0

	for y := 0; y <= m.H; y++ {
		for x := 0; x <= m.W; x++ {
			h := m.Get(x, y)
			pos := mgl.Vec3f{float32(x), float32(y), h}
			nor := m.Normal(x, y)
			vertices[i] = HeightMapVertex{pos, nor}
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

	helpers.BindLocations("height map", this.Program, &this.RenLoc)

	this.Data.VAO = gl.GenVertexArray()
	this.Data.VAO.Bind()

	this.Data.Indices = gl.GenBuffer()
	this.Data.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	this.Data.Vertices = gl.GenBuffer()
	this.Data.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&this.RenLoc, &HeightMapVertex{})

	this.Data.Numverts = len(indices)

	this.RenLoc.HeightMap.Uniform1i(4)
	this.RenLoc.Color.Uniform1i(3)
	this.RenLoc.Slope.Uniform1i(2)
	this.RenLoc.Texture.Uniform1i(1)
	this.RenLoc.LowerBound.Uniform3f(0, 0, min_h)
	this.RenLoc.UpperBound.Uniform3f(float32(heightMap.W), float32(heightMap.H), max_h)

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
	Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	numverts := wr.Data.Numverts

	ProjView := Proj.Mul4(View)
	ProjViewModel := ProjView.Mul4(Model)

	wr.RenLoc.Matrix.UniformMatrix4f(false, glMat(&ProjViewModel))
	wr.RenLoc.Model.UniformMatrix4f(false, glMat(&Model))

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
}
