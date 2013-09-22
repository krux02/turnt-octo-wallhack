package world

import (
	mgl "github.com/krux02/mathgl"
)

type World struct {
	HeightMap *HeightMap
	Portals   []Portal
}

const W, H = 64, 64

func NewWorld() (world *World) {
	heights := NewHeightMap(W, H)
	heights.DiamondSquare(W)

	PortalPositions := []mgl.Vec3f{mgl.Vec3f{10, 10, 15}, mgl.Vec3f{30, 30, 10}, mgl.Vec3f{60, 60, 9}}

	PortalMesh := LoadMesh("meshes/Portal.blend")

	Portals := make([]Portal, len(PortalPositions))
	for i, pos := range PortalPositions {
		Portals[i] = Portal{
			Position:    pos,
			Orientation: mgl.QuatIdentf(),
			Mesh:        PortalMesh,
			Target:      nil,
		}
	}
	for i := range Portals {
		Portals[i].Target = &Portals[(i+1)%len(Portals)]
	}

	world = &World{heights, Portals}
	return
}
