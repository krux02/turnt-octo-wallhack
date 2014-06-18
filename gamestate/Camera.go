package gamestate

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/math32"
)

type Camera struct {
	Entity
}

func NewCameraFromPos3f(p mgl.Vec3) *Camera {
	return NewCameraFromPos4f(mgl.Vec4{p[0], p[1], p[2], 1})
}

func NewCameraFromPos4f(position mgl.Vec4) *Camera {
	return &Camera{Entity{position, mgl.QuatIdent()}}
}

func NewCameraFromLookAt(eye mgl.Vec3, center mgl.Vec3, up mgl.Vec3) *Camera {
	return NewCameraFromMat4(mgl.LookAtV(eye, center, up))
}

func (camera *Camera) SetCameraFromLookAt(eye mgl.Vec3, center mgl.Vec3, up mgl.Vec3) {
	camera.SetView(mgl.LookAtV(eye, center, up))
}

func NewCameraFromMat4(view mgl.Mat4) (camera *Camera) {
	camera = new(Camera)
	camera.SetView(view)
	return camera
}

func (camera *Camera) MoveAbsolute(dist mgl.Vec4) {
	camera.Position[0] += dist[0]
	camera.Position[1] += dist[1]
	camera.Position[2] += dist[2]
}

func (camera *Camera) MoveRelative(dist mgl.Vec4) {
	dist_xyz := dist.Vec3()
	v := camera.Orientation.Rotate(dist_xyz)
	camera.MoveAbsolute(v.Vec4(0))
}

func (camera *Camera) Pos4f() mgl.Vec4 {
	p := camera.Position
	return mgl.Vec4{p[0], p[1], p[2], 1}
}

func (camera *Camera) Rotation2D() mgl.Mat4 {
	Orientation := camera.Orientation.Rotate(mgl.Vec3{0, 0, -1})
	angle := math32.Atan2(Orientation[1], Orientation[0])
	Rot2D := mgl.Rotate3DZ(angle)
	return helpers.Mat3toMat4(Rot2D)
}

func (camera *Camera) Rotate(angle float32, axis mgl.Vec3) {
	quat2 := mgl.QuatRotate(angle, axis).Normalize()
	camera.Orientation = camera.Orientation.Mul(quat2).Normalize()
}

func (camera *Camera) Yaw(yaw float32) {
	camera.Rotate(yaw, mgl.Vec3{1, 0, 0})
}

func (camera *Camera) Pitch(pitch float32) {
	camera.Rotate(pitch, mgl.Vec3{0, 1, 0})
}

func (camera *Camera) Roll(roll float32) {
	camera.Rotate(roll, mgl.Vec3{0, 0, 1})
}

func (camera *Camera) DirVec() (dir mgl.Vec3) {
	dir = camera.Orientation.Rotate(mgl.Vec3{0, 0, -1})
	return
}
