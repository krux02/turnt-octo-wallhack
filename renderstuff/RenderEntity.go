package renderstuff

import mgl "github.com/Jragonmiris/mathgl"

type IRenderEntity interface {
	Mesh() *Mesh
	Model() mgl.Mat4f
}
