package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/mathgl"
	"github.com/krux02/tw"
	"math"
	"os"
	"unsafe"
)

type GameState struct {
	Window         *glfw.Window
	Camera         *Camera
	Proj           mathgl.Mat4f
	HeightMap      *HeightMap
	PalmTrees      *PalmTrees
	ParticleSystem *ParticleSystem
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

	ps.Locations.Heights.Uniform1i(4)
	ps.Locations.LowerBound.Uniform3f(0, 0, min_h)
	ps.Locations.UpperBound.Uniform3f(w, h, max_h)

	gl.PointSize(4)

	gl.Enable(gl.CULL_FACE)

	gamestate := GameState{
		window,
		nil,
		mathgl.Perspective(90, 4.0/3.0, 0.01, 1000),
		heights,
		NewPalmTrees(heights, 128),
		ps,
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

	for !window.ShouldClose() && window.GetKey(glfw.KeyEscape) != glfw.Press {

		if glfw.GetTime() > time+1 {
			gamestate.fps = float32(frames)
			frames = 0
			time = glfw.GetTime()
		}

		Input(gamestate)

		gamestate.Player.Update(gamestate)

		View := gamestate.Camera.View()
		projView := gamestate.Proj.Mul4(View)

		gamestate.WordlRenderer.Program.Use()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		Loc := gamestate.WordlRenderer.WorldRenLoc

		Loc.Time.Uniform1f(float32(glfw.GetTime()))
		Loc.SeaLevel.Uniform1f(float32(math.Sin(glfw.GetTime()*0.1)*10 - 5))
		Loc.Highlight.Uniform1f(float32(highlight))

		gl.Disable(gl.BLEND)

		gamestate.WordlRenderer.Render(gamestate)

		gamestate.PalmTrees.Render(gamestate.Proj, View)

		gamestate.ParticleSystem.DoStep()

		finalMat := projView
		gamestate.ParticleSystem.Render(&finalMat)

		//RenderScreenQuad()

		gl.DepthMask(true)

		tw.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
