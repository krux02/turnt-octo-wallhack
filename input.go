package main

import glfw "github.com/go-gl/glfw3"
import "github.com/krux02/mathgl"
import "github.com/krux02/tw"
import "fmt"

var drag glfw.MouseButton = -1
var lastMousePos = mathgl.Vec2f{0, 0}
var highlight = 0

var currentMousePos func() mathgl.Vec2f
var updateLastMousePos func()

func InitInput(gamestate *GameState) {
	window := gamestate.Window

	MouseButton := func(window *glfw.Window, button glfw.MouseButton, state glfw.Action, modifiers glfw.ModifierKey) {
		if state == glfw.Press && drag == -1 {
			drag = button
		}
		if state == glfw.Release && drag == button {
			drag = -1
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
			default:
			}
		}

		tw.EventKeyGLFW(int(key), int(state)) // falsch, glfw3 hat scancodes
	}

	currentMousePos = func() mathgl.Vec2f {
		mx, my := window.GetCursorPosition()
		return mathgl.Vec2f{float32(mx), float32(my)}
	}

	updateLastMousePos = func() {
		lastMousePos = currentMousePos()
	}

	MouseMove := func(window *glfw.Window, mouseX, mouseY float64) {
		fmt.Println("mouse move", mouseX, mouseY)
		tw.EventMousePosGLFW(int(mouseX), int(mouseY))
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

func Input(gamestate *GameState) {

	window := gamestate.Window
	delta := currentMousePos().Sub(lastMousePos)
	inp := PlayerInput{}

	switch drag {
	case 0:
		if delta.Len() > 0 {
			inp.rotate[0] -= delta[1]
			inp.rotate[1] -= delta[0]
		}
	case 1:
		if delta.Len() > 0 {
			inp.rotate[1] -= delta[0]
			inp.rotate[2] -= delta[1]
		}
	}

	if window.GetKey(glfw.KeyE) == glfw.Press {
		inp.move[2] -= 1
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		inp.move[2] += 1
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		inp.move[0] -= 1
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		inp.move[0] += 1
	}
	if window.GetKey(glfw.KeyR) == glfw.Press {
		inp.rotate[2] -= 1
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		inp.rotate[2] += 1
	}

	gamestate.Player.SetInput(inp)
	updateLastMousePos()
}
