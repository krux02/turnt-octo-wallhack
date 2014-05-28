package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type TreeRenderer struct {
	Prog gl.Program
	Loc  RenderLocations
}

func (this *TreeRenderer) RenderLocations() *RenderLocations {
	return &this.Loc
}

func (this *TreeRenderer) Update(entiy gamestate.IRenderEntity) {}

func (this *TreeRenderer) UseProgram() {
	this.Prog.Use()
}

func (this *TreeRenderer) Delete() {
	this.Prog.Delete()
	*this = TreeRenderer{}
}

func NewTreeRenderer() *TreeRenderer {
	renderer := new(TreeRenderer)
	renderer.Prog = helpers.MakeProgram("Sprite.vs", "Sprite.fs")
	renderer.Prog.Use()
	helpers.BindLocations("palm sprite", renderer.Prog, &renderer.Loc)
	renderer.Loc.PalmTree.Uniform1i(5)
	return renderer
}

func (this *TreeRenderer) Render(meshData *RenderData, Proj, View, Model mgl.Mat4f, clippingPlane mgl.Vec4f, additionalUniforms map[string]int) {
	Rot2D := helpers.Mat4toMat3(Model)
	this.Prog.Use()
	meshData.VAO.Bind()

	this.Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	this.Loc.View.UniformMatrix4f(false, glMat4(&View))
	this.Loc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
	this.Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	for key, value := range additionalUniforms {
		loc := this.Prog.GetUniformLocation(key)
		if loc != -1 {
			loc.Uniform1i(value)
		}
	}

	gl.DrawArraysInstanced(gl.TRIANGLE_FAN, 0, meshData.Numverts, meshData.NumInstances)
}
