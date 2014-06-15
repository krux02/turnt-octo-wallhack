package gamestate

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type World struct {
	HeightMap      *HeightMap
	Water          *Water
	KdTree         KdTree
	Portals        []*Portal
	Trees          renderstuff.IRenderEntity
	ExampleObjects []renderstuff.IRenderEntity
	Player         *Player
}

type WrapVec4 mgl.Vec4

func (this WrapVec4) Dimension(dim int) float32 {
	return this[dim]
}

func portalFilter(kdElement KdElement) bool {
	switch kdElement.(type) {
	case *Portal:
		return true
	default:
		return false
	}
}

func (this *World) NearestPortal(pos mgl.Vec4) *Portal {
	return this.KdTree.NearestQuery(WrapVec4(pos), portalFilter).(*Portal)
}

func (this *World) NearestPortal2D(pos mgl.Vec2) *Portal {
	panic("not implemented")
}
