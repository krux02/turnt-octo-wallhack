package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type HeightMapRenderer struct {
	Program gl.Program
	RenLoc  RenderLocations
	Data    RenderData
}

func NewHeightMapRenderer(heightMap *gamestate.HeightMap) (this *HeightMapRenderer) {

	this = new(HeightMapRenderer)

	this.Program = helpers.MakeProgram("HeightMap.vs", "HeightMap.fs")
	this.Program.Use()

	helpers.BindLocations("height map", this.Program, &this.RenLoc)

	this.Data = CreateHeightMapRenderData(heightMap)

	helpers.SetAttribPointers(&this.RenLoc, &HeightMapVertex{})

	this.RenLoc.HeightMap.Uniform1i(4)
	this.RenLoc.ColorBand.Uniform1i(3)
	this.RenLoc.Slope.Uniform1i(2)
	this.RenLoc.Texture.Uniform1i(1)

	return
}

func (wr *HeightMapRenderer) Delete() {
	wr.Program.Delete()
	wr.Data.VAO.Delete()
	wr.Data.Indices.Delete()
	wr.Data.Vertices.Delete()
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

func CreateHeightMapRenderData(heightMap *gamestate.HeightMap) (data RenderData) {
	vertices := Vertices(heightMap)
	indices := TriangulationIndices(heightMap.W, heightMap.H)

	data.VAO = gl.GenVertexArray()
	data.VAO.Bind()

	data.Indices = gl.GenBuffer()
	data.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)

	data.Vertices = gl.GenBuffer()
	data.Vertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)

	data.Numverts = len(indices)

	return
}

func (this *HeightMapRenderer) Render(heightMap *gamestate.HeightMap, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, clippingPlane mgl.Vec4f) {
	this.Program.Use()
	this.Data.VAO.Bind()

	Loc := this.RenLoc
	Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	numverts := this.Data.Numverts

	if heightMap.HasChanges {
		min_h, max_h := heightMap.Bounds()
		this.RenLoc.LowerBound.Uniform3f(0, 0, min_h)
		this.RenLoc.UpperBound.Uniform3f(float32(heightMap.W), float32(heightMap.H), max_h)

		gl.ActiveTexture(gl.TEXTURE4)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R16, heightMap.W, heightMap.H, 0, gl.RED, gl.FLOAT, heightMap.TexturePixels())
		gl.ActiveTexture(gl.TEXTURE0)

		heightMap.HasChanges = false
	}

	this.RenLoc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	this.RenLoc.View.UniformMatrix4f(false, glMat4(&View))
	this.RenLoc.Model.UniformMatrix4f(false, glMat4(&Model))

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
}
