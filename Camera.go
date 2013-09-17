package main

//#include <string.h>
import "C"

import (
	"fmt"
	"github.com/krux02/mathgl"
)

func stringMat4f(m mathgl.Mat4f) {
	s1 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[0], m[4], m[8], m[12])
	s2 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[1], m[5], m[9], m[13])
	s3 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[2], m[6], m[10], m[14])
	s4 := fmt.Sprintf("%2.2f\t%2.2f\t%2.2f\t%2.2f\n", m[3], m[7], m[11], m[15])
	fmt.Sprintf("%s%s%s%s", s1, s2, s3, s4)
}

type Camera struct {
	position  mathgl.Vec3f
	direction mathgl.Quatf
}

func NewCamera(eye mathgl.Vec3f, center mathgl.Vec3f, up mathgl.Vec3f) *Camera {
	return nil
}

func (camera *Camera) MoveAbsolute(dist mathgl.Vec3f) {
	camera.position = camera.position.Add(dist)
}

func (camera *Camera) MoveRelative(dist mathgl.Vec3f) {
	camera.MoveAbsolute(camera.direction.Rotate(dist))
}

func (camera *Camera) View() mathgl.Mat4f {
	direction := camera.direction.Rotate(mathgl.Vec3f{0, 0, -1})
	center := camera.position.Add(direction)
	up := camera.direction.Rotate(mathgl.Vec3f{0, 1, 0})
	return mathgl.LookAtV(camera.position, center, up)
}

func (camera *Camera) Rotate(angle float32, axis mathgl.Vec3f) {
	quat2 := mathgl.QuatRotatef(angle, axis)
	camera.direction = camera.direction.Mul(quat2).Normalize()
	//camera.direction = MyMult(camera.direction, quat2).Normalize()
}

func (camera *Camera) Yaw(yaw float32) {
	camera.Rotate(yaw, mathgl.Vec3f{1, 0, 0})
}

func (camera *Camera) Pitch(pitch float32) {
	camera.Rotate(pitch, mathgl.Vec3f{0, 1, 0})
}

func (camera *Camera) Roll(roll float32) {
	camera.Rotate(roll, mathgl.Vec3f{0, 0, 1})
}
