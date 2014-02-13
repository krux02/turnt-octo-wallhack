package gamestate

//#include <string.h>
import "C"

import (
	mgl "github.com/Jragonmiris/mathgl"
	"math"
)

type Camera struct {
	Position    mgl.Vec3f
	Orientation mgl.Quatf
}

func NewCameraFromLookAt(eye mgl.Vec3f, center mgl.Vec3f, up mgl.Vec3f) *Camera {
	return NewCameraFromMat4(mgl.LookAtV(eye, center, up))
}

func (camera *Camera) SetCameraFromLookAt(eye mgl.Vec3f, center mgl.Vec3f, up mgl.Vec3f) {
	camera.SetFromMat4(mgl.LookAtV(eye, center, up))
}

func NewCameraFromMat4(view mgl.Mat4f) (camera *Camera) {
	camera = new(Camera)
	camera.SetFromMat4(view)
	return camera
}

func (camera *Camera) SetFromMat4(view mgl.Mat4f) {
	m00 := view[0]
	m10 := view[1]
	m20 := view[2]

	m01 := view[4]
	m11 := view[5]
	m21 := view[6]

	m02 := view[8]
	m12 := view[9]
	m22 := view[10]

	m03 := view[12]
	m13 := view[13]
	m23 := view[14]

	qw := float32(math.Sqrt(float64(1+m00+m11+m22))) / 2
	qx := (m21 - m12) / (4 * qw)
	qy := (m02 - m20) / (4 * qw)
	qz := (m10 - m01) / (4 * qw)

	dir := mgl.Quatf{qw, mgl.Vec3f{qx, qy, qz}}.Inverse()
	pos := dir.Rotate(mgl.Vec3f{-m03, -m13, -m23})

	camera.Position = pos
	camera.Orientation = dir
}

func (camera *Camera) MoveAbsolute(dist mgl.Vec3f) {
	camera.Position = camera.Position.Add(dist)
}

func (camera *Camera) MoveRelative(dist mgl.Vec3f) {
	camera.MoveAbsolute(camera.Orientation.Rotate(dist))
}

func (camera *Camera) View() mgl.Mat4f {
	Tx := camera.Position[0]
	Ty := camera.Position[1]
	Tz := camera.Position[2]
	return camera.Orientation.Inverse().Mat4().Mul4(mgl.Translate3D(-Tx, -Ty, -Tz))
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
