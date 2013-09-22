package world

import (
	mgl "github.com/krux02/mathgl"
)

type Portal struct {
	Position    mgl.Vec3f
	Orientation mgl.Quatf
	Mesh 		*Mesh
	Target		*Portal
}

