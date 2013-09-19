package main

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/tw"
	"fmt"
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
			fmt.Printf("\r%f       ", gamestate.Fps)
		}
		frames += 1

		Input(gamestate)

		gamestate.Player.Update(gamestate)

		Proj := gamestate.Proj
		View := gamestate.Camera.View()
		// ProjView := Proj.Mul4(View)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Disable(gl.BLEND)

		if gamestate.Options.Wireframe {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		if !gamestate.Options.DisableWorldRender {
			gamestate.WordlRenderer.Render(gamestate.World.HeightMap, Proj, View, currentTime, highlight)
		}
		if !gamestate.Options.DisableTreeRender {
			gamestate.PalmTrees.Render(Proj, View)
		}
		if !gamestate.Options.DisableParticlePhysics {
			gamestate.ParticleSystem.DoStep(currentTime)
		}
		if !gamestate.Options.DisableParticleRender {
			gamestate.ParticleSystem.Render(Proj, View)
		}

		tw.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
