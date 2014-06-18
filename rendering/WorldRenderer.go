package rendering

import (
	//	"fmt"
	"github.com/go-gl/gl"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	//	"github.com/krux02/turnt-octo-wallhack/math32"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type WorldRenderer struct {
	Proj, View         mgl.Mat4
	ClippingPlane_ws   mgl.Vec4
	Textures           *Textures
	HeightMapRenderer  *renderstuff.Renderer
	WaterRenderer      *renderstuff.Renderer
	WaterRendererA     *renderstuff.Renderer
	WaterRendererB     *renderstuff.Renderer
	MeshRenderer       *renderstuff.Renderer
	PortalRenderer     *renderstuff.Renderer
	TreeRenderer       *renderstuff.Renderer
	Skybox             *Skybox
	SkyboxRenderer     *renderstuff.Renderer
	ParticleSystem     *particles.ParticleSystem
	Framebuffer        [2]*FrameBuffer
	ScreenQuad         *ScreenQuad
	ScreenQuadRenderer *renderstuff.Renderer
	DebugRenderer      *LineRenderer
	MaxRecursion       int
	screenShot         bool
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

	return &WorldRenderer{
		Proj:               mgl.Perspective(90, float32(width)/float32(height), 0.3, 1000),
		View:               mgl.Ident4(),
		ClippingPlane_ws:   mgl.Vec4{1, 0, 0, -1000000},
		Textures:           NewTextures(w.HeightMap),
		HeightMapRenderer:  NewHeightMapRenderer(),
		WaterRendererA:     NewSurfaceWaterRenderer(),
		WaterRendererB:     NewDebugWaterRenderer(),
		MeshRenderer:       NewMeshRenderer(),
		PortalRenderer:     NewPortalRenderer(),
		TreeRenderer:       NewTreeRenderer(),
		SkyboxRenderer:     NewSkyboxRenderer(),
		Skybox:             &Skybox{},
		ParticleSystem:     particles.NewParticleSystem(w, 10000, mgl.Vec3{32, 32, 32}, 1, 250),
		Framebuffer:        [2]*FrameBuffer{NewFrameBuffer(window.GetSize()), NewFrameBuffer(window.GetSize())},
		ScreenQuad:         &ScreenQuad{},
		ScreenQuadRenderer: NewScreenQuadRenderer(),
		DebugRenderer:      NewLineRenderer(),
		MaxRecursion:       1,
	}
}

func (this *WorldRenderer) Delete() {
	this.Textures.Delete()
	this.HeightMapRenderer.Delete()
	this.MeshRenderer.Delete()
	this.PortalRenderer.Delete()
	this.TreeRenderer.Delete()
	this.ParticleSystem.Delete()
	this.SkyboxRenderer.Delete()
	this.WaterRendererA.Delete()
	this.WaterRendererB.Delete()
	for _, Framebuffer := range this.Framebuffer {
		Framebuffer.Delete()
	}
	this.ScreenQuadRenderer.Delete()
	this.DebugRenderer.Delete()
	*this = WorldRenderer{}
}

func (this *WorldRenderer) Render(ww *gamestate.World, options *settings.BoolOptions, window *sdl.Window) {
	camera := ww.Player.Camera
	p0 := camera.Pos4f()

	w, h := window.GetSize()

	camera.MoveRelative(mgl.Vec4{-0.1, 0, 0, 0})
	this.View = (ww.PortalTransform(p0, camera.Pos4f()).Mul4(camera.Model())).Inv()
	viewport := Viewport{0, 0, w / 2, h}
	viewport.Activate()
	this.render(ww, options, viewport, 0, nil)

	camera.MoveRelative(mgl.Vec4{+0.2, 0, 0, 0})
	this.View = (ww.PortalTransform(p0, camera.Pos4f()).Mul4(camera.Model())).Inv()
	viewport.X = w / 2
	viewport.Activate()
	this.render(ww, options, viewport, 0, nil)

	gl.Viewport(0, 0, w, h)

	gl.ActiveTexture(gl.TEXTURE0)
	this.Framebuffer[0].RenderTexture.Bind(gl.TEXTURE_RECTANGLE)
	this.ScreenQuadRenderer.Render(this.ScreenQuad, this.Proj, this.View, this.ClippingPlane_ws, nil)

	if this.screenShot {
		this.screenShot = false
		helpers.SaveTexture(gl.TEXTURE_RECTANGLE, 0, "screenshot.png")
	}
}

type Viewport struct {
	X, Y, W, H int
}

func (this *Viewport) Activate() {
	gl.Viewport(this.X, this.Y, this.W, this.H)
}

func (this *Viewport) ToPixel(pos mgl.Vec2) (X, Y int) {
	x := int(float32(this.W) * (pos[0] + 1) / 2)
	y := int(float32(this.H) * (pos[1] + 1) / 2)

	if x < 0 {
		x = 0
	}
	if x >= this.W {
		x = this.W - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= this.H {
		y = this.H - 1
	}
	return x + this.X, y + this.Y
}

func (this *WorldRenderer) render(ww *gamestate.World, options *settings.BoolOptions, viewport Viewport, recursion int, srcPortal *gamestate.Portal) {

	this.Framebuffer[recursion].Bind()
	defer this.Framebuffer[recursion].Unbind()

	gl.Clear(gl.DEPTH_BUFFER_BIT)

	camera := gamestate.NewCameraFromMat4(this.View)
	Rot2D := camera.Rotation2D()

	gl.CullFace(gl.BACK)

	time := float64(sdl.GetTicks()) / 1000

	if options.Wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	if options.Skybox {
		gl.Disable(gl.DEPTH_TEST)
		this.SkyboxRenderer.Render(this.Skybox, this.Proj, this.View, this.ClippingPlane_ws, nil)
		gl.Enable(gl.DEPTH_TEST)
	}

	gl.Enable(gl.CULL_FACE)

	if recursion != 0 {
		gl.Enable(gl.CLIP_DISTANCE0)
		defer gl.Disable(gl.CLIP_DISTANCE0)
	}

	for _, entity := range ww.ExampleObjects {
		this.MeshRenderer.Render(entity, this.Proj, this.View, this.ClippingPlane_ws, nil)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Disable(gl.CULL_FACE)

	if options.WorldRender {
		this.HeightMapRenderer.Render(ww.HeightMap, this.Proj, this.View, this.ClippingPlane_ws, nil)
	}
	PlayerPos := ww.Player.Position()
	ww.Water.Height = PlayerPos[2] - 15
	if options.WaterRender {
		this.WaterRendererA.Render(ww.Water, this.Proj, this.View, this.ClippingPlane_ws, WaterRenderUniforms{time, PlayerPos})
	}
	if options.WaterNormals {
		this.WaterRendererB.Render(ww.Water, this.Proj, this.View, this.ClippingPlane_ws, WaterRenderUniforms{time, PlayerPos})
	}

	gl.Disable(gl.CULL_FACE)

	gl.Disable(gl.BLEND)
	if options.TreeRender {
		this.TreeRenderer.Render(ww.Trees, this.Proj, this.View, this.ClippingPlane_ws, Rot2D)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	if options.ParticleRender {
		this.ParticleSystem.Render(this.Proj, this.View, this.ClippingPlane_ws)
	}

	gl.Disable(gl.BLEND)

	boxVertices := (*gamestate.TriangleMesh)(gamestate.QuadMesh()).MakeBoxVertices()

	pv := this.Proj.Mul4(this.View)

	// calculating nearest portal
	pos4f := this.View.Inv().Mul4x1(mgl.Vec4{0, 0, 0, 1})
	nearestPortal := ww.NearestPortal(pos4f)

	// draw  all portals except the nearest and the portal that we are looking throug
	for _, portal := range ww.Portals {
		// do not draw the nearest portal or the portal behind the source portal if available
		if (nearestPortal != portal) && (srcPortal == nil || srcPortal.Target != portal) {
			gl.Enable(gl.DEPTH_CLAMP)
			additionalUniforms := map[string]int{"Image": 7}
			this.PortalRenderer.Render(portal, this.Proj, this.View, this.ClippingPlane_ws, additionalUniforms)
		}
	}

	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	if options.DebugLines {

		if options.DepthTestDebugLines {
			gl.Disable(gl.DEPTH_TEST)
		}
		this.DebugRenderer.Render(this.Proj, this.View)
		gl.Enable(gl.DEPTH_TEST)
	}

	// draw
	if recursion < this.MaxRecursion {
		portal := nearestPortal
		pos := portal.Position
		rotation := portal.Orientation.Mat4()
		Model := mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotation)

		pvm := pv.Mul4(Model)
		meshMin := mgl.Vec4{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
		meshMax := mgl.Vec4{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
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

			p1x, p1y := viewport.ToPixel(meshMin.Vec2())
			p2x, p2y := viewport.ToPixel(meshMax.Vec2())
			pw, ph := p2x-p1x, p2y-p1y

			// do scissoring only when all vertices are in front of the camera
			scissor := meshMax[2] < 1
			scissor = scissor && (p1x != 0 || p1y != 0 || pw != viewport.W-1 || ph != viewport.H-1)

			if scissor {
				gl.Enable(gl.SCISSOR_TEST)
				gl.Scissor(p1x, p1y, pw, ph)
			}

			// omit rendering when portal is not in frustum at all
			// calculation View matrix that shows the target portal from the same angle as view shows the source portal

			//pos2 := portal.Target.Position
			Model2 := portal.Target.Model()
			// model matrix, so that portal 1 in camera 1 looks identical to portal 2 in camera
			oldView := this.View
			this.View = this.View.Mul4(Model).Mul4(Model2.Inv())

			normal_os := portal.Target.Normal
			normal_ws := Model.Mul4x1(normal_os)
			view_dir := helpers.HomogenDiff(portal.Position, camera.Position)
			sign := view_dir.Dot(normal_ws)

			oldClippingPlane := this.ClippingPlane_ws
			this.ClippingPlane_ws = portal.Target.ClippingPlane(sign > 0)

			this.render(ww, options, viewport, recursion+1, nearestPortal)
			this.ClippingPlane_ws = oldClippingPlane
			this.View = oldView

			gl.ActiveTexture(gl.TEXTURE0)
			this.Framebuffer[recursion+1].RenderTexture.Bind(gl.TEXTURE_RECTANGLE)

			if scissor {
				//gl.Scissor(0, 0, w, h)
				gl.Disable(gl.SCISSOR_TEST)
			}
			this.Framebuffer[recursion].Bind()
			gl.Enable(gl.DEPTH_CLAMP)
			additionalUniforms := map[string]int{"Image": 0}
			this.PortalRenderer.Render(nearestPortal, this.Proj, this.View, this.ClippingPlane_ws, additionalUniforms)
		}
	}
}
