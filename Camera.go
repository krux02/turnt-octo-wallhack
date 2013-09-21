package main

//#include <string.h>
import "C"

import (
	"fmt"
	mgl "github.com/krux02/mathgl"
	"math"
)

func stringMat4f(m mgl.Mat4f) {
	s1 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[0], m[4], m[8], m[12])
	s2 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[1], m[5], m[9], m[13])
	s3 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[2], m[6], m[10], m[14])
	s4 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[3], m[7], m[11], m[15])
	fmt.Sprintf("%s%s%s%s", s1, s2, s3, s4)
}

type Camera struct {
	position  mgl.Vec3f
	direction mgl.Quatf
}

func NewCamera(eye mgl.Vec3f, center mgl.Vec3f, up mgl.Vec3f) *Camera {
	return nil
}

func (camera *Camera) MoveAbsolute(dist mgl.Vec3f) {
	camera.position = camera.position.Add(dist)
}

func (camera *Camera) MoveRelative(dist mgl.Vec3f) {
	camera.MoveAbsolute(camera.direction.Rotate(dist))
}

func (camera *Camera) View() mgl.Mat4f {
	direction := camera.direction.Rotate(mgl.Vec3f{0, 0, -1})
	center := camera.position.Add(direction)
	up := camera.direction.Rotate(mgl.Vec3f{0, 1, 0})
	return mgl.LookAtV(camera.position, center, up)
}

func (camera *Camera) Rotation2D() (Rot2D mgl.Mat3f) {
	direction := camera.direction.Rotate(mgl.Vec3f{0, 0, -1})
	angle := math.Atan2(float64(direction[1]), float64(direction[0]))
	Rot2D = mgl.Rotate3DZ(angle/math.Pi*180)
	return
}

func (camera *Camera) Rotate(angle float32, axis mgl.Vec3f) {
	quat2 := mgl.QuatRotatef(angle, axis)
	camera.direction = camera.direction.Mul(quat2).Normalize()
	//camera.direction = MyMult(camera.direction, quat2).Normalize()
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
