package main

import "github.com/go-gl/glfw"
import "github.com/krux02/mathgl"
import "github.com/krux02/tw"
import "fmt"

var drag = -1
var lastMousePos = mathgl.Vec2f{0, 0}
var highlight = 0

func MouseButton(button int, state int) {
	if state == glfw.KeyPress && drag == -1 {
		drag = button
	}
	if state == glfw.KeyRelease && drag == button {
		drag = -1
	}

	tw.EventMouseButtonGLFW(button, state)
}

func MouseWheel(pos int) {
	highlight = pos

	tw.EventMouseWheelGLFW(pos)
}

var CursorEnabled = true

func KeyPress(key int, state int) {
	if state == glfw.KeyPress {
		switch key {
		case glfw.KeyKPAdd:
			highlight += 1
			// highlightLoc.Uniform1f(float32(highlight))
		case glfw.KeyKPSubtract:
			highlight -= 1
			// highlightLoc.Uniform1f(float32(highlight))
		case glfw.KeyEnter:
			if CursorEnabled {
				glfw.Disable(glfw.MouseCursor)
				CursorEnabled = false
			} else {
				glfw.Enable(glfw.MouseCursor)
				CursorEnabled = true
			}
		default:
		}
	}

	tw.EventKeyGLFW(key, state)
}

func currentMousePos() mathgl.Vec2f {
	mx, my := glfw.MousePos()
	return mathgl.Vec2f{float32(mx), float32(my)}
}

func updateLastMousePos() {
	lastMousePos = currentMousePos()
}

func MouseMove(mouseX, mouseY int) {
	fmt.Println("mouse move", mouseX, mouseY)
	tw.EventMousePosGLFW(mouseX, mouseY)
}

func InitInput(gamestate *GameState) {
	glfw.SetMouseWheelCallback(MouseWheel)
	glfw.SetKeyCallback(KeyPress)
	glfw.SetMouseButtonCallback(MouseButton)
	glfw.SetMousePosCallback(MouseMove)
	glfw.SetCharCallback(tw.EventCharGLFW)
}

func Input(gamestate *GameState) {

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

	if glfw.Key('L') == glfw.KeyPress {
		inp.move[2] -= 1
	}
	if glfw.Key('A') == glfw.KeyPress {
		inp.move[2] += 1
	}
	if glfw.Key('I') == glfw.KeyPress {
		inp.move[0] -= 1
	}
	if glfw.Key('E') == glfw.KeyPress {
		inp.move[0] += 1
	}
	if glfw.Key('C') == glfw.KeyPress {
		inp.rotate[2] -= 1
	}
	if glfw.Key('V') == glfw.KeyPress {
		inp.rotate[2] += 1
	}

	gamestate.Player.SetInput(inp)
	updateLastMousePos()
}
