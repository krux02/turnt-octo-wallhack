package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/mathgl"
	"github.com/krux02/tw"
	"os"
	"unsafe"
)

type GameState struct {
	Window          *glfw.Window
	Camera          *Camera
	Proj            mathgl.Mat4f
	HeightMap       *HeightMap
	PalmTrees       *PalmTrees
	ParticleSystem  *ParticleSystem
	WordlRenderer   *WorldRenderer
	Player          Player
	Fps             float32
	ParticleRender  bool
	ParticlePhysics bool
	WorldRender     bool
	TreeRender      bool
	PlayerPhysics   bool
}

const vertexStride = int(unsafe.Sizeof(Vertex{}))

const w = 64
const h = 64

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

	var fps float32
	bar.AddVarRO("fps", tw.TYPE_FLOAT, unsafe.Pointer(&fps), "")

	initDebugContext()

	// heights := NewHeightMapFramFile("test.png")
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

	ps.TransformProg.Use()

	ps.TransformLoc.Heights.Uniform1i(4)
	ps.TransformLoc.LowerBound.Uniform3f(0, 0, min_h)
	ps.TransformLoc.UpperBound.Uniform3f(w, h, max_h)

	gl.PointSize(4)

	//gl.Enable(gl.CULL_FACE)

	gamestate := GameState{
		window,
		nil,
		mathgl.Perspective(90, 4.0/3.0, 0.01, 1000),
		heights,
		NewPalmTrees(heights, 25000),
		ps,
		wr,
		&MyPlayer{Camera{mathgl.Vec3f{5, 5, 10}, mathgl.QuatIdentf()}, PlayerInput{}, mathgl.Vec3f{}},
		0,
		true,
		true,
		true,
		true,
		true,
	}
	gamestate.Camera = gamestate.Player.GetCamera()

	bar.AddVarRW("ParticleRender", tw.TYPE_BOOL8, unsafe.Pointer(&gamestate.ParticleRender), "")
	bar.AddVarRW("ParticlePhysics", tw.TYPE_BOOL8, unsafe.Pointer(&gamestate.ParticlePhysics), "")
	bar.AddVarRW("WorldRender", tw.TYPE_BOOL8, unsafe.Pointer(&gamestate.WorldRender), "")
	bar.AddVarRW("TreeRender", tw.TYPE_BOOL8, unsafe.Pointer(&gamestate.TreeRender), "")
	bar.AddVarRW("PlayerPhysics", tw.TYPE_BOOL8, unsafe.Pointer(&gamestate.PlayerPhysics), "")
	bar.AddButton("save image", func() { SaveImage("test.png", heights.ExportImage()) }, "")

	wireframe := false

	setCallback := func(value unsafe.Pointer) {
		if( *(*bool)(value) ) {
			gl.PolygonMode( gl.FRONT_AND_BACK, gl.LINE );
			wireframe = true
		} else {
			gl.PolygonMode( gl.FRONT_AND_BACK, gl.FILL );
			wireframe = false
		}
	}

	getCallback := func(value unsafe.Pointer) {
		*(*bool)(value) = wireframe
	}

	bar.AddVarCB("wireframe", tw.TYPE_BOOL8, setCallback, getCallback, nil, "")

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

	for !window.ShouldClose() && window.GetKey(glfw.KeyEscape) != glfw.Press {
		currentTime := glfw.GetTime()

		if currentTime > time+1 {
			gamestate.Fps = float32(frames)
			frames = 0
			time = currentTime
		}

		Input(gamestate)

		gamestate.Player.Update(gamestate)

		Proj := gamestate.Proj
		View := gamestate.Camera.View()
		ProjView := Proj.Mul4(View)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Disable(gl.BLEND)

		if gamestate.WorldRender {
			gamestate.WordlRenderer.Render(gamestate.HeightMap, Proj, View, currentTime, highlight)
		}
		if gamestate.TreeRender {
			gamestate.PalmTrees.Render(Proj, View)
		}
		if gamestate.ParticlePhysics {
			gamestate.ParticleSystem.DoStep(currentTime)
		}
		if gamestate.ParticleRender {
			gamestate.ParticleSystem.Render(&ProjView)
		}
		//RenderScreenQuad()

		// gl.DepthMask(true)

		tw.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
