package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type WaterRenderer struct{ Renderer }

type DebugWaterRenderer struct{ WaterRenderer }

type SurfaceWaterRenderer struct{ WaterRenderer }

type WaterRenderUniforms struct {
	Time         float64
	CameraPos_ws mgl.Vec4f
}

func (this *WaterRenderer) Update(entity gamestate.IRenderEntity, etc interface{}) {
	water := entity.(*gamestate.Water)
	this.Program.Use()
	uniforms := etc.(WaterRenderUniforms)
	this.RenLoc.Time.Uniform1f(float32(uniforms.Time))
	lb, ub := water.LowerBound, water.UpperBound
	this.RenLoc.LowerBound.Uniform3f(lb[0], lb[1], lb[2])
	this.RenLoc.UpperBound.Uniform3f(ub[0], ub[1], ub[2])
	v := uniforms.CameraPos_ws
	this.RenLoc.CameraPos_ws.Uniform4f(v[0], v[1], v[2], v[3])
	this.RenLoc.WaterHeight.Uniform1f(water.Height)
}

func (this *WaterRenderer) Init() {
	this.Program.Use()
	helpers.BindLocations("water", this.Program, &this.RenLoc)
	this.RenLoc.HeightMap.Uniform1i(4)
	this.RenLoc.GroundTexture.Uniform1i(1)
	this.RenLoc.Skybox.Uniform1i(7)
}

func NewSurfaceWaterRenderer() (this *WaterRenderer) {
	this = new(WaterRenderer)
	this.Program = helpers.MakeProgram("Water.vs", "Water.fs")
	this.Init()
	return
}

func NewDebugWaterRenderer() (this *WaterRenderer) {
	this = new(WaterRenderer)
	this.Program = helpers.MakeProgram3("Water.vs", "Normal.gs", "Line.fs")
	this.Init()
	this.OverrideModeToPoints = true
	return
}
