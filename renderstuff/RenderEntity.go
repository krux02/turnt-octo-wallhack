package renderstuff

import mgl "github.com/krux02/mathgl/mgl32"

type IRenderEntity interface {
	Mesh() *Mesh
	Model() mgl.Mat4
}
