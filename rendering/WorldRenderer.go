package rendering

import (
	//	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"math"
	//"math/rand"
)

type WorldRenderer struct {
	Textures          *Textures
	HeightMapRenderer *HeightMapRenderer
	WaterRenderer     *WaterRenderer
	MeshRenderer      *MeshRenderer
	PortalRenderer    *PortalRenderer
	Portal            PortalRenderData
	PalmTrees         *PalmTrees
	ParticleSystem    *particles.ParticleSystem
	Framebuffer       [2]*FrameBuffer
	ScreenQuad        *ScreenQuadRenderer
	MaxRecursion      int
}

func NewWorldRenderer(w *gamestate.World) *WorldRenderer {

	portalData := w.Portals[0].Mesh
	mr := NewMeshRenderer()
	pr := NewPortalRenderer()
	return &WorldRenderer{
		Textures:          NewTextures(w.HeightMap),
		HeightMapRenderer: NewHeightMapRenderer(w.HeightMap),
		WaterRenderer:     NewWaterRenderer(w.HeightMap),
		MeshRenderer:      mr,
		PortalRenderer:    pr,
		Portal:            pr.CreateRenderData(portalData),
		PalmTrees:         NewPalmTrees(w.HeightMap, 5000),
		ParticleSystem:    particles.NewParticleSystem(w, 10000, mgl.Vec3f{32, 32, 32}, 1, 250),
		Framebuffer:       [2]*FrameBuffer{NewFrameBuffer(), NewFrameBuffer()},
		ScreenQuad:        NewScreenQuadRenderer(),
		MaxRecursion:      1,
	}
}

func (this *WorldRenderer) Delete() {
	this.Textures.Delete()
	this.HeightMapRenderer.Delete()
	this.MeshRenderer.Delete()
	this.PortalRenderer.Delete()
	// TODO delete portal data
	this.PalmTrees.Delete()
	this.ParticleSystem.Delete()
	for _, Framebuffer := range this.Framebuffer {
		Framebuffer.Delete()
	}
	this.ScreenQuad.Delete()
	*this = WorldRenderer{}
}

func (this *WorldRenderer) Render(ww *gamestate.World, options *settings.BoolOptions, Proj mgl.Mat4f, View mgl.Mat4f, window *glfw.Window) {
	this.render(ww, options, Proj, View, window, 0, mgl.Vec4f{3 / 5.0, 4 / 5.0, 0, math.MaxFloat32}, nil)

	gl.ActiveTexture(gl.TEXTURE9)
	this.Framebuffer[0].RenderTexture.Bind(gl.TEXTURE_RECTANGLE)
	this.ScreenQuad.Render()
}

func (this *WorldRenderer) render(ww *gamestate.World, options *settings.BoolOptions, Proj mgl.Mat4f, View mgl.Mat4f, window *glfw.Window, recursion int, clippingPlane mgl.Vec4f, srcPortal *gamestate.Portal) {
	this.Framebuffer[recursion].Bind()
	defer this.Framebuffer[recursion].Unbind()

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	camera := gamestate.NewCameraFromMat4(View)
	Rot2D := camera.Rotation2D()

	gl.CullFace(gl.BACK)

	currentTime := glfw.GetTime()

	if !options.NoParticlePhysics {
		this.ParticleSystem.DoStep(currentTime)
	}

	if options.Wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.CLIP_DISTANCE0)

	if !options.NoWorldRender {
		this.HeightMapRenderer.Render(Proj, View, mgl.Ident4f(), clippingPlane)
		this.WaterRenderer.Render(Proj, View, mgl.Ident4f(), currentTime, clippingPlane)
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	if !options.NoParticleRender {
		this.ParticleSystem.Render(Proj, View, clippingPlane)
	}

	gl.Disable(gl.CULL_FACE)

	if !options.NoTreeRender {
		this.PalmTrees.Render(Proj, View, Rot2D, clippingPlane)
	}

	boxVertices := ww.Portals[0].Mesh.MakeBoxVertices()
	pv := Proj.Mul4(View)

	// calculating nearest portal
	pos4f := View.Inv().Mul4x1(mgl.Vec4f{0, 0, 0, 1})
	nearestPortal := ww.NearestPortal(pos4f)

	// draw  all portals except the nearest and the portal that we are looking throug
	for _, portal := range ww.Portals {
		// do not draw the nearest portal or the portal behind the source portal if available
		if (nearestPortal != portal) && (srcPortal == nil || srcPortal.Target != portal) {
			pos := portal.Position
			rotation := portal.Orientation.Mat4()
			Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)
			this.PortalRenderer.Render(&this.Portal, Proj, View, Model, 7)
		}
	}

	// draw
	if recursion < this.MaxRecursion {
		portal := nearestPortal
		pos := portal.Position
		rotation := portal.Orientation.Mat4()
		Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)

		pvm := pv.Mul4(Model)
		meshMin := mgl.Vec4f{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
		meshMax := mgl.Vec4f{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
		for _, v := range boxVertices {
			v = pvm.Mul4x1(v)
			v = v.Mul(1 / v[3])
			meshMin = gamestate.Min(meshMin, v)
			meshMax = gamestate.Max(meshMax, v)
		}

		// at least partially visible
		if -1 < meshMax[0] && meshMin[0] < 1 &&
			-1 < meshMax[1] && meshMin[1] < 1 &&
			-1 < meshMax[2] && meshMin[2] < 1 {

			w, h := window.GetSize()
			p1x, p1y := convertToPixelCoords(mgl.Vec2f{meshMin[0], meshMin[1]}, w, h)
			p2x, p2y := convertToPixelCoords(mgl.Vec2f{meshMax[0], meshMax[1]}, w, h)
			pw, ph := p2x-p1x, p2y-p1y

			// do scissoring only when all vertices are in front of the camera
			scissor := meshMax[2] < 1
			scissor = scissor && (p1x != 0 || p1y != 0 || pw != w-1 || ph != h-1)

			if scissor {
				gl.Enable(gl.SCISSOR_TEST)
				gl.Scissor(p1x, p1y, pw, ph)
			}

			// omit rendering when portal is not in frustum at all
			// calculation View matrix that shows the target portal from the same angle as view shows the source portal

			//pos2 := portal.Target.Position
			Model2 := portal.Target.Model()
			// model matrix, so that portal 1 in camera 1 looks identical to portal 2 in camera
			View2 := View.Mul4(Model).Mul4(Model2.Inv())

			normal_os := portal.Target.Normal
			normal_ws := Model.Mul4x1(normal_os)
			view_dir := helpers.HomogenDiff(portal.Position, camera.Position)
			sign := view_dir.Dot(normal_ws)

			clippingPlane = portal.Target.ClippingPlane(sign > 0)

			this.render(ww, options, Proj, View2, window, recursion+1, clippingPlane, nearestPortal)
			gl.ActiveTexture(gl.TEXTURE8)
			this.Framebuffer[recursion+1].RenderTexture.Bind(gl.TEXTURE_RECTANGLE)

			if scissor {
				gl.Scissor(0, 0, w, h)
				gl.Disable(gl.SCISSOR_TEST)
			}

			this.Framebuffer[recursion].Bind()
			pos := nearestPortal.Position
			rotation := nearestPortal.Orientation.Mat4()
			Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)
			this.PortalRenderer.Render(&this.Portal, Proj, View, Model, 8)

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
