package main

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/krux02/turnt-octo-wallhack/debug"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/tw"
	//"math"
)

var drag uint8

func GrabCursor() {
	if sdl.GetRelativeMouseMode() {
		sdl.SetRelativeMouseMode(false)
		//window.SetCursorPosition(0, 0)
	} else {
		sdl.SetRelativeMouseMode(true)
	}
}

// returns false if the player wants to quit
func Input(gs *gamestate.GameState, worldRenderer *rendering.WorldRenderer) bool {
	running := true

	window := gs.Window
	inp := gamestate.PlayerInput{}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		tw.EventSDL(event, 1, 3)
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

			if e.State == sdl.PRESSED && button == 1 {
				GrabCursor()
			}

			// ray cast testing
			if e.State == sdl.PRESSED && button == 0 {
				var dir_cs mgl.Vec4f
				if !sdl.GetRelativeMouseMode() {
					mx, my := e.X, e.Y
					W, H := window.GetSize()
					x := (2*float32(mx) - float32(W)) / float32(H)
					y := (float32(H) - 2*float32(my)) / float32(H)
					dir_cs = mgl.Vec4f{x, y, -1, 0}
				} else {
					dir_cs = mgl.Vec4f{0, 0, -1, 0}
				}

				m := gs.Camera.Model()
				pos_ws := gs.Camera.Position
				dir_ws := m.Mul4x1(dir_cs)

				pos_ws_v3 := helpers.XYZ(pos_ws)
				dir_ws_v3 := helpers.XYZ(dir_ws)

				factor, hit := gs.World.HeightMap.RayCast(pos_ws_v3, dir_ws_v3)

				if hit {
					out := pos_ws.Add(dir_ws.Mul(factor))
					n := helpers.Vector(gs.World.HeightMap.Normal2f(out[0], out[1]))
					debug.Color(mgl.Vec4f{0, 1, 0, 1})
					debug.Line(out, out.Add(n))

					heightMap := gs.World.HeightMap
					heightMap.Bump(mgl.Vec2f{out[0], out[1]}, 3)
				}
			}
		case *sdl.KeyDownEvent:
			switch e.Keysym.Scancode {
			case sdl.SCANCODE_RETURN:
				GrabCursor()
			case sdl.SCANCODE_SPACE:
				gs.Player.Camera.Position = gs.Options.StartPosition
			case sdl.SCANCODE_ESCAPE:
				running = false
			}
		case *sdl.KeyUpEvent:
		}
	}

	if sdl.GetRelativeMouseMode() {
		x, y, _ := sdl.GetRelativeMouseState()
		inp.Rotate[0] = -float32(y) / 5
		inp.Rotate[1] = -float32(x) / 5
	}

	state := sdl.GetKeyboardState()

	if state[sdl.SCANCODE_E] == 1 {
		inp.Move[2] -= 1
	}
	if state[sdl.SCANCODE_D] == 1 {
		inp.Move[2] += 1
	}
	if state[sdl.SCANCODE_S] == 1 {
		inp.Move[0] -= 1
	}
	if state[sdl.SCANCODE_F] == 1 {
		inp.Move[0] += 1
	}
	if state[sdl.SCANCODE_R] == 1 {
		inp.Rotate[2] -= 1
	}
	if state[sdl.SCANCODE_W] == 1 {
		inp.Rotate[2] += 1
	}

	gs.Player.Input = inp

	return running
}
