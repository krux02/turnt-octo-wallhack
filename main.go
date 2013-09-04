package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
	"github.com/krux02/mathgl"
	"github.com/krux02/tw"
	"io/ioutil"
	"math"
	"os"
	"unsafe"
)

func MakeProgram(vertFname, fragFname string) gl.Program {
	vertSource, err := ioutil.ReadFile(vertFname)
	if err != nil {
		panic(err)
	}

	fragSource, err := ioutil.ReadFile(fragFname)
	if err != nil {
		panic(err)
	}

	return glh.NewProgram(glh.Shader{gl.VERTEX_SHADER, string(vertSource)}, glh.Shader{gl.FRAGMENT_SHADER, string(fragSource)})
}

type GameState struct {
	Window         *glfw.Window
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

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	glfw.Init()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglDebugContext, gl.TRUE)

	window, err := glfw.CreateWindow(1024, 768, "gogog", nil, nil)
	if window == nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	window.MakeContextCurrent()

	gl.Init()
	gl.GetError() // Ignore error

	window.SetInputMode(glfw.StickyKeys, gl.TRUE)

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
		window,
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

	window.SetSizeCallback(func(window *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, width, height)
		gamestate.Proj = mathgl.Perspective(90, float64(width)/float64(height), 0.1, 1000)
		tw.WindowSize(width, height)
	})

	tw.WindowSize(1024, 768)

	InitInput(&gamestate)

	MainLoop(&gamestate)
}

func MainLoop(gamestate *GameState) {
	var frames int
	time := glfw.GetTime()

	window := gamestate.Window

	ok := true

	window.SetCloseCallback(func(window *glfw.Window) { ok = false })

	for ok {
		ok = ok && window.GetKey(glfw.KeyEscape) != glfw.Press

		if glfw.GetTime() > time+1 {
			gamestate.fps = float32(frames)
			frames = 0
			time = glfw.GetTime()
		}

		Input(gamestate)

		gamestate.Player.Update(gamestate)

		view := gamestate.Camera.View()
		projView := gamestate.Proj.Mul4(view)

		gamestate.WordlRenderer.Program.Use()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		Loc := gamestate.WordlRenderer.WorldRenLoc

		Loc.time.Uniform1f(float32(glfw.GetTime()))
		Loc.seaLevel.Uniform1f(float32(math.Sin(glfw.GetTime()*0.1)*10 - 5))
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

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
