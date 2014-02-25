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
	Bar     *tw.Bar
	World   *World
	Player  *Player
	Fps     float32
	Options settings.BoolOptions
}

func NewGameState(window *glfw.Window, world *World) (gamestate *GameState) {
	gl.ClearColor(0., 0., 0.4, 0.0)

	gl.ActiveTexture(gl.TEXTURE5)

	gl.Enable(gl.DEPTH_TEST)

	bar := tw.NewBar("TweakBar")

	startPos := mgl.Vec4f{5, 5, 10, 1}

	gamestate = &GameState{
		Window:  window,
		Camera:  nil,
		Bar:     bar,
		World:   world,
		Player:  &Player{*NewCameraFromPos4f(startPos), PlayerInput{}, mgl.Vec4f{}},
		Fps:     0,
		Options: settings.BoolOptions{StartPosition: startPos},
	}

	opt := &gamestate.Options
	opt.Load()
	gamestate.Camera = gamestate.Player.GetCamera()

	tw.Define(" GLOBAL help='This example shows how to integrate AntTweakBar with GLFW and OpenGL.' ")

	bar.AddVarRO("fps", tw.TYPE_FLOAT, unsafe.Pointer(&gamestate.Fps), "")
	bar.AddVarRW("NoParticleRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoParticleRender), "")
	bar.AddVarRW("NoParticlePhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoParticlePhysics), "")
	bar.AddVarRW("NoWorldRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoWorldRender), "")
	bar.AddVarRW("NoTreeRender", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoTreeRender), "")
	bar.AddVarRW("NoPlayerPhysics", tw.TYPE_BOOL8, unsafe.Pointer(&opt.NoPlayerPhysics), "")
	bar.AddVarRW("Wireframe", tw.TYPE_BOOL8, unsafe.Pointer(&opt.Wireframe), "")
	bar.AddVarRW("DepthClamp", tw.TYPE_BOOL8, unsafe.Pointer(&opt.DepthClamp), "")

	for i, portal := range gamestate.World.KdTree {
		ptr := &(portal.(*Portal).Orientation)
		bar.AddVarRW(fmt.Sprintf("Rotation %d", i), tw.TYPE_QUAT4F, unsafe.Pointer(ptr), "")
	}

	bar.AddButton("save image", func() { helpers.SaveImage("test.png", world.HeightMap.ExportImage()) }, "")

	w, h := window.GetSize()
	tw.WindowSize(w, h)

	return
}

func (this *GameState) Delete() {
	this.Options.Save()
	this.Bar.Delete()
}
