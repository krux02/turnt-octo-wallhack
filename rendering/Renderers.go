package rendering

import (
	"github.com/go-gl/gl"
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/constants"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

func NewSkyboxRenderer() *renderstuff.Renderer {
	program := helpers.MakeProgram("Skybox.vs", "Skybox.fs")
	return renderstuff.NewRenderer(program, "Skybox", nil, nil)
}

func NewTreeRenderer() *renderstuff.Renderer {
	program := helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	return renderstuff.NewRenderer(program, "TreeSprite", nil, TreeUpdate)
}

func TreeUpdate(loc *renderstuff.RenderLocations, entiy renderstuff.IRenderEntity, additionalUniforms interface{}) {
	Rot2D := helpers.Mat4toMat3(additionalUniforms.(mgl.Mat4))
	loc.Rot2D.UniformMatrix3f(false, renderstuff.GlMat3(&Rot2D))
}

type WaterRenderUniforms struct {
	Time         float64
	CameraPos_ws mgl.Vec4
}

func WaterUpdate(loc *renderstuff.RenderLocations, entity renderstuff.IRenderEntity, etc interface{}) {
	water := entity.(*gamestate.Water)
	uniforms := etc.(WaterRenderUniforms)

	loc.Time.Uniform1f(float32(uniforms.Time))
	lb, ub := water.LowerBound, water.UpperBound
	loc.LowerBound.Uniform3f(lb[0], lb[1], lb[2])
	loc.UpperBound.Uniform3f(ub[0], ub[1], ub[2])
	v := uniforms.CameraPos_ws
	loc.CameraPos_ws.Uniform4f(v[0], v[1], v[2], v[3])
	loc.WaterHeight.Uniform1f(water.Height)
}

func NewSurfaceWaterRenderer() *renderstuff.Renderer {
	program := helpers.MakeProgram("Water.vs", "Water.fs")
	name := "Water"
	return renderstuff.NewRenderer(program, name, nil, WaterUpdate)
}

func NewDebugWaterRenderer() *renderstuff.Renderer {
	program := helpers.MakeProgram3("Water.vs", "Normal.gs", "Line.fs")
	name := "Water Normals"
	renderer := renderstuff.NewRenderer(program, name, nil, WaterUpdate)
	renderer.OverrideModeToPoints = true
	return renderer
}

type PortalRenderUniforms struct {
	Viewport Viewport
	Image    int
}

func NewPortalRenderer() *renderstuff.Renderer {
	program := helpers.MakeProgram("Portal.vs", "Portal.fs")
	return renderstuff.NewRenderer(program, "Portal", nil, PortalRenderUpdate)
}

func PortalRenderUpdate(loc *renderstuff.RenderLocations, entity renderstuff.IRenderEntity, etc interface{}) {
	uniforms := etc.(PortalRenderUniforms)
	loc.Viewport.Uniform4f(float32(uniforms.Viewport.X), float32(uniforms.Viewport.Y), float32(uniforms.Viewport.W), float32(uniforms.Viewport.H))
	loc.Image.Uniform1i(uniforms.Image)
}

func NewMeshRenderer() (this *renderstuff.Renderer) {
	return renderstuff.NewRenderer(helpers.MakeProgram("Mesh.vs", "Mesh.fs"), "mesh", nil, nil)
}

func NewScreenQuadRenderer() (this *renderstuff.Renderer) {
	program := helpers.MakeProgram("ScreenQuad.vs", "ScreenQuad.fs")
	return renderstuff.NewRenderer(program, "ScreenQuad", nil, ScreenQuadUpdate)
}

type ScreenQuadData struct {
	TextureSize  mgl.Vec2
	ViewportSize mgl.Vec2
}

func ScreenQuadUpdate(loc *renderstuff.RenderLocations, entity renderstuff.IRenderEntity, etc interface{}) {
	data := etc.(ScreenQuadData)
	loc.ViewPortSize.Uniform2f(data.ViewportSize.Elem())
	loc.TextureSize.Uniform2f(data.TextureSize.Elem())
}

func NewHeightMapRenderer() *renderstuff.Renderer {
	return renderstuff.NewRenderer(
		helpers.MakeProgram("HeightMap.vs", "HeightMap.fs"),
		"height map",
		nil,
		HeightMapUpdate,
	)
}

func HeightMapUpdate(loc *renderstuff.RenderLocations, entity renderstuff.IRenderEntity, etc interface{}) {
	heightMap := entity.(*gamestate.HeightMap)

	if heightMap.HasChanges {
		min_h, max_h := heightMap.Bounds()
		loc.LowerBound.Uniform3f(0, 0, min_h)
		loc.UpperBound.Uniform3f(float32(heightMap.W), float32(heightMap.H), max_h)
		gl.ActiveTexture(gl.TEXTURE0 + constants.TextureHeightMap)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, heightMap.W, heightMap.H, gl.RED, gl.FLOAT, heightMap.TexturePixels())
		gl.ActiveTexture(gl.TEXTURE0)
		heightMap.HasChanges = false
	}
}
