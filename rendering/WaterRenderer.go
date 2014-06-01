package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/constants"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type WaterRenderer struct{ Renderer }

type WaterRenderUniforms struct {
	Time         float64
	CameraPos_ws mgl.Vec4f
}

func WaterUpdate(loc *RenderLocations, entity gamestate.IRenderEntity, etc interface{}) {
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

func WaterInit(loc *RenderLocations) {
	loc.HeightMap.Uniform1i(constants.TextureHeightMap)
	loc.GroundTexture.Uniform1i(constants.TextureGround)
	loc.Skybox.Uniform1i(constants.TextureSkybox)
}

func NewSurfaceWaterRenderer() *Renderer {
	program := helpers.MakeProgram("Water.vs", "Water.fs")
	name := "Water"
	return NewRenderer(program, name, WaterInit, WaterUpdate)
}

func NewDebugWaterRenderer() *Renderer {
	program := helpers.MakeProgram3("Water.vs", "Normal.gs", "Line.fs")
	name := "Water Normals"
	renderer := NewRenderer(program, name, WaterInit, WaterUpdate)
	renderer.OverrideModeToPoints = true
	return renderer
}
