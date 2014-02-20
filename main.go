package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/debugContext"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/generation"
	"github.com/krux02/turnt-octo-wallhack/rendering"
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

	glfw.SwapInterval(60)

	window, err := glfw.CreateWindow(1024, 768, "Turnt Octo Wallhack", nil, nil)
	if window == nil {
		fmt.Println("error")
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	maj, min, rev := glfw.GetVersion()
	fmt.Printf("glVersion %d.%d.%d\n", maj, min, rev)

	window.MakeContextCurrent()

	tw.Init(tw.OPENGL_CORE, nil)
	defer tw.Terminate()

	gl.Init()
	gl.GetError() // Ignore error

	window.SetInputMode(glfw.StickyKeys, gl.TRUE)

	debugContext.InitDebugContext()

	world := generation.GenerateWorld(128, 128, 16)
	gs := gamestate.NewGameState(world, window)
	defer gs.Delete()
	renderer := rendering.NewWorldRenderer(gs.World)
	defer renderer.Delete()

	InitInput(gs)

	MainLoop(gs, renderer)
}
