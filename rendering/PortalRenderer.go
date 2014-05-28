package rendering

import (
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

func NewPortalRenderer() (mr *Renderer) {
	mr = new(Renderer)
	mr.Program = helpers.MakeProgram("Portal.vs", "Portal.fs")
	mr.Program.Use()
	helpers.BindLocations("portal", mr.Program, &mr.RenLoc)
	return
}
