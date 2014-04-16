package main

import (
	//	"fmt"
	//mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/turnt-octo-wallhack/simulation"
	"github.com/krux02/tw"
	//"math"
)

func MainLoop(gs *gamestate.GameState, renderer *rendering.WorldRenderer) {
	var frames int
	time := glfw.GetTime()

	window := gs.Window

	for !window.ShouldClose() && window.GetKey(glfw.KeyEscape) != glfw.Press {

		currentTime := glfw.GetTime()

		if currentTime > time+1 {
			gs.Fps = float32(frames)
			frames = 0
			time = currentTime
		}
		frames += 1

		Input(gs)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Disable(gl.BLEND)

		simulation.Simulate(gs, renderer.ParticleSystem)
		renderer.Render(gs.World, &gs.Options, gs.Camera.View(), window)

		tw.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
