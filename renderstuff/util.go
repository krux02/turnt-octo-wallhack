package renderstuff

import (
	mgl "github.com/krux02/mathgl/mgl32"
)

func GlMat4(mat *mgl.Mat4) *[16]float32 {
	return (*[16]float32)(mat)
}

func GlMat3(mat *mgl.Mat3) *[9]float32 {
	return (*[9]float32)(mat)
}

func GlMat2(mat *mgl.Mat2) *[4]float32 {
	return (*[4]float32)(mat)
}
