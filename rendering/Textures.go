package rendering

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type Textures struct {
	Textures []gl.Texture
}

func NewTextures(heightMap *gamestate.HeightMap) *Textures {
	textures := make([]gl.Texture, 0, 7)

	// TEXTURE0 is only used for temporarly bound textures

	gl.ActiveTexture(gl.TEXTURE1)
	detailTexture, err := helpers.LoadTexture2D("textures/GravelCobbleS.jpg")
	if err != nil {
		fmt.Println(err)
	} else {
		textures = append(textures, detailTexture)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	}

	gl.ActiveTexture(gl.TEXTURE2)
	slopeTexture, err := helpers.LoadTexture2D("textures/Cliffs0149_18_S.png")
	if err != nil {
		fmt.Println("cant load GravelCobble0003_2_S.jpg")
		fmt.Println(err)
	} else {
		textures = append(textures, slopeTexture)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	}

	gl.ActiveTexture(gl.TEXTURE3)
	colorTexture, err := helpers.LoadTexture1D("textures/gradient.png")
	if err != nil {
		fmt.Println(err)
	} else {
		textures = append(textures, colorTexture)
		gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}

	gl.ActiveTexture(gl.TEXTURE4)
	heightMapTexture := gl.GenTexture()
	textures = append(textures, heightMapTexture)
	heightMapTexture.Bind(gl.TEXTURE_2D)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R16, heightMap.W, heightMap.H, 0, gl.RED, gl.FLOAT, heightMap.TexturePixels())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.ActiveTexture(gl.TEXTURE5)
	palmTexture, err := helpers.LoadTexture2D("textures/palme.png")
	if err != nil {
		fmt.Println("can't load palme.png")
		fmt.Println(err)
	} else {
		textures = append(textures, palmTexture)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	}

	gl.ActiveTexture(gl.TEXTURE6)
	firebullTexture, err := helpers.LoadTexture2D("textures/fireball.png")
	if err != nil {
		panic("fireball.png")
	} else {
		textures = append(textures, firebullTexture)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	}

	gl.ActiveTexture(gl.TEXTURE7)

	skybox, err := helpers.LoadTextureCube("textures/Above_The_Sea.jpg")
	if err != nil {
		panic("Above_The_Sea.jpg")
	} else {
		textures = append(textures, skybox)
		gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		//gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
	}

	gl.ActiveTexture(gl.TEXTURE0)

	return &Textures{textures}
}

func (this *Textures) Delete() {
	gl.DeleteTextures(this.Textures)
}
