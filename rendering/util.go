package rendering

import mgl "github.com/Jragonmiris/mathgl"

func glMat(mat *mgl.Mat4f) *[16]float32 {
	return (*[16]float32)(mat)
}
