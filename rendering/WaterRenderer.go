package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"math"
)

type WaterVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

type WaterRenderer struct {
	Program gl.Program
	RenLoc  WorldRenderLocations
	Data    WorldRenderData
}

type WaterRenderData struct {
	VAO      gl.VertexArray
	Indices  gl.Buffer
	Vertices gl.Buffer
	Numverts int
}

type WaterRenderLocations struct {
	Vertex_ms, Normal_ms                                    gl.AttribLocation
	Matrix, Model, Time, SeaLevel, Highlight                gl.UniformLocation
	Min_h, Max_h, U_color, U_texture, U_slope, U_screenRect gl.UniformLocation
	U_clippingPlane                                         gl.UniformLocation
}

func WaterVertices(W, H int) []WaterVertex {
	vertices := make([]WaterVertex, (W+1)*(H+1))

	i := 0
	for y := 0; y <= H; y++ {
		for x := 0; x <= W; x++ {
			pos := mgl.Vec3f{float32(x), float32(y), 0}
			nor := mgl.Vec3f{0, 0, 1}
			vertices[i] = WaterVertex{pos, nor}
			i += 1
		}
	}

	return vertices
}

func NewWaterRenderer(heightMap *gamestate.HeightMap) (this *WaterRenderer) {
	vertices := WaterVertices(heightMap.W, heightMap.H)
	indices := heightMap.Triangulate()
	min_h, max_h := heightMap.Bounds()

	this = new(WaterRenderer)

	this.Program = helpers.MakeProgram("Water.vs", "Water.fs")
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

	helpers.SetAttribPointers(&this.RenLoc, &WorldVertex{}, true)

	this.Data.Numverts = len(indices)

	this.RenLoc.U_color.Uniform1i(0)
	this.RenLoc.U_texture.Uniform1i(1)
	this.RenLoc.U_slope.Uniform1i(2)
	this.RenLoc.U_screenRect.Uniform1i(3)
	this.RenLoc.Min_h.Uniform1f(min_h)
	this.RenLoc.Max_h.Uniform1f(max_h)

	return
}

func (wr *WaterRenderer) Delete() {
	wr.Program.Delete()
	wr.Data.VAO.Delete()
	wr.Data.Indices.Delete()
	wr.Data.Vertices.Delete()
}

func (wr *WaterRenderer) Render(Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, time float64, clippingPlane mgl.Vec4f) {
	wr.Program.Use()
	wr.Data.VAO.Bind()

	Loc := wr.RenLoc
	Loc.Time.Uniform1f(float32(time))
	Loc.SeaLevel.Uniform1f(float32(math.Sin(time*0.1)*10 - 5))
	Loc.Highlight.Uniform1f(1)
	Loc.U_clippingPlane.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	numverts := wr.Data.Numverts
	ProjView := Proj.Mul4(View)
	ProjViewModel := ProjView.Mul4(Model)

	wr.RenLoc.Matrix.UniformMatrix4f(false, (*[16]float32)(&ProjViewModel))
	wr.RenLoc.Model.UniformMatrix4f(false, (*[16]float32)(&Model))

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
}
