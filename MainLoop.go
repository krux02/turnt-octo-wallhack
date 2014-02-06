package main

import (
	//	"fmt"
	//mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/tw"
	//"math"
)

func MainLoop(gamestate *GameState) {
	var frames int
	time := glfw.GetTime()

	window := gamestate.Window

	for !window.ShouldClose() && window.GetKey(glfw.KeyEscape) != glfw.Press {

		currentTime := glfw.GetTime()

		if currentTime > time+1 {
			gamestate.Fps = float32(frames)
			frames = 0
			time = currentTime
		}
		frames += 1

		Input(gamestate)

		gamestate.Player.Update(gamestate)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Disable(gl.BLEND)

		gamestate.WorldRenderer.Render(gamestate.World, &gamestate.Options, gamestate.Proj, gamestate.Camera.View(), window)

		tw.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
