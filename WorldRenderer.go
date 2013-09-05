package main

import "github.com/go-gl/gl"
import "github.com/krux02/mathgl"
import "unsafe"

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

func NewWorldRenderer(heightMap *HeightMap) *WorldRenderer {
	vertices := heightMap.Vertices()
	indices := heightMap.Triangulate()
	min_h, max_h := heightMap.Bounds()

	prog := MakeProgram("World.vs", "World.fs")
	prog.Use()

	vao_A := gl.GenVertexArray()
	vao_A.Bind()

	Loc := WorldRenderLocations{}
	BindLocations(prog, &Loc)

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(int(0))), indices, gl.STATIC_DRAW)

	verticesBuffer := gl.GenBuffer()
	verticesBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(Vertex{})), vertices, gl.STATIC_DRAW)

	SetAttribPointers(&Loc, &Vertex{})

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

func (wr *WorldRenderer) Render(gamestate *GameState) {
	wr.WorldVAO.Bind()
	numverts := wr.WorldNumverts

	view := gamestate.Camera.View()
	projView := gamestate.Proj.Mul4(view)

	w := gamestate.HeightMap.W
	h := gamestate.HeightMap.H

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			modelMat := mathgl.Translate3D(float64(i*w), float64(j*h), 0)
			finalMat := projView.Mul4(modelMat)

			wr.WorldRenLoc.Matrix.UniformMatrix4f(false, (*[16]float32)(&finalMat))
			wr.WorldRenLoc.Model.UniformMatrix4f(false, (*[16]float32)(&modelMat))

			gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
		}
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)
	gl.DepthMask(false)
}
