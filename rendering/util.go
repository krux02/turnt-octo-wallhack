package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
)

func glMat4(mat *mgl.Mat4f) *[16]float32 {
	return (*[16]float32)(mat)
}

func glMat3(mat *mgl.Mat3f) *[9]float32 {
	return (*[9]float32)(mat)
}

func glMat2(mat *mgl.Mat2f) *[4]float32 {
	return (*[4]float32)(mat)
}
