package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	//"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type TreeRenderer struct{ Renderer }

func NewTreeRenderer() *TreeRenderer {
	renderer := new(TreeRenderer)
	renderer.Program = helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	renderer.Program.Use()
	helpers.BindLocations("palm sprite", renderer.Program, &renderer.RenLoc)
	renderer.RenLoc.PalmTree.Uniform1i(5)
	return renderer
}

func (this *TreeRenderer) Update(entiy gamestate.IRenderEntity, additionalUniforms interface{}) {
	Rot2D := helpers.Mat4toMat3(additionalUniforms.(mgl.Mat4f))
	this.RenLoc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
}
