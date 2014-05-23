package rendering

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

func (this *WorldRenderer) RenderEntity(View mgl.Mat4f, ClippingPlane_ws mgl.Vec4f, entity interface{}) {
	switch e := entity.(type) {
	case *gamestate.Npc:
		mesh := e.Mesh()
		Model := e.Entity().Model()
		this.MeshRenderer.Render(mesh, this.Proj, View, Model, ClippingPlane_ws)
	default:
		panic(fmt.Sprintf("unknown entity type %v", entity))
	}
}
