package settings

import mgl "github.com/Jragonmiris/mathgl"

type BoolOptions struct {
	NoParticleRender,
	NoParticlePhysics,
	NoWorldRender,
	NoTreeRender,
	NoPlayerPhysics,
	DepthClamp,
	Wireframe bool
	StartPosition mgl.Vec4f
}
