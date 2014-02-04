package settings

import mgl "github.com/Jragonmiris/mathgl"

type BoolOptions struct {
	NoParticleRender,
	NoParticlePhysics,
	NoWorldRender,
	NoTreeRender,
	NoPlayerPhysics,
	Wireframe bool
	portalOrientations []*mgl.Quatf
}
