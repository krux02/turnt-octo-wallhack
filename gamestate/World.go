package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"math"
)

type World struct {
	HeightMap *HeightMap
	Portals   []*Portal
}

func (this *World) NearestPortal(pos mgl.Vec4f) *Portal {
	dist := float32(math.MaxFloat32)
	var nearestPortal *Portal
	for _, portal := range this.Portals {
		p := portal.Position
		p = p.Mul(1 / p[3])
		newDist := helpers.HomogenDist(portal.Position, pos)
		if newDist < dist {
			dist = newDist
			nearestPortal = portal
		}
	}
	return nearestPortal
}

func (this *World) NearestPortal2D(pos mgl.Vec2f) *Portal {
	dist := float32(math.MaxFloat32)
	var nearestPortal *Portal
	for _, portal := range this.Portals {
		p := mgl.Vec2f{portal.Position[0], portal.Position[1]}
		newDist := pos.Sub(p).Len()
		if newDist < dist {
			dist = newDist
			nearestPortal = portal
		}
	}
	return nearestPortal
}
