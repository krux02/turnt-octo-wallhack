package rendering

import (
	"fmt"
	"github.com/go-gl/gl"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/settings"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

func (this *WorldRenderer) render(ww *gamestate.World, options *settings.BoolOptions, viewport Viewport, recursion int, srcPortal *gamestate.Portal) {

	fmt.Println(recursion)
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
			this.PortalRenderer.Render(nearestPortal, this.Proj, this.View, this.ClippingPlane_ws, PortalRenderUniforms{viewport, 0})
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
	if recursion < MaxRecursion {
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

			if scissor {
				//gl.Scissor(0, 0, w, h)
				gl.Disable(gl.SCISSOR_TEST)
			}
			this.Framebuffer[recursion].Bind()
			gl.Enable(gl.DEPTH_CLAMP)

			gl.ActiveTexture(gl.TEXTURE0)
			this.Framebuffer[recursion+1].RenderTexture.Bind(target)

			this.PortalRenderer.Render(nearestPortal, this.Proj, this.View, this.ClippingPlane_ws, PortalRenderUniforms{viewport, 0})
		}
	}
}
