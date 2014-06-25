package rendering

import (
	"fmt"
	"github.com/go-gl/gl"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	//"github.com/krux02/turnt-octo-wallhack/math32"
	"github.com/krux02/libovr"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/veandco/go-sdl2/sdl"
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
	OvrStuff           *OvrStuff
	FrameIndex         int
	MaxRecursion       int
	width, height      int
	riftRender         bool
	screenShot         bool
}

func (this *WorldRenderer) Resize(width, height int) {
	if !this.riftRender {
		this.Proj = mgl.Perspective(90, float32(width)/float32(height), 0.3, 1000)
		this.width = width
		this.height = height
		for _, fb := range this.Framebuffer {
			fb.Resize(width, height)
		}
	}
}

func (this *WorldRenderer) ScreenShot() {
	this.screenShot = true
}

func (this *WorldRenderer) ToggleRift() {
	this.riftRender = !this.riftRender

	if this.riftRender {
		for _, fb := range this.Framebuffer {
			fb.Resize(1920, 1080)
		}
		this.width, this.height = 1920, 1080
	} else {
		this.Resize(1280, 800)
	}
}

func (this *WorldRenderer) RiftRender() bool {
	return this.riftRender
}

func NewWorldRenderer(window *sdl.Window, w *gamestate.World) *WorldRenderer {

	width, height := window.GetSize()

	//framebufferWidth, FrameBufferHeight := 1920, 1080
	framebuffers := [2]*FrameBuffer{
		NewFrameBuffer(width, height),
		NewFrameBuffer(width, height),
	}

	ovrStuff := new(OvrStuff).Init(width, height, framebuffers[0])

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
		Framebuffer:        framebuffers,
		ScreenQuad:         &ScreenQuad{},
		ScreenQuadRenderer: NewScreenQuadRenderer(),
		DebugRenderer:      NewLineRenderer(),
		OvrStuff:           ovrStuff,
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

	p0 := ww.Player.Camera.Pos4f()

	if this.riftRender {
		//w, h := this.Framebuffer[0].W, this.Framebuffer[0].H
		proj := this.Proj
		this.OvrStuff.Hmd.BeginFrame(0)
		for i := 0; i < 2; i++ {
			eye := this.OvrStuff.HmdDesc.EyeRenderOrder[i]
			pose := this.OvrStuff.Hmd.BeginEyeRender(eye)

			v0 := Vec3(this.OvrStuff.EyeRenderDesc[eye].ViewAdjust)
			v3 := v0.Add(Vec3(pose.Position)).Vec4(0)

			camera := ww.Player.Camera
			camera.MoveRelative(v3)
			this.Proj = this.OvrStuff.Proj[eye]
			this.View = (ww.PortalTransform(p0, camera.Pos4f()).Mul4(camera.Model())).Inv()
			viewport := this.OvrStuff.ViewportsFramebuffer[eye]
			viewport.Activate()
			fmt.Println(viewport)
			this.render(ww, options, viewport, 0, nil)
			this.OvrStuff.Hmd.EndEyeRender(eye, pose, this.OvrStuff.Textures[eye].Texture())
		}
		this.OvrStuff.Hmd.EndFrame()
		this.Proj = proj
	} else {
		w, h := window.GetSize()
		viewport := Viewport{0, 0, w, h}
		this.render(ww, options, viewport, 0, nil)
		//fmt.Println(viewport)
		viewport.Activate()
		gl.ActiveTexture(gl.TEXTURE0)
		this.Framebuffer[0].RenderTexture.Bind(target)
		this.ScreenQuadRenderer.Render(this.ScreenQuad, this.Proj, this.View, this.ClippingPlane_ws, mgl.Vec2{float32(w), float32(h)})
	}

	if this.screenShot {
		this.screenShot = false
		helpers.SaveTexture(gl.TEXTURE_RECTANGLE, 0, "screenshot.png")
	}

	this.FrameIndex++
}

func Vec3(v ovr.Vector3f) mgl.Vec3 {
	return mgl.Vec3{v.X, v.Y, v.Z}
}

func Vec4(v ovr.Vector3f, w float32) mgl.Vec4 {
	return mgl.Vec4{v.X, v.Y, v.Z, w}
}

type Viewport struct {
	X, Y, W, H int
}

func (this *Viewport) ToOvrRecti() (rect ovr.Recti) {
	rect.Pos.X = int32(this.X)
	rect.Pos.Y = int32(this.Y)
	rect.Size.W = int32(this.W)
	rect.Size.H = int32(this.H)
	return
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
