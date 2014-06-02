package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

var NpcMesh = LoadMeshManaged("meshes/Torso.obj")

type Npc Entity

func (this *Npc) GetMesh() renderstuff.IMesh {
	return NpcMesh
}

func (this *Npc) GetModel() mgl.Mat4f {
	return (*Entity)(this).Model()
}

func (this *Npc) Entity() *Entity {
	return (*Entity)(this)
}
