package main

import (
	// 	// "fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	mgl "github.com/krux02/mathgl"
	// "github.com/krux02/turnt-octo-wallhack/debugContext"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/turnt-octo-wallhack/world"
	"github.com/krux02/tw"
	"unsafe"
)

type GameState struct {
	Window          *glfw.Window
	Camera          *Camera
	Proj            mgl.Mat4f
	Textures        []gl.Texture
	Bar             *tw.Bar
	World           *world.World
	Portal          *rendering.MeshRenderData
	PortalPositions []mgl.Vec3f
	PalmTrees       *rendering.PalmTrees
	ParticleSystem  *particles.ParticleSystem
	WordlRenderer   *rendering.WorldRenderer
	MeshRenderer    *rendering.MeshRenderer
	Player          Player
	Fps             float32
	Options         BoolOptions
}

func NewGameState(window *glfw.Window) (gamestate *GameState) { 	
	gl.ClearColor(0., 0., 0.4, 0.0)
	

	World := world.NewWorld()
	min_h, max_h := World.HeightMap.Bounds()

	wr := rendering.NewWorldRenderer(World.HeightMap)
	rendering.InitScreenQuad()

	textures := initTextures()
	gl.ActiveTexture(gl.TEXTURE4)
	World.HeightMap.Texture()
	gl.ActiveTexture(gl.TEXTURE5)

	ps := particles.NewParticleSystem(75000, mgl.Vec3f{32, 32, 32}, 1, 250)

	gl.Enable(gl.DEPTH_TEST)

	ps.TransformProg.Use()
	ps.TransformLoc.Heights.Uniform1i(4)
	ps.TransformLoc.LowerBound.Uniform3f(0, 0, min_h)
	ps.TransformLoc.UpperBound.Uniform3f(world.W, world.H, max_h)

	gl.Enable(gl.CULL_FACE)

	meshRenderer := rendering.NewMeshRenderer()
	portalData := meshRenderer.CreateMeshRenderData(World.Portal)

	bar := tw.NewBar("TweakBar")

	gamestate = &GameState{
		Window:          window,
		Camera:          nil,
		Proj:            mgl.Perspective(90, 4.0/3.0, 0.01, 1000),
		Textures:        textures,
		Bar:             bar,
		World:           World,
		Portal:          &portalData,
		PortalPositions: []mgl.Vec3f{mgl.Vec3f{10, 10, 15}},//, mgl.Vec3f{30, 30, 10}, mgl.Vec3f{60, 60, 9}},
		PalmTrees:       rendering.NewPalmTrees(World.HeightMap, 10000),
		ParticleSystem:  ps,
		WordlRenderer:   wr,
		MeshRenderer:    &meshRenderer,
		Player:          &MyPlayer{Camera{mgl.Vec3f{5, 5, 10}, mgl.QuatIdentf()}, PlayerInput{}, mgl.Vec3f{}},
		Fps:             0,
		Options:         BoolOptions{},
	}

	gamestate.Camera = gamestate.Player.GetCamera()
	opt := &gamestate.Options

	tw.Define(" GLOBAL help='This example shows how to integrate AntTweakBar with GLFW and OpenGL.' ")
	bar.AddVarRO("fps", tw.TYPE_FLOAT, unsafe.Pointer(&gamestate.Fps), "")
	bar.AddVarRW("DisableParticleRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableParticleRender), "")
	bar.AddVarRW("DisableParticlePhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableParticlePhysics), "")
	bar.AddVarRW("DisableWorldRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableWorldRender), "")
	bar.AddVarRW("DisableTreeRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisableTreeRender), "")
	bar.AddVarRW("DisablePlayerPhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DisablePlayerPhysics), "")
	bar.AddVarRW("Wireframe", tw.TYPE_BOOL8, unsafe.Pointer(&opt.Wireframe), "")


	bar.AddButton("save image", func() { helpers.SaveImage("test.png", World.HeightMap.ExportImage()) }, "")

	window.SetSizeCallback(func(window *glfw.Window, width, height int) {
		gl.Viewport(0, 0, width, height)
		gamestate.Proj = mgl.Perspective(90, float64(width)/float64(height), 0.1, 1000)
		tw.WindowSize(width, height)
	})

	w, h := window.GetSize()
	tw.WindowSize(w, h)

	return
}

func (this *GameState) Delete() {
	gl.DeleteTextures(this.Textures)
	this.Bar.Delete()
}
