package rendering

import (
	//	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type WorldRenderer struct {
	Proj              mgl.Mat4f
	Textures          *Textures
	HeightMapRenderer *HeightMapRenderer
	WaterRenderer     *WaterRenderer
	MeshRenderer      *MeshRenderer
	PortalRenderer    *PortalRenderer
	Portal            PortalRenderData
	PalmRenderer      *PalmRenderer
	ParticleSystem    *particles.ParticleSystem
	SkyboxRenderer    *SkyboxRenderer
	Framebuffer       [2]*FrameBuffer
	ScreenQuad        *ScreenQuadRenderer
	DebugRenderer     *LineRenderer
	MaxRecursion      int
	screenShot        bool
}

func (this *WorldRenderer) Resize(width, height int) {
	this.Proj = mgl.Perspective(90, float32(width)/float32(height), 0.3, 1000)
	gl.Viewport(0, 0, width, height)
	for _, fb := range this.Framebuffer {
		fb.Resize(width, height)
	}
}

func (this *WorldRenderer) ScreenShot() {
	this.screenShot = true
}

func NewWorldRenderer(window *sdl.Window, w *gamestate.World) *WorldRenderer {
	width, height := window.GetSize()

	mr := NewMeshRenderer()
	return &WorldRenderer{
		Proj:              mgl.Perspective(90, float32(width)/float32(height), 0.3, 1000),
		Textures:          NewTextures(w.HeightMap),
		HeightMapRenderer: NewHeightMapRenderer(w.HeightMap),
		WaterRenderer:     NewWaterRenderer(w.HeightMap),
		MeshRenderer:      mr,
		PortalRenderer:    NewPortalRenderer(),
		PalmRenderer:      NewPalmRenderer(&w.Palms),
		ParticleSystem:    particles.NewParticleSystem(w, 10000, mgl.Vec3f{32, 32, 32}, 1, 250),
		SkyboxRenderer:    NewSkyboxRenderer(),
		Framebuffer:       [2]*FrameBuffer{NewFrameBuffer(window.GetSize()), NewFrameBuffer(window.GetSize())},
		ScreenQuad:        NewScreenQuadRenderer(),
		DebugRenderer:     NewLineRenderer(),
		MaxRecursion:      1,
	}
}

func (this *WorldRenderer) Delete() {
	this.Textures.Delete()
	this.HeightMapRenderer.Delete()
	this.MeshRenderer.Delete()
	this.PortalRenderer.Delete()
	// TODO delete portal data
	this.PalmRenderer.Delete()
	this.ParticleSystem.Delete()
	this.SkyboxRenderer.Delete()
	this.WaterRenderer.Delete()
	for _, Framebuffer := range this.Framebuffer {
		Framebuffer.Delete()
	}
	this.ScreenQuad.Delete()
	this.DebugRenderer.Delete()
	*this = WorldRenderer{}
}

func (this *WorldRenderer) Render(ww *gamestate.World, options *settings.BoolOptions, View mgl.Mat4f, window *sdl.Window) {
	this.render(ww, options, View, window, 0, mgl.Vec4f{3 / 5.0, 4 / 5.0, 0, math.MaxFloat32}, nil)

	gl.ActiveTexture(gl.TEXTURE0)
	this.Framebuffer[0].RenderTexture.Bind(gl.TEXTURE_RECTANGLE)
	this.ScreenQuad.Render(0)

	if this.screenShot {
		this.screenShot = false
		helpers.SaveTexture(gl.TEXTURE_RECTANGLE, 0, "screenshot.png")
	}
}

func (this *WorldRenderer) render(ww *gamestate.World, options *settings.BoolOptions, View mgl.Mat4f, window *sdl.Window, recursion int, clippingPlane mgl.Vec4f, srcPortal *gamestate.Portal) {

	this.Framebuffer[recursion].Bind()
	defer this.Framebuffer[recursion].Unbind()

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	camera := gamestate.NewCameraFromMat4(View)
	Rot2D := camera.Rotation2D()

	gl.CullFace(gl.BACK)

	currentTime := float64(sdl.GetTicks()) / 1000

	if options.Wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	if options.Skybox {
		gl.Disable(gl.DEPTH_TEST)
		this.SkyboxRenderer.Render(this.Proj, View, 7)
		gl.Enable(gl.DEPTH_TEST)
	}

	gl.Enable(gl.CULL_FACE)

	for _, entity := range ww.ExampleObjects {
		this.RenderEntity(View, entity)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	if recursion != 0 {
		gl.Enable(gl.CLIP_DISTANCE0)
		defer gl.Disable(gl.CLIP_DISTANCE0)
	}

	if options.WorldRender {
		this.HeightMapRenderer.Render(this.Proj, View, mgl.Ident4f(), ww.HeightMap, clippingPlane)
	}
	if options.WaterRender {
		this.WaterRenderer.Render(this.Proj, View, mgl.Ident4f(), currentTime, clippingPlane, options.WaterNormals)
	}

	gl.Disable(gl.CULL_FACE)

	gl.Disable(gl.BLEND)
	if options.TreeRender {
		this.PalmRenderer.Render(this.Proj, View, Rot2D, clippingPlane)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	if options.ParticleRender {
		this.ParticleSystem.Render(this.Proj, View, clippingPlane)
	}

	gl.Disable(gl.BLEND)

	boxVertices := gamestate.QuadMesh().MakeBoxVertices()

	pv := this.Proj.Mul4(View)

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
			this.PortalRenderer.Render(this.Proj, View, Model, clippingPlane, 7)
		}
	}

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	if options.DebugLines {

		if options.DepthTestDebugLines {
			gl.Disable(gl.DEPTH_TEST)
		}
		this.DebugRenderer.Render(this.Proj, View)
		gl.Enable(gl.DEPTH_TEST)
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

			this.render(ww, options, View2, window, recursion+1, clippingPlane, nearestPortal)

			gl.ActiveTexture(gl.TEXTURE0)
			this.Framebuffer[recursion+1].RenderTexture.Bind(gl.TEXTURE_RECTANGLE)

			if scissor {
				gl.Scissor(0, 0, w, h)
				gl.Disable(gl.SCISSOR_TEST)
			}

			this.Framebuffer[recursion].Bind()
			pos := nearestPortal.Position
			rotation := nearestPortal.Orientation.Mat4()
			Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)
			this.PortalRenderer.Render(this.Proj, View, Model, clippingPlane, 0)

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
