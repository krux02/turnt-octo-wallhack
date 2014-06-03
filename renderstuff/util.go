package renderstuff

import (
	mgl "github.com/Jragonmiris/mathgl"
)

func GlMat4(mat *mgl.Mat4f) *[16]float32 {
	return (*[16]float32)(mat)
}

func GlMat3(mat *mgl.Mat3f) *[9]float32 {
	return (*[9]float32)(mat)
}

func GlMat2(mat *mgl.Mat2f) *[4]float32 {
	return (*[4]float32)(mat)
}
