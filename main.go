package main

import (
	"fmt"
	//	"github.com/Jragonmiris/mathgl/examples/opengl-tutorial/helper"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/krux02/mathgl"
	"github.com/krux02/tw"
	"math"
	"os"
	"unsafe"
)

type GameState struct {
	Camera         *Camera
	Proj           mathgl.Mat4f
	HeightMap      *HeightMap
	ParticleSystem *ParticleSystem
	ParticlesVAO   gl.VertexArray
	WordlRenderer  *WorldRenderer
	Player         Player
	fps            float32
}

const vertexStride = int(unsafe.Sizeof(Vertex{}))

const w = 128
const h = 128

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

	var fps float32
	bar.AddButton("but1", func() { fmt.Printf("but1 %d\n", xxx); xxx += 1 }, "")
	bar.AddButton("but2", func() { fmt.Printf("but2 %d\n", xxx); xxx += 1 }, "")
	bar.AddVarRO("fps", tw.TYPE_FLOAT, unsafe.Pointer(&fps), "")

	initDebugContext()

	heights := NewHeightMap(w, h)
	heights.DiamondSquare(w)
	min_h, max_h := heights.Bounds()

	wr := NewWorldRenderer(heights)

	gl.ClearColor(0., 0., 0.4, 0.)

	InitScreenQuad()

	releaseTextures := initTextures()
	defer releaseTextures()
	gl.ActiveTexture(gl.TEXTURE4)
	heights.Texture()
	gl.ActiveTexture(gl.TEXTURE5)

	ps := NewParticleSystem(100000, mathgl.Vec3f{32, 32, 32}, 0.1, 500)

	gl.Enable(gl.DEPTH_TEST)

	ps.Program.Use()
	ps.Program.GetUniformLocation("heights").Uniform1i(4)
	ps.Program.GetUniformLocation("lowerBound").Uniform3f(0, 0, min_h)
	ps.Program.GetUniformLocation("upperBound").Uniform3f(w, h, max_h)

	gl.PointSize(4)

	vao_C := gl.GenVertexArray()

	gl.Enable(gl.CULL_FACE)

	gamestate := GameState{
		nil,
		mathgl.Perspective(90, 4.0/3.0, 0.01, 1000),
		heights,
		ps,
		vao_C,
		wr,
		&MyPlayer{Camera{mathgl.Vec3f{5, 5, 10}, mathgl.QuatIdentf()}, PlayerInput{}, mathgl.Vec3f{}},
		0,
	}
	gamestate.Camera = gamestate.Player.GetCamera()

	glfw.SetWindowSizeCallback(func(width int, height int) {
		gl.Viewport(0, 0, width, height)
		gamestate.Proj = mathgl.Perspective(90, float64(width)/float64(height), 0.1, 1000)
		tw.WindowSize(width, height)
	})

	InitInput(&gamestate)

	MainLoop(&gamestate)
}

func MainLoop(gamestate *GameState) {
	var frames int
	time := glfw.Time()
	for ok := true; ok; ok = (glfw.Key(glfw.KeyEsc) != glfw.KeyPress && glfw.WindowParam(glfw.Opened) == gl.TRUE) {
		frames += 1

		if glfw.Time() > time+1 {
			gamestate.fps = float32(frames)
			frames = 0
			time = glfw.Time()
		}

		Input(gamestate)

		gamestate.Player.Update(gamestate)

		view := gamestate.Camera.View()
		projView := gamestate.Proj.Mul4(view)

		gamestate.WordlRenderer.Program.Use()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		Loc := gamestate.WordlRenderer.WorldRenLoc

		Loc.time.Uniform1f(float32(glfw.Time()))
		Loc.seaLevel.Uniform1f(float32(math.Sin(glfw.Time()*0.1)*10 - 5))
		Loc.highlight.Uniform1f(float32(highlight))

		gl.Disable(gl.BLEND)

		gamestate.WordlRenderer.Render(gamestate)

		gamestate.ParticlesVAO.Bind()

		gamestate.ParticleSystem.DoStep()

		for i := -2; i <= 2; i++ {
			for j := -2; j <= 2; j++ {
				modelMat := mathgl.Translate3D(float64(i*w), float64(j*h), 0)
				finalMat := projView.Mul4(modelMat)
				gamestate.ParticleSystem.Render(&finalMat)
			}
		}

		//RenderScreenQuad()

		gl.DepthMask(true)

		tw.Draw()

		glfw.SwapBuffers()
	}
}
