package main

import (
	"flag"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/krux02/libovr"
	"github.com/krux02/turnt-octo-wallhack/debugContext"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/generation"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/krux02/tw"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"runtime"
	"runtime/pprof"
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	runtime.LockOSThread()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	options := new(settings.BoolOptions)
	options.Load()
	defer options.Save()

	width := 1280
	height := 800
	window := sdl.CreateWindow("TOW", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if window == nil {
		panic("cant create window")
	}
	defer window.Destroy()

	//defer renderer.Destroy()

	glcontext := sdl.GL_CreateContext(window)
	if glcontext == nil {
		panic(fmt.Sprintf("can't create context %v", sdl.GetError()))
	}
	defer sdl.GL_DeleteContext(glcontext)
	sdl.GL_MakeCurrent(window, glcontext)

	//err := gl.GlewInit()
	//fmt.Println(gl.GlewGetErrorString(err))
	//glew init
	gl.Init()
	sdl.GL_SetSwapInterval(1)

	fmt.Println("glVersion", gl.GetString(gl.VERSION))

	tw.Init(tw.OPENGL_CORE, nil)
	defer tw.Terminate()

	ok := ovr.Initialize()
	if !ok {
		panic("cant't initialize ovr")
	}
	defer ovr.Shutdown()
	OvrTest()

	gl.GetError() // Ignore error
	debugContext.InitDebugContext()

	world := generation.GenerateWorld(64, 64, 2)
	gs := gamestate.NewGameState(options, window, world)
	defer gs.Delete()
	worldRenderer := rendering.NewWorldRenderer(window, gs.World)
	defer worldRenderer.Delete()

	gs.Bar.AddButton("screen shot", worldRenderer.ScreenShot, "")

	MainLoop(gs, worldRenderer)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			panic(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}

func OvrTest() {
	numDevices := ovr.HmdDetect()
	fmt.Printf("libovr found %d connected devices\n", numDevices)

	if numDevices > 0 {
		hmd := ovr.HmdCreate(0)
		defer hmd.Destroy()
		desc := hmd.GetDesc()
		fmt.Println("%+v", desc)

		fmt.Println(hmd.GetEnabledCaps())

	}
}
