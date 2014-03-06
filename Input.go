package main

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/tw"
	"math"
)

// import "fmt"

var drag glfw.MouseButton = -1
var lastMousePos = mgl.Vec2f{0, 0}
var highlight = 0

var currentMousePos func() mgl.Vec2f
var updateLastMousePos func()

func InitInput(gs *gamestate.GameState) {
	window := gs.Window

	MouseButton := func(window *glfw.Window, button glfw.MouseButton, state glfw.Action, modifiers glfw.ModifierKey) {
		if state == glfw.Press && drag == -1 {
			drag = button
		}
		if state == glfw.Release && drag == button {
			drag = -1
		}

		// ray cast testing
		if state == glfw.Press {

			mx, my := window.GetCursorPosition()
			_, H := window.GetSize()
			x := 2*float32(mx)/float32(H) - 1
			y := 2*float32(my)/float32(H) - 1

			dir_cs := mgl.Vec4f{x, y, -1, 0}
			m := gs.Camera.Model()
			pos_ws := gs.Camera.Position
			dir_ws := m.Mul4x1(dir_cs)

			pos_ws_v3 := mgl.Vec3f{pos_ws[0], pos_ws[1], pos_ws[2]}
			dir_ws_v3 := mgl.Vec3f{dir_ws[0], dir_ws[1], dir_ws[2]}

			fmt.Println(pos_ws_v3, dir_ws_v3)
			out, hit := gs.World.HeightMap.RayCast(pos_ws_v3, dir_ws_v3)
			fmt.Println(out, hit)
		}

		tw.EventMouseButtonGLFW(int(button), int(state))
	}

	MouseWheel := func(window *glfw.Window, xoffset, yoffset float64) {
		highlight += int(yoffset)
		tw.MouseWheel(int(yoffset)) // falsch glfw3 ist relativ
	}

	KeyPress := func(window *glfw.Window, key glfw.Key, _ int, state glfw.Action, modifiers glfw.ModifierKey) {
		if state == glfw.Press {
			switch key {
			case glfw.KeyKpAdd:
				highlight += 1
				// highlightLoc.Uniform1f(float32(highlight))
			case glfw.KeyKpSubtract:
				highlight -= 1
				// highlightLoc.Uniform1f(float32(highlight))
			case glfw.KeyEnter:
				switch window.GetInputMode(glfw.Cursor) {
				case glfw.CursorNormal:
					window.SetInputMode(glfw.Cursor, glfw.CursorDisabled)
				default:
					window.SetInputMode(glfw.Cursor, glfw.CursorNormal)
				}
			case glfw.KeySpace:
				gs.Player.Camera.Position = gs.Options.StartPosition
			default:
			}
		}

		tw.EventKeyGLFW(int(key), int(state)) // falsch, glfw3 hat scancodes
	}

	currentMousePos = func() mgl.Vec2f {
		mx, my := window.GetCursorPosition()
		return mgl.Vec2f{float32(mx), float32(my)}
	}

	updateLastMousePos = func() {
		lastMousePos = currentMousePos()
	}

	MouseMove := func(window *glfw.Window, mouseX, mouseY float64) {
		tw.EventMousePosGLFW(int(math.Floor(mouseX)), int(math.Floor(mouseY)))
	}

	CharacterType := func(window *glfw.Window, char uint) {
		tw.EventCharGLFW(int(char), int(glfw.Press))
	}

	window.SetCursorPositionCallback(MouseMove)
	window.SetScrollCallback(MouseWheel)
	window.SetKeyCallback(KeyPress)
	window.SetMouseButtonCallback(MouseButton)

	window.SetCharacterCallback(CharacterType)
}

func Input(gs *gamestate.GameState) {

	window := gs.Window
	delta := currentMousePos().Sub(lastMousePos)
	inp := gamestate.PlayerInput{}

	switch drag {
	case 1:
		if delta.Len() > 0 {
			inp.Rotate[0] -= delta[1]
			inp.Rotate[1] -= delta[0]
		}
	case 2:
		if delta.Len() > 0 {
			inp.Rotate[1] -= delta[0]
			inp.Rotate[2] -= delta[1]
		}
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
	updateLastMousePos()
}
