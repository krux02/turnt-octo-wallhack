package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

type IRenderer interface {
	Render(renderData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f)
	UseProgram()
	RenderLocations() *RenderLocations
	Update(entiy gamestate.IRenderEntity, etc interface{})
}

type Renderer struct {
	Program gl.Program
	RenLoc  RenderLocations
}

func (this *Renderer) Delete() {
	this.Program.Delete()
	*this = Renderer{}
}

func (this *Renderer) UseProgram() {
	this.Program.Use()
}

func (this *Renderer) RenderLocations() *RenderLocations {
	return &this.RenLoc
}

func (this *Renderer) Update(entiy gamestate.IRenderEntity, etc interface{}) {}

func (this *Renderer) Render(renData *RenderData, Proj, View, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f) {
	this.Program.Use()
	renData.VAO.Bind()

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Model.UniformMatrix4f(false, glMat4(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))

	//Loc.Image.Uniform1i(additionalUniforms["Image"])
	Loc.ClippingPlane_ws.Uniform4f(ClippingPlane_ws[0], ClippingPlane_ws[1], ClippingPlane_ws[2], ClippingPlane_ws[3])
	numverts := renData.Numverts

	if renData.InstanceDataBuffer == 0 {
		if renData.Indices == 0 {
			gl.DrawArrays(renData.Mode, 0, renData.Numverts)
		} else {
			gl.DrawElements(renData.Mode, numverts, renData.IndexType, uintptr(0))
		}
	} else {
		if renData.Indices == 0 {
			gl.DrawArraysInstanced(renData.Mode, 0, renData.Numverts, renData.NumInstances)
		} else {
			gl.DrawElementsInstanced(renData.Mode, renData.Numverts, renData.IndexType, 0, renData.NumInstances)
		}
	}
}

func (this *WorldRenderer) RenderEntity(entity gamestate.IRenderEntity, View mgl.Mat4f, ClippingPlane_ws mgl.Vec4f, additionalUniforms interface{}) {
	var renderer IRenderer
	switch entity.(type) {
	case *gamestate.Npc:
		renderer = this.MeshRenderer
	case *gamestate.HeightMap:
		renderer = this.HeightMapRenderer
	case *gamestate.Forest:
		renderer = this.TreeRenderer
	case *gamestate.Portal:
		renderer = this.PortalRenderer
	case *gamestate.Water:
		renderer = this.WaterRenderer
		//renderer = this.DebugWaterRenderer
	default:
		panic(fmt.Sprintf("unknown entity type %v", entity))
	}

	renderer.UseProgram()
	renderer.Update(entity, additionalUniforms)

	mesh := entity.GetMesh()
	meshData := this.RenData[mesh]
	if meshData == nil {
		md := CreateMeshRenderData(mesh, renderer.RenderLocations())
		meshData = &md
		this.RenData[mesh] = &md
	}
	Model := entity.GetModel()
	renderer.Render(meshData, this.Proj, View, Model, ClippingPlane_ws)
}
