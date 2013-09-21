package world

type World struct {
	HeightMap *HeightMap
	Portal    *Mesh
}

const W,H = 64,64

func NewWorld() (world *World) {
	heights := NewHeightMap(W, H)
	heights.DiamondSquare(W)

	PortalMesh := LoadMesh("meshes/Portal.blend")

	world = &World{heights, PortalMesh}

	return
}
