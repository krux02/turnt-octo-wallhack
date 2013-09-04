package main

import "github.com/go-gl/gl"
import "github.com/krux02/mathgl"
import "unsafe"
import "github.com/Jragonmiris/mathgl/examples/opengl-tutorial/helper"


type Vertex struct {
	position mathgl.Vec3f
	normal   mathgl.Vec3f
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
	vertexPosition_modelspace gl.AttribLocation
	vertexNormal_modelspace   gl.AttribLocation
	matrix                    gl.UniformLocation
	model                     gl.UniformLocation
	time                      gl.UniformLocation
	seaLevel                  gl.UniformLocation
	highlight                 gl.UniformLocation
	// u_color                   gl.UniformLocation
	// u_texture                 gl.UniformLocation
	// u_slope                   gl.UniformLocation
	// u_screenRect              gl.UniformLocation
	// min_h                     gl.UniformLocation
	// max_h                     gl.UniformLocation
}

func NewWorldRenderer(heightMap *HeightMap) *WorldRenderer {
	vertices := heightMap.Vertices()
	indices := heightMap.Triangulate()
	min_h, max_h := heightMap.Bounds()

	prog := helper.MakeProgram("World.vs", "World.fs")
	prog.Use()

	vao_A := gl.GenVertexArray()
	vao_A.Bind()
	vertexPosLoc := prog.GetAttribLocation("vertexPosition_modelspace")
	vertexPosLoc.EnableArray()
	vertexNormLoc := prog.GetAttribLocation("vertexNormal_modelspace")
	vertexNormLoc.EnableArray()

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(int(0))), indices, gl.STATIC_DRAW)

	verticesBuffer := gl.GenBuffer()

	verticesBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(Vertex{})), vertices, gl.STATIC_DRAW)
	vertexPosLoc.AttribPointer(3, gl.FLOAT, false, vertexStride, unsafe.Offsetof(Vertex{}.position))
	vertexNormLoc.AttribPointer(3, gl.FLOAT, false, vertexStride, unsafe.Offsetof(Vertex{}.normal))

	matrixLoc := prog.GetUniformLocation("matrix")
	modelLoc := prog.GetUniformLocation("model")
	timeLoc := prog.GetUniformLocation("time")
	seaLevelLoc := prog.GetUniformLocation("seaLevel")
	highlightLoc := prog.GetUniformLocation("highlight")

	prog.GetUniformLocation("u_color").Uniform1i(0)
	prog.GetUniformLocation("u_texture").Uniform1i(1)
	prog.GetUniformLocation("u_slope").Uniform1i(2)
	prog.GetUniformLocation("u_screenRect").Uniform1i(3)
	prog.GetUniformLocation("min_h").Uniform1f(min_h)
	prog.GetUniformLocation("max_h").Uniform1f(max_h)

	wrl := WorldRenderLocations{
		vertexPosLoc,
		vertexNormLoc,
		matrixLoc,
		modelLoc,
		timeLoc,
		seaLevelLoc,
		highlightLoc,
	}

	return &WorldRenderer{
		prog,
		wrl,
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

			wr.WorldRenLoc.matrix.UniformMatrix4f(false, (*[16]float32)(&finalMat))
			wr.WorldRenLoc.model.UniformMatrix4f(false, (*[16]float32)(&modelMat))

			gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
		}
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)
	gl.DepthMask(false)
}

		
