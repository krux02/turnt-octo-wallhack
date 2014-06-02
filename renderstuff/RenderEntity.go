package renderstuff

import mgl "github.com/Jragonmiris/mathgl"

type IRenderEntity interface {
	GetMesh() IMesh
	GetModel() mgl.Mat4f
}
