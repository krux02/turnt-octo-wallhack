package simulation

import (
	"github.com/krux02/turnt-octo-wallhack/rendering"
	"github.com/krux02/turnt-octo-wallhack/world"
)

func EnterPortal(portal *world.Portal, camera *rendering.Camera) {
	portal.SetFromMat4(portal.Target.ModelMat4().Mul(portal.ModelMat4().Inv()).Mul(camera.View()))
}

func PortalPassed(portal *Portal, pos1, pos2 mgl.Vec3f) bool {
	pos1w := mgl.Vec4f{pos1[0], pos1[1], pos1[2], 1}
	pos2w := mgl.Vec4f{pos2[0], pos2[1], pos2[2], 1}
	plane := portal.ClippingPlane(true)
	return pos1w.Dot(plane)*pos2w.Dot(plane) < 0
}
