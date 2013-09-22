package main

import (
	// "fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	mgl "github.com/krux02/mathgl"
	"github.com/krux02/turnt-octo-wallhack/world"
	"github.com/krux02/tw"
	"math"
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

		Proj := gamestate.Proj
		View := gamestate.Camera.View()
		Rot2D := gamestate.Camera.Rotation2D()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Disable(gl.BLEND)

		if gamestate.Options.Wireframe {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}
		if !gamestate.Options.DisableWorldRender {
			gamestate.WorldRenderer.HeightMapRenderer.Render(gamestate.World.HeightMap, Proj, View, currentTime, highlight)
		}
		if !gamestate.Options.DisableTreeRender {
			gamestate.WorldRenderer.PalmTrees.Render(Proj, View, Rot2D)
		}
		if !gamestate.Options.DisableParticlePhysics {
			gamestate.WorldRenderer.ParticleSystem.DoStep(currentTime)
		}
		if !gamestate.Options.DisableParticleRender {
			gamestate.WorldRenderer.ParticleSystem.Render(Proj, View)
		}

		rotation := mgl.HomogRotate3D(currentTime, mgl.Vec3f{0, 0, 1})
		portal := gamestate.World.Portals[0].Mesh

		boxVertices := portal.MakeBoxVertices()

		pv := Proj.Mul4(View)

		for _, portal := range gamestate.World.Portals {
			pos := portal.Position
			Model := mgl.Translate3D(float64(pos[0]), float64(pos[1]), float64(pos[2])).Mul4(rotation)
			gamestate.WorldRenderer.MeshRenderer.Render(&gamestate.WorldRenderer.Portal, Proj, View, Model)

			pvm := pv.Mul4(Model)

			meshMin := mgl.Vec4f{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
			meshMax := mgl.Vec4f{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}

			for _, v := range boxVertices {
				v = pvm.Mul4x1(v)
				v = v.Mul(1 / v[3])

				meshMin = world.Min(meshMin, v)
				meshMax = world.Max(meshMax, v)
			}

			if meshMin[0] < 1 && meshMin[1] < 1 && meshMin[2] < 1 &&
				meshMax[0] > -1 && meshMax[1] > -1 && meshMax[2] > -1 {
				w, h := gamestate.Window.GetSize()

				p1x, p1y := convertToPixelCoords(mgl.Vec2f{meshMin[0], meshMin[1]}, w, h)
				p2x, p2y := convertToPixelCoords(mgl.Vec2f{meshMax[0], meshMax[1]}, w, h)
				pw, ph := p2x-p1x, p2y-p1y

				if p1x != 0 || p1y != 0 || pw != w-1 || ph != h-1 {
					//gl.Viewport(p1x, p1y, pw, ph)
					gl.Enable(gl.SCISSOR_TEST)
					gl.Scissor(p1x, p1y, pw, ph)
					gl.ClearColor(1, 0, 0, 1)
					gl.Clear(gl.COLOR_BUFFER_BIT)
					gl.ClearColor(0, 0, 0, 1)
					//gl.Viewport(0,0,w,h)
					gl.Scissor(0, 0, w, h)
					gl.Disable(gl.SCISSOR_TEST)
				}
			}
		}

		tw.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func convertToPixelCoords(pos mgl.Vec2f, w, h int) (x, y int) {
	x = int(float32(w) * (pos[0] + 1) / 2)
	y = int(float32(h) * (pos[1] + 1) / 2)

	if x < 0 {
		x = 0
	}
	if x >= w {
		x = w - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= h {
		y = h - 1
	}
	return
}
