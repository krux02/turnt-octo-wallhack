package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/mathgl"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/turnt-octo-wallhack/world"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/tw"
	"os"
	"unsafe"
)

type GameState struct {
	Window         *glfw.Window
	Camera         *Camera
	Proj           mathgl.Mat4f
	HeightMap      *world.HeightMap
	PalmTrees      *rendering.PalmTrees
	ParticleSystem *particles.ParticleSystem
	WordlRenderer  *rendering.WorldRenderer
	Player         Player
	Fps            float32
	Options        BoolOptions
}

type BoolOptions struct {
	DisableParticleRender,
	DisableParticlePhysics,
	DisableWorldRender,
	DisableTreeRender,
	DisablePlayerPhysics,
	Wireframe bool
}

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
	heights := world.NewHeightMap(w, h)
	heights.DiamondSquare(w)
	min_h, max_h := heights.Bounds()

	wr := rendering.NewWorldRenderer(heights)

	gl.ClearColor(0., 0., 0.4, 0.)

	rendering.InitScreenQuad()

	releaseTextures := initTextures()
	defer releaseTextures()
	gl.ActiveTexture(gl.TEXTURE4)
	heights.Texture()


	gl.ActiveTexture(gl.TEXTURE5)

	ps := particles.NewParticleSystem(100000, mathgl.Vec3f{32, 32, 32}, 0.1, 500)

	gl.Enable(gl.DEPTH_TEST)

	ps.TransformProg.Use()

	ps.TransformLoc.Heights.Uniform1i(4)
	ps.TransformLoc.LowerBound.Uniform3f(0, 0, min_h)
	ps.TransformLoc.UpperBound.Uniform3f(w, h, max_h)

	//gl.Enable(gl.CULL_FACE)

	gamestate := GameState{
		window,
		nil,
		mathgl.Perspective(90, 4.0/3.0, 0.01, 1000),
		heights,
		rendering.NewPalmTrees(heights, 25000),
		ps,
		wr,
		&MyPlayer{Camera{mathgl.Vec3f{5, 5, 10}, mathgl.QuatIdentf()}, PlayerInput{}, mathgl.Vec3f{}},
		0,
		BoolOptions{},
	}
	gamestate.Camera = gamestate.Player.GetCamera()
	opt := &gamestate.Options

	bar.AddVarRW("DisableParticleRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableParticleRender), "")
	bar.AddVarRW("DisableParticlePhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableParticlePhysics), "")
	bar.AddVarRW("DisableWorldRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableWorldRender), "")
	bar.AddVarRW("DisableTreeRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableTreeRender), "")
	bar.AddVarRW("DisablePlayerPhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisablePlayerPhysics), "")
	bar.AddVarRW("Wireframe", tw.TYPE_BOOL8, unsafe.Pointer(&opt.Wireframe), "")
	bar.AddButton("save image", func() { helpers.SaveImage("test.png", heights.ExportImage()) }, "")

	window.SetSizeCallback(func(window *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, width, height)
		gamestate.Proj = mathgl.Perspective(90, float64(width)/float64(height), 0.1, 1000)
		tw.WindowSize(width, height)
	})

	tw.WindowSize(1024, 768)

	InitInput(&gamestate)

	MainLoop(&gamestate)
}
