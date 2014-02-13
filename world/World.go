package world

import (
	mgl "github.com/Jragonmiris/mathgl"
	"math"
)

type World struct {
	HeightMap *HeightMap
	Portals   []*Portal
}

const W, H = 64, 64

func NewWorld() (world *World) {
	heights := NewHeightMap(W, H)
	heights.DiamondSquare(W)

	PortalPositions := []mgl.Vec3f{mgl.Vec3f{10, 10, 15}, mgl.Vec3f{30, 30, 10}, mgl.Vec3f{60, 60, 9}}

	//PortalMesh := LoadMesh("meshes/Portal.blend")
	PortalMesh := PortalRect()

	Portals := make([]*Portal, len(PortalPositions))
	for i, pos := range PortalPositions {
		Portals[i] = &Portal{
			Position:    pos,
			Orientation: mgl.QuatIdentf(),
			Mesh:        PortalMesh,
			Target:      nil,
		}
	}
	for i := range Portals {
		Portals[i].Target = Portals[(i+1)%len(Portals)]
	}

	world = &World{heights, Portals}
	return
}

func (this *World) NearestPortal(pos mgl.Vec3f) *Portal {
	var dist float32 = math.MaxFloat32
	var nearestPortal *Portal
	for _, portal := range this.Portals {
		newDist := pos.Sub(portal.Position).Len()
		if newDist < dist {
			dist = newDist
			nearestPortal = portal
		}
	}
	return nearestPortal
}
