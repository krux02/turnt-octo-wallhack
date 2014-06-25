package main

import (
	//"fmt"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/debug"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/tw"
	"github.com/veandco/go-sdl2/sdl"
	//"math"
)

var drag uint8

func GrabCursor() {
	if sdl.GetRelativeMouseMode() {
		sdl.SetRelativeMouseMode(false)
	} else {
		sdl.SetRelativeMouseMode(true)
	}
}

func GetMouseDirection(window *sdl.Window, mx, my int) (dir_cs mgl.Vec4) {
	if !sdl.GetRelativeMouseMode() {
		W, H := window.GetSize()
		x := (2*float32(mx) - float32(W)) / float32(H)
		y := (float32(H) - 2*float32(my)) / float32(H)
		dir_cs = mgl.Vec4{x, y, -1, 0}
	} else {
		dir_cs = mgl.Vec4{0, 0, -1, 0}
	}
	return
}

func RayCastInCameraSpace(gs *gamestate.GameState, dir_cs mgl.Vec4) (hit_ws mgl.Vec4, hit bool) {
	pos_ws := gs.Camera.Position
	dir_ws := gs.Camera.Model().Mul4x1(dir_cs)

	pos_ws_v3 := helpers.XYZ(pos_ws)
	dir_ws_v3 := helpers.XYZ(dir_ws)

	factor, hit := gs.World.HeightMap.RayCast(pos_ws_v3, dir_ws_v3)
	if hit {
		hit_ws = pos_ws.Add(dir_ws.Mul(factor))
	}
	return
}

// returns false if the player wants to quit
func Input(gs *gamestate.GameState, worldRenderer *rendering.WorldRenderer) bool {
	running := true
	window := gs.Window
	inp := gamestate.PlayerInput{}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		consumeEvent := true
		if !sdl.GetRelativeMouseMode() {
			consumeEvent = tw.EventSDL(event, 2, 0)
		}
		if consumeEvent {
			switch e := event.(type) {
			case *sdl.WindowEvent:
				switch e.Event {
				case sdl.WINDOWEVENT_CLOSE:
					running = false
				case sdl.WINDOWEVENT_RESIZED:
					width, height := int(e.Data1), int(e.Data2)
					worldRenderer.Resize(width, height)
					tw.WindowSize(width, height)
				}
			case *sdl.MouseButtonEvent:
				button := e.Button

				if e.State == sdl.PRESSED && drag == 255 {
					drag = e.Button
				}
				if e.State == sdl.RELEASED && drag == button {
					drag = 255
				}

				if e.State == sdl.PRESSED && button == sdl.BUTTON_RIGHT {
					GrabCursor()
				}

				// ray cast testing
				if e.State == sdl.PRESSED && button == sdl.BUTTON_LEFT {
					dir_cs := GetMouseDirection(window, int(e.X), int(e.Y))
					out, hit := RayCastInCameraSpace(gs, dir_cs)
					if hit {
						n := helpers.Vector(gs.World.HeightMap.Normal2f(out[0], out[1]))
						debug.Color(mgl.Vec4{0, 1, 0, 1})
						debug.Line(out, out.Add(n))
					}
				}
			case *sdl.KeyDownEvent:
				switch e.Keysym.Scancode {
				case sdl.SCANCODE_RETURN:
					GrabCursor()
				case sdl.SCANCODE_SPACE:
				case sdl.SCANCODE_ESCAPE:
					running = false
				case sdl.SCANCODE_F3:
					worldRenderer.ToggleRift()
					if worldRenderer.RiftRender() {
						window.SetSize(1280, 800)
					}
				}
			case *sdl.KeyUpEvent:
			}
		}
	}

	if sdl.GetRelativeMouseMode() {
		x, y, _ := sdl.GetRelativeMouseState()
		inp.Rotate[0] = -float32(y) / 500
		inp.Rotate[1] = -float32(x) / 500
	}

	keyState := sdl.GetKeyboardState()

	if keyState[sdl.SCANCODE_E] == 1 {
		inp.Move[2] -= 1
	}
	if keyState[sdl.SCANCODE_D] == 1 {
		inp.Move[2] += 1
	}
	if keyState[sdl.SCANCODE_S] == 1 {
		inp.Move[0] -= 1
	}
	if keyState[sdl.SCANCODE_F] == 1 {
		inp.Move[0] += 1
	}
	if keyState[sdl.SCANCODE_R] == 1 {
		inp.Rotate[2] -= 0.01
	}
	if keyState[sdl.SCANCODE_W] == 1 {
		inp.Rotate[2] += 0.01
	}

	x, y, button_state := sdl.GetMouseState()
	if button_state&sdl.BUTTON_LEFT != 0 {
		dir_cs := GetMouseDirection(window, x, y)
		out, hit := RayCastInCameraSpace(gs, dir_cs)
		if hit {
			heightMap := gs.World.HeightMap
			heightMap.Bump(mgl.Vec2{out[0], out[1]}, 3)
		}
	}

	gs.World.Player.Input = inp

	return running
}
