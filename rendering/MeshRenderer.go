package rendering

import (
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

func NewMeshRenderer() (this *Renderer) {
	return NewRenderer(helpers.MakeProgram("Mesh.vs", "Mesh.fs"), "mesh", nil, nil)
}
