package main

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/debug"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/tw"
	"math"
)

var drag glfw.MouseButton = -1

/*
var lastMousePos = mgl.Vec2f{0, 0}
var currentMousePos func() mgl.Vec2f
var updateLastMousePos func()
*/

func GrabCursor(window *glfw.Window) {
	switch window.GetInputMode(glfw.Cursor) {
	case glfw.CursorNormal:
		window.SetInputMode(glfw.Cursor, glfw.CursorDisabled)
		window.SetCursorPosition(0, 0)
	case glfw.CursorDisabled:
		window.SetInputMode(glfw.Cursor, glfw.CursorNormal)
	}
}

func InitInput(gs *gamestate.GameState) {
	window := gs.Window

	MouseButton := func(window *glfw.Window, button glfw.MouseButton, state glfw.Action, modifiers glfw.ModifierKey) {
		if state == glfw.Press && drag == -1 {
			drag = button
		}
		if state == glfw.Release && drag == button {
			drag = -1
		}

		if state == glfw.Press && button == 1 {
			GrabCursor(window)
		}

		// ray cast testing
		if state == glfw.Press && button == 0 {

			var dir_cs mgl.Vec4f
			if window.GetInputMode(glfw.Cursor) == glfw.CursorNormal {
				mx, my := window.GetCursorPosition()
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

		tw.EventMouseButtonGLFW(int(button), int(state))
	}

	KeyPress := func(window *glfw.Window, key glfw.Key, _ int, state glfw.Action, modifiers glfw.ModifierKey) {
		if state == glfw.Press {
			switch key {
			case glfw.KeyEnter:
				GrabCursor(window)
			case glfw.KeySpace:
				gs.Player.Camera.Position = gs.Options.StartPosition
			default:
			}
		}

		tw.EventKeyGLFW(int(key), int(state)) // falsch, glfw3 hat scancodes
	}

	MouseMove := func(window *glfw.Window, mouseX, mouseY float64) {
		tw.EventMousePosGLFW(int(math.Floor(mouseX)), int(math.Floor(mouseY)))
	}

	CharacterType := func(window *glfw.Window, char uint) {
		tw.EventCharGLFW(int(char), int(glfw.Press))
	}

	window.SetCursorPositionCallback(MouseMove)
	window.SetKeyCallback(KeyPress)
	window.SetMouseButtonCallback(MouseButton)

	window.SetCharacterCallback(CharacterType)
}

func Input(gs *gamestate.GameState) {
	window := gs.Window
	inp := gamestate.PlayerInput{}

	if window.GetInputMode(glfw.Cursor) == glfw.CursorDisabled {
		mx64, my64 := window.GetCursorPosition()
		window.SetCursorPosition(0, 0)
		inp.Rotate[0] = -float32(my64) / 5
		inp.Rotate[1] = -float32(mx64) / 5
	}

	if window.GetKey(glfw.KeyE) == glfw.Press {
		inp.Move[2] -= 1
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		inp.Move[2] += 1
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		inp.Move[0] -= 1
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		inp.Move[0] += 1
	}
	if window.GetKey(glfw.KeyR) == glfw.Press {
		inp.Rotate[2] -= 1
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		inp.Rotate[2] += 1
	}

	gs.Player.Input = inp

	//updateLastMousePos()
}
