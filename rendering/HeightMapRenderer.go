package rendering

import (
	//"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type HeightMapRenderer struct {
	Program gl.Program
	RenLoc  RenderLocations
}

func NewHeightMapRenderer() (this *HeightMapRenderer) {

	this = new(HeightMapRenderer)

	this.Program = helpers.MakeProgram("HeightMap.vs", "HeightMap.fs")
	this.Program.Use()

	helpers.BindLocations("height map", this.Program, &this.RenLoc)

	this.RenLoc.HeightMap.Uniform1i(4)
	this.RenLoc.ColorBand.Uniform1i(3)
	this.RenLoc.Slope.Uniform1i(2)
	this.RenLoc.Texture.Uniform1i(1)

	return
}

func (wr *HeightMapRenderer) Delete() {
	wr.Program.Delete()
}

func (this *HeightMapRenderer) Render(renderData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, clippingPlane mgl.Vec4f, _ map[string]int) {
	this.Program.Use()
	renderData.VAO.Bind()

	this.RenLoc.ClippingPlane_ws.Uniform4f(clippingPlane[0], clippingPlane[1], clippingPlane[2], clippingPlane[3])

	numverts := renderData.Numverts

	this.RenLoc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	this.RenLoc.View.UniformMatrix4f(false, glMat4(&View))
	this.RenLoc.Model.UniformMatrix4f(false, glMat4(&Model))

	gl.DrawElements(gl.TRIANGLES, numverts, gl.UNSIGNED_INT, uintptr(0))
}

func (this *HeightMapRenderer) RenderLocations() *RenderLocations {
	return &this.RenLoc
}

func (this *HeightMapRenderer) UseProgram() {
	this.Program.Use()
}

func (this *HeightMapRenderer) Update(entity gamestate.IRenderEntity, etc interface{}) {
	heightMap := entity.(*gamestate.HeightMap)

	if heightMap.HasChanges {

		min_h, max_h := heightMap.Bounds()
		this.RenLoc.LowerBound.Uniform3f(0, 0, min_h)
		this.RenLoc.UpperBound.Uniform3f(float32(heightMap.W), float32(heightMap.H), max_h)

		gl.ActiveTexture(gl.TEXTURE4)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R16, heightMap.W, heightMap.H, 0, gl.RED, gl.FLOAT, heightMap.TexturePixels())
		gl.ActiveTexture(gl.TEXTURE0)

		heightMap.HasChanges = false
	}
}
