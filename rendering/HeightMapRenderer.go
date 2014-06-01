package rendering

import (
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type HeightMapRenderer struct{ Renderer }

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
