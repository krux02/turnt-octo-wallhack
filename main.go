package main

import (
	"fmt"
	"github.com/Jragonmiris/mathgl/examples/opengl-tutorial/helper"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/krux02/mathgl"
	"github.com/krux02/tw"
	"math"
	"os"
	"unsafe"
)

type Vertex struct {
	position mathgl.Vec3f
	normal   mathgl.Vec3f
}

type GameState struct {
	Camera         *Camera
	HeightMap      *HeightMap
	ParticleSystem *ParticleSystem
	Player         Player
}

const vertexStride = int(unsafe.Sizeof(Vertex{}))

var mat1 = mathgl.Perspective(90, 4.0/3.0, 0.01, 1000)

func Resize(width int, height int) {
	gl.Viewport(0, 0, width, height)
	mat1 = mathgl.Perspective(90, float64(width)/float64(height), 0.1, 1000)

	tw.WindowSize(width, height)
	// RandomNoiseRectangle(width, height)
}

func main() {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.FsaaSamples, 4)
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.OpenWindowHint(glfw.OpenGLDebugContext, gl.TRUE)

	if err := glfw.OpenWindow(1024, 768, 0, 0, 0, 0, 32, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	gl.Init()
	gl.GetError() // Ignore error

	glfw.SetWindowTitle("gogogo")
	glfw.Enable(glfw.StickyKeys)

	speed := 0.0

	tw.Init(tw.OPENGL_CORE, nil)
	defer tw.Terminate()
	bar := tw.NewBar("TweakBar")
	tw.Define(" GLOBAL help='This example shows how to integrate AntTweakBar with GLFW and OpenGL.' ")
	bar.AddVarRW("speed", tw.TYPE_DOUBLE, unsafe.Pointer(&speed), " label='Rot speed' min=0 max=2 step=0.01 keyIncr=s keyDecr=S help='Rotation speed (turns/second)' ")

	xxx := 0

	bar.AddButton("but1", func() { fmt.Printf("but1 %d\n", xxx); xxx += 1 }, "")
	bar.AddButton("but2", func() { fmt.Printf("but2 %d\n", xxx); xxx += 1 }, "")

	//C.TwWindowSize(1024, 768);

	glfw.SetWindowSizeCallback(Resize)

	initDebugContext()

	const w = 128
	const h = 128

	heights := NewHeightMap(w, h)
	heights.DiamondSquare(w)
	vertices := heights.Vertices()
	indices := heights.Triangulate()
	min_h, max_h := heights.Bounds()

	gl.ClearColor(0., 0., 0.4, 0.)

	prog := helper.MakeProgram("World.vs", "World.fs")
	defer prog.Delete()
	prog.Use()

	vao_A := gl.GenVertexArray()
	defer vao_A.Delete()
	vao_A.Bind()
	vertexPosLoc := prog.GetAttribLocation("vertexPosition_modelspace")
	vertexPosLoc.EnableArray()
	vertexNormLoc := prog.GetAttribLocation("vertexNormal_modelspace")
	vertexNormLoc.EnableArray()

	indexBuffer := gl.GenBuffer()
	defer indexBuffer.Delete()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(int(0))), indices, gl.STATIC_DRAW)

	verticesBuffer := gl.GenBuffer()
	defer verticesBuffer.Delete()
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

	InitScreenQuad()

	releaseTextures := initTextures()
	defer releaseTextures()
	gl.ActiveTexture(gl.TEXTURE4)
	heights.Texture()
	gl.ActiveTexture(gl.TEXTURE5)

	ps := NewParticleSystem(100000, mathgl.Vec3f{32, 32, 32}, 0.1, 500)

	gamestate := GameState{
		nil,
		heights,
		ps,
		&MyPlayer{Camera{mathgl.Vec3f{5, 5, 10}, mathgl.QuatIdentf()}, PlayerInput{}, mathgl.Vec3f{}},
	}

	gamestate.Camera = gamestate.Player.GetCamera()

	InitInput(&gamestate)

	gl.Enable(gl.DEPTH_TEST)

	ps.Program.Use()
	ps.Program.GetUniformLocation("heights").Uniform1i(4)
	ps.Program.GetUniformLocation("lowerBound").Uniform3f(0, 0, min_h)
	ps.Program.GetUniformLocation("upperBound").Uniform3f(w, h, max_h)

	gl.PointSize(4)

	vao_C := gl.GenVertexArray()

	gl.Enable(gl.CULL_FACE)

	for ok := true; ok; ok = (glfw.Key(glfw.KeyEsc) != glfw.KeyPress && glfw.WindowParam(glfw.Opened) == gl.TRUE) {
		Input(&gamestate)

		gamestate.Player.Update(&gamestate)

		mat2 := gamestate.Camera.View()

		prog.Use()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		timeLoc.Uniform1f(float32(glfw.Time()))
		seaLevelLoc.Uniform1f(float32(math.Sin(glfw.Time()*0.1)*10 - 5))
		highlightLoc.Uniform1f(float32(highlight))

		gl.Disable(gl.BLEND)

		prog.Use()
		vao_A.Bind()
		numverts := len(indices)

		mathgl.Translate3D(w, 0, 0)

		projView := mat1.Mul4(mat2)

		for i := -2; i <= 2; i++ {
			for j := -2; j <= 2; j++ {
				modelMat := mathgl.Translate3D(float64(i*w), float64(j*h), 0)
				finalMat := projView.Mul4(modelMat)
				matrixLoc.UniformMatrix4f(false, (*[16]float32)(&finalMat))
				modelLoc.UniformMatrix4f(false, (*[16]float32)(&modelMat))
				gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
			}
		}

		finalMat := projView

		vao_C.Bind()
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
		gl.Enable(gl.BLEND)
		gl.DepthMask(false)

		ps.DoStep()
		ps.Render(&finalMat)

		//RenderScreenQuad()

		gl.DepthMask(true)

		tw.Draw()

		glfw.SwapBuffers()
	}
}
