package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type World struct {
	HeightMap *HeightMap
	KdTree    KdTree
	Portals   []*Portal
	Palms     PalmTreesInstanceData
}

type WrapVec4f mgl.Vec4f

func (this WrapVec4f) Dimension(dim int) float32 {
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

func (this *World) NearestPortal(pos mgl.Vec4f) *Portal {
	return this.KdTree.NearestQuery(WrapVec4f(pos), portalFilter).(*Portal)
}

func (this *World) NearestPortal2D(pos mgl.Vec2f) *Portal {
	panic("not implemented")
}

func (this *World) Save(filename string) {

}
