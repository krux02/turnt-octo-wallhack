package rendering

import (
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

func NewPortalRenderer() *Renderer {
	program := helpers.MakeProgram("Portal.vs", "Portal.fs")
	return NewRenderer(program, "Portal", nil, nil)
}
