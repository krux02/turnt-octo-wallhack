package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

func NewSkyboxRenderer() *Renderer {
	program := helpers.MakeProgram("Skybox.vs", "Skybox.fs")
	return NewRenderer(program, "Skybox", nil, nil)
}

func NewTreeRenderer() *Renderer {
	program := helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	return NewRenderer(program, "TreeSprite", nil, TreeUpdate)
}

func TreeUpdate(loc *RenderLocations, entiy renderstuff.IRenderEntity, additionalUniforms interface{}) {
	Rot2D := helpers.Mat4toMat3(additionalUniforms.(mgl.Mat4f))
	loc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
}

type WaterRenderUniforms struct {
	Time         float64
	CameraPos_ws mgl.Vec4f
}

func WaterUpdate(loc *RenderLocations, entity renderstuff.IRenderEntity, etc interface{}) {
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

func NewSurfaceWaterRenderer() *Renderer {
	program := helpers.MakeProgram("Water.vs", "Water.fs")
	name := "Water"
	return NewRenderer(program, name, nil, WaterUpdate)
}

func NewDebugWaterRenderer() *Renderer {
	program := helpers.MakeProgram3("Water.vs", "Normal.gs", "Line.fs")
	name := "Water Normals"
	renderer := NewRenderer(program, name, nil, WaterUpdate)
	renderer.OverrideModeToPoints = true
	return renderer
}

func NewPortalRenderer() *Renderer {
	program := helpers.MakeProgram("Portal.vs", "Portal.fs")
	return NewRenderer(program, "Portal", nil, nil)
}

func NewMeshRenderer() (this *Renderer) {
	return NewRenderer(helpers.MakeProgram("Mesh.vs", "Mesh.fs"), "mesh", nil, nil)
}

func NewScreenQuadRenderer() (this *Renderer) {
	program := helpers.MakeProgram("ScreenQuad.vs", "ScreenQuad.fs")
	return NewRenderer(program, "ScreenQuad", nil, nil)
}
