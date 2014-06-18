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

func (this *World) PortalTransform(pos1, pos2 mgl.Vec4) mgl.Mat4 {
	portal := this.NearestPortal(pos1)
	m := portal.View()
	pos1 = m.Mul4x1(pos1)
	pos2 = m.Mul4x1(pos2)
	n := portal.Normal

	a := n.Dot(pos1)
	b := n.Dot(pos2)

	var portalPassed bool
	if a*b < 0 {
		c := a / (a - b)
		pos3 := pos1.Mul(c).Add(pos2.Mul(1 - c))
		portalPassed = -1 < pos3[0] && pos3[0] < 1 && -1 < pos3[1] && pos3[1] < 1
	} else {
		portalPassed = false
	}

	if portalPassed {
		return portal.Transform()
	} else {
		return mgl.Ident4()
	}
}
