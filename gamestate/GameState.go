package gamestate

import (
	"fmt"
	"github.com/go-gl-legacy/gl"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/krux02/tw"
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
)

type GameState struct {
	Window  *sdl.Window
	Camera  *Camera
	Bar     *tw.Bar
	World   *World
	Fps     float32
	Options settings.BoolOptions
}

func NewGameState(window *sdl.Window, world *World) (gamestate *GameState) {
	gl.ClearColor(0., 0., 0.4, 0.0)

	gl.Enable(gl.DEPTH_TEST)

	bar := tw.NewBar("TweakBar")

	gamestate = &GameState{
		Window:  window,
		Camera:  nil,
		Bar:     bar,
		World:   world,
		Fps:     0,
		Options: settings.BoolOptions{},
	}

	opt := &gamestate.Options
	opt.Load()
	gamestate.Camera = gamestate.World.Player.GetCamera()

	tw.Define(" GLOBAL help='This example shows how to integrate AntTweakBar with SDL2 and OpenGL.' ")
	bar.AddVarRO("fps", tw.TYPE_FLOAT, unsafe.Pointer(&gamestate.Fps), "")
	opt.CreateGui(bar)

	for i, portal := range gamestate.World.KdTree {
		ptr := &(portal.(*Portal).Orientation)
		bar.AddVarRW(fmt.Sprintf("Rotation %d", i), tw.TYPE_QUAT4F, unsafe.Pointer(ptr), "")
	}

	//window.GetSize(w, h)
	w, h := window.GetSize()
	tw.WindowSize(w, h)

	return
}

func (this *GameState) Delete() {
	this.Options.Save()
	this.Bar.Delete()
}
