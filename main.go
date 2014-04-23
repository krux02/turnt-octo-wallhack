package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/krux02/turnt-octo-wallhack/debugContext"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/generation"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/tw"
	"runtime"
)

var counter = 1

func SdlError() {
	fmt.Println("errtest", counter)
	counter = counter + 1
	err := sdl.GetError()
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.LockOSThread()

	if sdl.Init(sdl.INIT_EVERYTHING) < 0 {
		panic("Unable to initialize SDL")
	}
	defer sdl.Quit()

	sdl.GL_SetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GL_SetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_DEBUG_FLAG|sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)

	window := sdl.CreateWindow("TOW", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1024, 768, sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if window == nil {
		panic("cant create window")
	}
	defer window.Destroy()

	//defer renderer.Destroy()

	glcontext := sdl.GL_CreateContext(window)
	if glcontext == nil {
		fmt.Sprintf("can't create context %v", sdl.GetError())
	}

	sdl.GL_MakeCurrent(window, glcontext)

	//err := gl.GlewInit()
	//fmt.Println(gl.GlewGetErrorString(err))
	//glew init
	gl.Init()

	defer sdl.GL_DeleteContext(glcontext)

	sdl.GL_SetSwapInterval(1)

	fmt.Println("glVersion", gl.GetString(gl.VERSION))

	tw.Init(tw.OPENGL_CORE, nil)
	defer tw.Terminate()

	gl.GetError() // Ignore error
	debugContext.InitDebugContext()

	world := generation.GenerateWorld(64, 64, 2)
	gs := gamestate.NewGameState(window, world)
	defer gs.Delete()
	worldRenderer := rendering.NewWorldRenderer(window, gs.World)
	defer worldRenderer.Delete()

	gs.Bar.AddButton("screen shot", worldRenderer.ScreenShot, "")

	MainLoop(gs, worldRenderer)
}
