package gamestate

import mgl "github.com/Jragonmiris/mathgl"

var NpcMesh = LoadMeshManaged("meshes/Torso.obj")

type Npc Entity

func (this *Npc) GetMesh() IMesh {
	return NpcMesh
}

func (this *Npc) GetModel() mgl.Mat4f {
	return (*Entity)(this).Model()
}

func (this *Npc) Entity() *Entity {
	return (*Entity)(this)
}
