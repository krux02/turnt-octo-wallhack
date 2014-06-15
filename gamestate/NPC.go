package gamestate

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

var NpcMesh = (*renderstuff.Mesh)(LoadMeshManaged("meshes/Torso.obj"))

type Npc Entity

func (this *Npc) Mesh() *renderstuff.Mesh {
	return NpcMesh
}

func (this *Npc) Model() mgl.Mat4 {
	return (*Entity)(this).Model()
}

func (this *Npc) Entity() *Entity {
	return (*Entity)(this)
}
