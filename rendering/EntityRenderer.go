package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

type Renderer interface {
	Render(renderData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f)
	UseProgram()
	RenderLocations() *RenderLocations
}

func (this *WorldRenderer) RenderEntity(entity gamestate.IRenderEntity, View mgl.Mat4f, ClippingPlane_ws mgl.Vec4f) {
	var renderer Renderer
	switch entity.(type) {
	case *gamestate.Npc:
		renderer = this.MeshRenderer
	default:
		panic(fmt.Sprintf("unknown entity type %v", entity))
	}
	renderer.UseProgram()
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
