package world

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type Portal struct {
	Position    mgl.Vec3f
	Orientation mgl.Quatf
	Mesh 		*Mesh
	Target		*Portal
}

