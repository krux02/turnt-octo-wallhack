package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

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

func (this *WorldRenderer) ClearRenderData() {
	for _, value := range this.RenData {
		value.Delete()
	}
}
