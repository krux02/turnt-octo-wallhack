package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
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

func (this *TreeRenderer) Render(meshData *RenderData, Proj, View, Model mgl.Mat4f, clippingPlane mgl.Vec4f, additionalUniforms map[string]int) {
	Rot2D := helpers.Mat4toMat3(Model)
	this.Program.Use()
	meshData.VAO.Bind()

	this.RenLoc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	this.RenLoc.View.UniformMatrix4f(false, glMat4(&View))
	this.RenLoc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
	this.RenLoc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	for key, value := range additionalUniforms {
		loc := this.Program.GetUniformLocation(key)
		if loc != -1 {
			loc.Uniform1i(value)
		}
	}

	gl.DrawArraysInstanced(meshData.Mode, 0, meshData.Numverts, meshData.NumInstances)
}
