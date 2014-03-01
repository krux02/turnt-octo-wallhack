package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type WaterVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

type WaterRenderer struct {
	Program      gl.Program
	RenLoc       WaterRenderLocations
	DebugProgram gl.Program
	DebugRenLoc  WaterRenderLocations
	Data         WaterRenderData
}

type WaterRenderData struct {
	VAO, DebugVao gl.VertexArray
	Indices       gl.Buffer
	Vertices      gl.Buffer
	Numverts      int
}

type WaterRenderLocations struct {
	Vertex_ms, Normal_ms              gl.AttribLocation
	HeightMap, LowerBound, UpperBound gl.UniformLocation
	ClippingPlane_ws, CameraPos_ws    gl.UniformLocation
	GroundTexture, Skybox             gl.UniformLocation
	Time, Model, View, Proj           gl.UniformLocation
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
	indices := TriangulationIndices(heightMap.W, heightMap.H)
	min_h, max_h := heightMap.Bounds()
	W, H := float32(heightMap.W), float32(heightMap.H)

	this = new(WaterRenderer)

	this.Program = helpers.MakeProgram("Water.vs", "Water.fs")
	this.Program.Use()
	helpers.BindLocations("water", this.Program, &this.RenLoc)

	this.Data.VAO = gl.GenVertexArray()
	this.Data.VAO.Bind()

	this.Data.Indices = gl.GenBuffer()
	this.Data.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	this.Data.Vertices = gl.GenBuffer()
	this.Data.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)

	helpers.SetAttribPointers(&this.RenLoc, &WaterVertex{})
	this.Data.Numverts = len(indices)

	this.RenLoc.HeightMap.Uniform1i(4)
	this.RenLoc.LowerBound.Uniform3f(0, 0, min_h)
	this.RenLoc.UpperBound.Uniform3f(W, H, max_h)
	this.RenLoc.GroundTexture.Uniform1i(1)
	this.RenLoc.Skybox.Uniform1i(7)

	this.DebugProgram = helpers.MakeProgram3("Water.vs", "Normal.gs", "Line.fs")
	this.DebugProgram.Use()

	this.Data.DebugVao = gl.GenVertexArray()
	this.Data.DebugVao.Bind()
	this.Data.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	this.Data.Vertices.Bind(gl.ARRAY_BUFFER)

	helpers.BindLocations("water debug", this.DebugProgram, &this.DebugRenLoc)

	this.DebugRenLoc.HeightMap.Uniform1i(4)
	this.DebugRenLoc.LowerBound.Uniform3f(0, 0, min_h)
	this.DebugRenLoc.UpperBound.Uniform3f(W, H, max_h)
	this.DebugRenLoc.GroundTexture.Uniform1i(1)

	helpers.SetAttribPointers(&this.RenLoc, &WaterVertex{})

	return
}

func (wr *WaterRenderer) Delete() {
	wr.Program.Delete()
	wr.DebugProgram.Delete()
	wr.Data.VAO.Delete()
	wr.Data.Indices.Delete()
	wr.Data.Vertices.Delete()
}

func (wr *WaterRenderer) Render(Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, time float64, clippingPlane mgl.Vec4f, normals bool) {
	wr.Program.Use()
	wr.Data.VAO.Bind()

	numverts := wr.Data.Numverts

	Loc := wr.RenLoc
	Loc.Time.Uniform1f(float32(time))
	Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])
	Loc.Proj.UniformMatrix4f(false, glMat(&Proj))
	Loc.Model.UniformMatrix4f(false, glMat(&Model))
	Loc.View.UniformMatrix4f(false, glMat(&View))
	v := View.Inv().Mul4x1(mgl.Vec4f{0, 0, 0, 1})
	Loc.CameraPos_ws.Uniform4f(v[0], v[1], v[2], v[3])

	gl.Disable(gl.CULL_FACE)

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))

	// debug rendering
	if normals {
		wr.DebugProgram.Use()
		wr.Data.DebugVao.Bind()

		Loc = wr.DebugRenLoc
		Loc.Time.Uniform1f(float32(time))
		Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])
		Loc.Proj.UniformMatrix4f(false, glMat(&Proj))
		Loc.Model.UniformMatrix4f(false, glMat(&Model))
		Loc.View.UniformMatrix4f(false, glMat(&View))
		Loc.CameraPos_ws.Uniform4f(v[0], v[1], v[2], v[3])

		gl.Disable(gl.CULL_FACE)

		gl.DrawElements(gl.POINTS, numverts, gl.UNSIGNED_INT, uintptr(0))
	}
}
