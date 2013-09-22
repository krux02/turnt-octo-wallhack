package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/debugContext"
	"github.com/krux02/tw"
	"os"
	"runtime"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	runtime.LockOSThread()

	glfw.Init()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglDebugContext, gl.TRUE)

	window, err := glfw.CreateWindow(1024, 768, "Turnt Octo Wallhack", nil, nil)
	if window == nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	window.MakeContextCurrent()

	tw.Init(tw.OPENGL_CORE, nil)
	defer tw.Terminate()

	gl.Init()
	gl.GetError() // Ignore error

	window.SetInputMode(glfw.StickyKeys, gl.TRUE)

	debugContext.InitDebugContext()

	gamestate := NewGameState(window)
	defer gamestate.Delete()

	InitInput(gamestate)

	MainLoop(gamestate)
}
