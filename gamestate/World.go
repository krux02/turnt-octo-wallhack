package gamestate

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

	PortalPositions := []mgl.Vec4f{mgl.Vec4f{10, 10, 15, 1}, mgl.Vec4f{30, 30, 10, 1}} // , mgl.Vec3f{60, 60, 9}

	//PortalMesh := LoadMesh("meshes/Portal.blend")
	//PortalMesh := PortalRect()
	PortalMesh := QuadMesh()
	normal := mgl.Vec4f{0, 0, 1, 0}

	Portals := make([]*Portal, len(PortalPositions))
	for i, pos := range PortalPositions {
		Portals[i] = &Portal{
			Entity{Position: pos, Orientation: mgl.QuatIdentf()},
			PortalMesh,
			normal,
			nil,
		}
	}
	for i := range Portals {
		j := i ^ 1
		Portals[i].Target = Portals[j]
	}

	world = &World{heights, Portals}
	return
}

func (this *World) NearestPortal(pos mgl.Vec4f) *Portal {
	pos = pos.Mul(1 / pos[3])
	var dist float32 = math.MaxFloat32
	var nearestPortal *Portal
	for _, portal := range this.Portals {
		p := portal.Position
		p = p.Mul(1 / p[3])
		newDist := pos.Sub(p).Len()
		if newDist < dist {
			dist = newDist
			nearestPortal = portal
		}
	}
	return nearestPortal
}
