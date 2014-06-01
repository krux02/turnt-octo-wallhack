package rendering

import (
	//"fmt"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

func (this *WorldRenderer) RenderEntity(renderer *Renderer, entity gamestate.IRenderEntity, additionalUniforms interface{}) {
	renderer.Program.Use()

	mesh := entity.GetMesh()
	meshData := this.RenData[mesh]
	if meshData == nil {
		md := LoadMeshToGpu(mesh, &renderer.RenLoc)
		meshData = &md
		this.RenData[mesh] = &md
	}
	Model := entity.GetModel()

	renderer.Update(entity, additionalUniforms)
	renderer.Render(meshData, this.Proj, this.View, Model, this.ClippingPlane_ws)
}

func (this *WorldRenderer) ClearRenderData() {
	for _, value := range this.RenData {
		value.Delete()
	}
}
