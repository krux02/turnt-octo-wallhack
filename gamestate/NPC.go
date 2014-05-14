package gamestate

var NpcMesh = LoadMeshManaged("meshes/Torso.obj")

type MeshEntity interface {
	Mesh() *Mesh
	Entity() *Entity
}

type Npc Entity

func (this *Npc) Mesh() *Mesh {
	return NpcMesh
}

func (this *Npc) Entity() *Entity {
	return (*Entity)(this)
}
