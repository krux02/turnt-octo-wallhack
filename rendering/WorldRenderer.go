package rendering

import (
	// "fmt"
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
	PortalRenderer    *PortalRenderer
	Portal            PortalRenderData
	PalmTrees         *PalmTrees
	ParticleSystem    *particles.ParticleSystem
	Framebuffer       *FrameBuffer
	Textures          *Textures
	ScreenQuad        *ScreenQuadRenderer
}

func NewWorldRenderer(w *world.World) *WorldRenderer {
	textures := NewTextures()
	portalData := w.Portals[0].Mesh
	mr := NewMeshRenderer()
	pr := NewPortalRenderer()
	return &WorldRenderer{
		HeightMapRenderer: NewHeightMapRenderer(w.HeightMap),
		MeshRenderer:      mr,
		PortalRenderer:    pr,
		Portal:            pr.CreateRenderData(portalData),
		PalmTrees:         NewPalmTrees(w.HeightMap, 5000),
		ParticleSystem:    particles.NewParticleSystem(w, 10000, mgl.Vec3f{32, 32, 32}, 1, 250),
		Framebuffer:       NewFrameBuffer(),
		Textures:          textures,
		ScreenQuad:        NewScreenQuadRenderer(),
	}
}

func (this *WorldRenderer) Delete() {
	this.HeightMapRenderer.Delete()
	this.MeshRenderer.Delete()
	this.PortalRenderer.Delete()
	this.Portal.Indices.Delete()
	this.Portal.Vertices.Delete()
	this.Portal.VAO.Delete()
	this.PalmTrees.Delete()
	this.ParticleSystem.Delete()
	this.Framebuffer.Delete()
	this.Textures.Delete()
	this.ScreenQuad.Delete()
	*this = WorldRenderer{}
}

func (this *WorldRenderer) Render(ww *world.World, options *settings.BoolOptions, Proj mgl.Mat4f, View mgl.Mat4f, window *glfw.Window, max_recursion int, clippingPlane mgl.Vec4f) {

	this.Framebuffer.Bind()

	camera := NewCameraM(View)
	Rot2D := camera.Rotation2D()

	currentTime := glfw.GetTime()
	//rotation := mgl.HomogRotate3D(float32(currentTime), mgl.Vec3f{0, 0, 1})
	//rotation := options.Rotation.Mat4()

	W := ww.HeightMap.W
	H := ww.HeightMap.H

	if !options.NoParticlePhysics {
		this.ParticleSystem.DoStep(currentTime)
	}

	allPortals := make([]world.Portal, len(ww.Portals)*9)
	k := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			OffsetM := mgl.Translate3D(float32(i*W), float32(j*H), 0)
			OffsetV := mgl.Vec3f{float32(i * W), float32(j * H), 0}

			for _, portal := range ww.Portals {
				portal.Position = portal.Position.Add(OffsetV)
				allPortals[k] = portal
				k++
			}

			if options.Wireframe {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			}
			if !options.NoWorldRender {
				this.HeightMapRenderer.Render(Proj, View, OffsetM, currentTime, clippingPlane)
			}
			if !options.NoTreeRender {
				this.PalmTrees.Render(Proj, View.Mul4(OffsetM), Rot2D)
			}
			if !options.NoParticleRender {
				this.ParticleSystem.Render(Proj, View.Mul4(OffsetM))
			}
		}
	}

	boxVertices := ww.Portals[0].Mesh.MakeBoxVertices()
	pv := Proj.Mul4(View)

	// calculating nearest portal
	var dist float32 = math.MaxFloat32
	var pos4f = View.Inv().Mul4x1(mgl.Vec4f{0, 0, 0, 1})
	var pos3f = mgl.Vec3f{pos4f[0], pos4f[1], pos4f[2]}
	var nearestPortal world.Portal

	for _, portal := range allPortals {
		newDist := pos3f.Sub(portal.Position).Len()
		if newDist < dist {
			dist = newDist
			nearestPortal = portal
		}
	}

	// drawing  portals
	for _, portal := range allPortals {
		pos := portal.Position
		rotation := portal.Orientation.Mat4()
		Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)
		this.PortalRenderer.Render(&this.Portal, Proj, View, Model)

		if max_recursion > 0 && (portal == nearestPortal) {
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

					gl.ClearColor(0, 1, 1, 1)
					gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
					gl.ClearColor(0, 0, 0, 1)

					//gl.Viewport(0,0,w,h)
					// calculation View matrix that shows the target portal from the same angle as view shows the source portal
					pos2 := portal.Target.Position
					Model2 := mgl.Translate3D(pos2[0], pos2[1], pos2[2]).Mul4(portal.Target.Orientation.Mat4())
					View2 := View.Mul4(Model).Mul4(Model2.Inv())

					normal_os := mgl.Vec4f{0, 1, 0, 0}
					normal_ws := Model.Mul4x1(normal_os)
					// normal_cs := View.Mul4x1(normal_ws)
					view_dir := portal.Position.Sub(camera.Position)
					sign := view_dir.Dot(mgl.Vec3f{normal_ws[0], normal_ws[1], normal_ws[2]})

					if sign > 0 {
						clippingPlane = Model2.Mul4x1(mgl.Vec4f{0, 1, 0, 0})
					} else {
						clippingPlane = Model2.Mul4x1(mgl.Vec4f{0, -1, 0, 0})
					}
					clippingPlane[3] = -clippingPlane.Dot(mgl.Vec4f{pos2[0], pos2[1], pos2[2], 0})

					//camera2 := NewCameraM(View2)

					this.Render(ww, options, Proj, View2, window, max_recursion-1, clippingPlane)

					gl.Scissor(0, 0, w, h)
					gl.Disable(gl.SCISSOR_TEST)
				}
			}
		}
	}

	this.Framebuffer.Unbind()

	this.ScreenQuad.Render()

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
