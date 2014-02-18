package gamestate

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/krux02/tw"
	"unsafe"
)

type GameState struct {
	Window  *glfw.Window
	Camera  *Camera
	Proj    mgl.Mat4f
	Bar     *tw.Bar
	World   *World
	Player  *Player
	Fps     float32
	Options settings.BoolOptions
}

func NewGameState(window *glfw.Window) (gamestate *GameState) {
	gl.ClearColor(0., 0., 0.4, 0.0)

	World := NewWorld()

	gl.ActiveTexture(gl.TEXTURE4)
	World.HeightMap.Texture()
	gl.ActiveTexture(gl.TEXTURE5)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	bar := tw.NewBar("TweakBar")

	startPos := mgl.Vec3f{5, 5, 10}

	gamestate = &GameState{
		Window:  window,
		Camera:  nil,
		Proj:    mgl.Perspective(90, 4.0/3.0, 0.001, 1000),
		Bar:     bar,
		World:   World,
		Player:  &Player{Camera{startPos, mgl.QuatIdentf()}, PlayerInput{}, mgl.Vec3f{}},
		Fps:     0,
		Options: settings.BoolOptions{StartPosition: startPos},
	}

	gamestate.Camera = gamestate.Player.GetCamera()
	opt := &gamestate.Options
	opt.NoTreeRender = true
	opt.NoParticlePhysics = true
	opt.NoParticleRender = true

	tw.Define(" GLOBAL help='This example shows how to integrate AntTweakBar with GLFW and OpenGL.' ")

	bar.AddVarRO("fps", tw.TYPE_FLOAT, unsafe.Pointer(&gamestate.Fps), "")
	bar.AddVarRW("NoParticleRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoParticleRender), "")
	bar.AddVarRW("NoParticlePhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoParticlePhysics), "")
	bar.AddVarRW("NoWorldRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoWorldRender), "")
	bar.AddVarRW("NoTreeRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoTreeRender), "")
	bar.AddVarRW("NoPlayerPhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoPlayerPhysics), "")
	bar.AddVarRW("Wireframe", tw.TYPE_BOOL8, unsafe.Pointer(&opt.Wireframe), "")
	bar.AddVarRW("DepthClamp", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DepthClamp), "")

	N := len(gamestate.World.Portals)
	for i := 0; i < N; i++ {
		ptr := &gamestate.World.Portals[i].Orientation
		bar.AddVarRW(fmt.Sprintf("Rotation %d", i), tw.TYPE_QUAT4F, unsafe.Pointer(ptr), "")
	}

	bar.AddButton("save image", func() { helpers.SaveImage("test.png", World.HeightMap.ExportImage()) }, "")

	window.SetSizeCallback(func(window *glfw.Window, width, height int) {
		gl.Viewport(0, 0, width, height)
		gamestate.Proj = mgl.Perspective(90, float32(width)/float32(height), 0.1, 1000)
		tw.WindowSize(width, height)
	})

	w, h := window.GetSize()
	tw.WindowSize(w, h)

	return
}

func (this *GameState) Delete() {
	this.Bar.Delete()
}
