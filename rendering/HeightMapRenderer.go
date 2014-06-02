package rendering

import (
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/constants"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

func NewHeightMapRenderer() *Renderer {
	return NewRenderer(
		helpers.MakeProgram("HeightMap.vs", "HeightMap.fs"),
		"height map",
		nil,
		HeightMapUpdate,
	)
}

func HeightMapUpdate(loc *RenderLocations, entity renderstuff.IRenderEntity, etc interface{}) {
	heightMap := entity.(*gamestate.HeightMap)

	if heightMap.HasChanges {
		min_h, max_h := heightMap.Bounds()
		loc.LowerBound.Uniform3f(0, 0, min_h)
		loc.UpperBound.Uniform3f(float32(heightMap.W), float32(heightMap.H), max_h)
		gl.ActiveTexture(gl.TEXTURE0 + constants.TextureHeightMap)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, heightMap.W, heightMap.H, gl.RED, gl.FLOAT, heightMap.TexturePixels())
		gl.ActiveTexture(gl.TEXTURE0)
		heightMap.HasChanges = false
	}
}
