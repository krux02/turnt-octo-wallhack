package main

import (
	//	"fmt"
	"github.com/go-gl/gl"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/turnt-octo-wallhack/simulation"
	"github.com/krux02/tw"
)

func MainLoop(gs *gamestate.GameState, renderer *rendering.WorldRenderer) {
	var frames int
	time := float32(sdl.GetTicks()) / 1000
	window := gs.Window
	running := true

	for running {
		currentTime := float32(sdl.GetTicks()) / 1000
		if currentTime > time+1 {
			gs.Fps = float32(frames)
			frames = 0
			time = currentTime
		}
		frames += 1
		running = Input(gs, renderer)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Disable(gl.BLEND)

		simulation.Simulate(gs, renderer.ParticleSystem)
		renderer.Render(gs.World, &gs.Options, gs.Camera.View(), window)

		tw.Draw()
		sdl.GL_SwapWindow(window)

		helpers.UpdateManagers()
	}
}
