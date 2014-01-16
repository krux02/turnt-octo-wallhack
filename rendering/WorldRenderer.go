package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/krux02/turnt-octo-wallhack/world"
	"math"
)

type WorldRenderer struct {
	HeightMapRenderer *HeightMapRenderer
	MeshRenderer      *MeshRenderer
	Portal            MeshRenderData
	PalmTrees         *PalmTrees
	ParticleSystem    *particles.ParticleSystem
}

func NewWorldRenderer(w *world.World) *WorldRenderer {
	portalData := w.Portals[0].Mesh
	mr := NewMeshRenderer()
	return &WorldRenderer{
		HeightMapRenderer: NewHeightMapRenderer(w.HeightMap),
		MeshRenderer:      mr,
		Portal:            mr.CreateMeshRenderData(portalData),
		PalmTrees:         NewPalmTrees(w.HeightMap, 5000),
		ParticleSystem:    particles.NewParticleSystem(w, 10000, mgl.Vec3f{32, 32, 32}, 1, 250),
	}
}

func (this *WorldRenderer) Render(w *world.World, options *settings.BoolOptions, Proj mgl.Mat4f, View mgl.Mat4f, Rot2D mgl.Mat3f, window *glfw.Window) {
	currentTime := glfw.GetTime()
	rotation := mgl.HomogRotate3D(float32(currentTime), mgl.Vec3f{0, 0, 1})

	W := w.HeightMap.W
	H := w.HeightMap.H

	if !options.DisableParticlePhysics {
		this.ParticleSystem.DoStep(currentTime)
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			Offset := mgl.Translate3D(float32(i*W), float32(j*H), 0)

			if options.Wireframe {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			}
			if !options.DisableWorldRender {
				this.HeightMapRenderer.Render(Proj, View, Offset, currentTime)
			}
			if !options.DisableTreeRender {
				this.PalmTrees.Render(Proj, View.Mul4(Offset), Rot2D)
			}
			if !options.DisableParticleRender {
				this.ParticleSystem.Render(Proj, View.Mul4(Offset))
			}

			portal := w.Portals[0].Mesh
			boxVertices := portal.MakeBoxVertices()
			pv := Proj.Mul4(View)

			var dist float32 = 1e38
			var pos4f = View.Inv().Mul4x1(mgl.Vec4f{0, 0, 0, 1})
			var pos3f = mgl.Vec3f{pos4f[0], pos4f[1], pos4f[2]}
			var nearestPortal world.Portal

			for _, portal := range w.Portals {
				newDist := pos3f.Sub(portal.Position).Len()
				if newDist < dist {
					dist = newDist
					nearestPortal = portal
				}
			}

			for _, portal := range w.Portals {
				pos := portal.Position
				Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)
				this.MeshRenderer.Render(&this.Portal, Proj, View.Mul4(Offset), Model)

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
					w, h := window.GetSize()
					p1x, p1y := convertToPixelCoords(mgl.Vec2f{meshMin[0], meshMin[1]}, w, h)
					p2x, p2y := convertToPixelCoords(mgl.Vec2f{meshMax[0], meshMax[1]}, w, h)
					pw, ph := p2x-p1x, p2y-p1y
					if p1x != 0 || p1y != 0 || pw != w-1 || ph != h-1 {
						//gl.Viewport(p1x, p1y, pw, ph)
						gl.Enable(gl.SCISSOR_TEST)
						gl.Scissor(p1x, p1y, pw, ph)
						if portal == nearestPortal {
							gl.ClearColor(0, 1, 1, 1)
						} else {
							gl.ClearColor(1, 0, 0, 1)
						}
						gl.Clear(gl.COLOR_BUFFER_BIT)
						gl.ClearColor(0, 0, 0, 1)
						//gl.Viewport(0,0,w,h)
						gl.Scissor(0, 0, w, h)
						gl.Disable(gl.SCISSOR_TEST)
					}
				}
			}
		}
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
