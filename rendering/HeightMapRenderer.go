package rendering

import (
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

func NewHeightMapRenderer() *Renderer {
	return NewRenderer(
		helpers.MakeProgram("HeightMap.vs", "HeightMap.fs"),
		"height map",
		HeightMapInit,
		HeightMapUpdate,
	)
}

func HeightMapInit(loc *RenderLocations) {
	loc.HeightMap.Uniform1i(4)
	loc.ColorBand.Uniform1i(3)
	loc.Slope.Uniform1i(2)
	loc.Texture.Uniform1i(1)
}

func HeightMapUpdate(loc *RenderLocations, entity gamestate.IRenderEntity, etc interface{}) {
	heightMap := entity.(*gamestate.HeightMap)

	if heightMap.HasChanges {
		min_h, max_h := heightMap.Bounds()
		loc.LowerBound.Uniform3f(0, 0, min_h)
		loc.UpperBound.Uniform3f(float32(heightMap.W), float32(heightMap.H), max_h)
		gl.ActiveTexture(gl.TEXTURE4)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, heightMap.W, heightMap.H, gl.RED, gl.FLOAT, heightMap.TexturePixels())
		gl.ActiveTexture(gl.TEXTURE0)
		heightMap.HasChanges = false
	}
}
