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

func (pt *TreeRenderer) Render(meshData *RenderData, Proj, View, Model mgl.Mat4f, clippingPlane mgl.Vec4f) {
	Rot2D := helpers.Mat4toMat3(Model)
	pt.Prog.Use()
	meshData.VAO.Bind()
	pt.Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	pt.Loc.View.UniformMatrix4f(false, glMat4(&View))
	pt.Loc.Rot2D.UniformMatrix3f(false, glMat3(&Rot2D))
	pt.Loc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])
	gl.DrawArraysInstanced(gl.TRIANGLE_FAN, 0, meshData.Numverts, meshData.NumInstances)
}
