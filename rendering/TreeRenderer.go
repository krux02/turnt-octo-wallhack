package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

func NewTreeRenderer() *Renderer {
	program := helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	return NewRenderer(program, "TreeSprite", TreeInit, TreeUpdate)
}

func TreeInit(loc *RenderLocations) {
	loc.PalmTree.Uniform1i(5)
}

func TreeUpdate(loc *RenderLocations, entiy gamestate.IRenderEntity, additionalUniforms interface{}) {
	Rot2D := helpers.Mat4toMat3(additionalUniforms.(mgl.Mat4f))
	loc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
}
