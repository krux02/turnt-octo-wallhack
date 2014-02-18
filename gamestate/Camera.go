package gamestate

//#include <string.h>
import "C"

import (
	mgl "github.com/Jragonmiris/mathgl"
	"math"
)

type Camera struct {
	Entity
}

func NewCameraFromPos3f(p mgl.Vec3f) *Camera {
	return NewCameraFromPos4f(mgl.Vec4f{p[0], p[1], p[2], 1})
}

func NewCameraFromPos4f(position mgl.Vec4f) *Camera {
	return &Camera{Entity{position, mgl.QuatIdentf()}}
}

func NewCameraFromLookAt(eye mgl.Vec3f, center mgl.Vec3f, up mgl.Vec3f) *Camera {
	return NewCameraFromMat4(mgl.LookAtV(eye, center, up))
}

func (camera *Camera) SetCameraFromLookAt(eye mgl.Vec3f, center mgl.Vec3f, up mgl.Vec3f) {
	camera.SetView(mgl.LookAtV(eye, center, up))
}

func NewCameraFromMat4(view mgl.Mat4f) (camera *Camera) {
	camera = new(Camera)
	camera.SetView(view)
	return camera
}

func (camera *Camera) MoveAbsolute(dist mgl.Vec3f) {
	camera.Position[0] += dist[0]
	camera.Position[1] += dist[1]
	camera.Position[2] += dist[2]
}

func (camera *Camera) MoveRelative(dist mgl.Vec3f) {
	camera.MoveAbsolute(camera.Orientation.Rotate(dist))
}

func (camera *Camera) Pos4f() mgl.Vec4f {
	p := camera.Position
	return mgl.Vec4f{p[0], p[1], p[2], 1}
}

func (camera *Camera) Rotation2D() (Rot2D mgl.Mat3f) {
	Orientation := camera.Orientation.Rotate(mgl.Vec3f{0, 0, -1})
	angle := math.Atan2(float64(Orientation[1]), float64(Orientation[0]))
	Rot2D = mgl.Rotate3DZ(float32(angle / math.Pi * 180))
	return
}

func (camera *Camera) Rotate(angle float32, axis mgl.Vec3f) {
	quat2 := mgl.QuatRotatef(angle, axis).Normalize()
	camera.Orientation = camera.Orientation.Mul(quat2).Normalize()
}

func (camera *Camera) Yaw(yaw float32) {
	camera.Rotate(yaw, mgl.Vec3f{1, 0, 0})
}

func (camera *Camera) Pitch(pitch float32) {
	camera.Rotate(pitch, mgl.Vec3f{0, 1, 0})
}

func (camera *Camera) Roll(roll float32) {
	camera.Rotate(roll, mgl.Vec3f{0, 0, 1})
}

func (camera *Camera) DirVec() (dir mgl.Vec3f) {
	dir = camera.Orientation.Rotate(mgl.Vec3f{0, 0, -1})
	return
}
