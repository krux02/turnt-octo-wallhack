package rendering

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_ttf"
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
	detailTexture := gl.GenTexture()
	detailTexture.Bind(gl.TEXTURE_2D)
	err := helpers.LoadTexture2D("textures/GravelCobbleS.jpg")
	if err != nil {
		detailTexture.Delete()
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
	slopeTexture := gl.GenTexture()
	slopeTexture.Bind(gl.TEXTURE_2D)
	err = helpers.LoadTexture2D("textures/Cliffs0149_18_S.png")
	if err != nil {
		slopeTexture.Delete()
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
	colorTexture := gl.GenTexture()
	colorTexture.Bind(gl.TEXTURE_1D)
	err = helpers.LoadTexture1D("textures/gradient.png")
	if err != nil {
		colorTexture.Delete()
		fmt.Println(err)
	} else {
		textures = append(textures, colorTexture)
		gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}

	gl.ActiveTexture(gl.TEXTURE4)
	heightMapTexture := gl.GenTexture()
	heightMapTexture.Bind(gl.TEXTURE_2D)
	textures = append(textures, heightMapTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R16, heightMap.W, heightMap.H, 0, gl.RED, gl.FLOAT, heightMap.TexturePixels())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.ActiveTexture(gl.TEXTURE5)
	palmTexture := gl.GenTexture()
	palmTexture.Delete()
	err = helpers.LoadTexture2DWatched("textures/palme.png")
	if err != nil {
		palmTexture.Delete()
		fmt.Println("can't load palme.png")
		fmt.Println(err)
	} else {
		textures = append(textures, palmTexture)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	gl.ActiveTexture(gl.TEXTURE8)
	fireballTexture := gl.GenTexture()
	fireballTexture.Bind(gl.TEXTURE_2D)
	err = helpers.LoadTexture2D("textures/fireball.png")
	if err != nil {
		fireballTexture.Delete()
		panic("fireball.png")
	} else {
		textures = append(textures, fireballTexture)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	gl.ActiveTexture(gl.TEXTURE7)
	skybox := gl.GenTexture()
	skybox.Bind(gl.TEXTURE_CUBE_MAP)
	err = helpers.LoadTextureCube("textures/Above_The_Sea.jpg")
	if err != nil {
		skybox.Delete()
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

	gl.ActiveTexture(gl.TEXTURE6)

	ttf.Init()
	defer ttf.Quit()
	font, _ := ttf.OpenFont("fonts/Symbola.ttf", 64)
	color := sdl.Color{255, 255, 255, 255}
	surface := font.RenderText_Blended("Bla", color)
	defer surface.Free()
	textTexture := gl.GenTexture()

	textTexture.Bind(gl.TEXTURE_2D)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int(surface.W), int(surface.H), 0, gl.RGBA, gl.UNSIGNED_BYTE, uintptr(surface.Data()))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	surface.Data()

	gl.ActiveTexture(gl.TEXTURE0)

	return &Textures{textures}
}

func TexturesTest() {
	out := make([]int32, 1)
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, out)
	maxUnits := int(out[0])

	for i := 0; i < maxUnits; i++ {
		gl.ActiveTexture(gl.GLenum(gl.TEXTURE0 + i))

		gl.GetIntegerv(gl.TEXTURE_BINDING_1D, out)
		texture := gl.Texture(out[0])
		if texture != 0 {
			fmt.Println("unit: ", i, " texture1d: ", texture)
		}

		gl.GetIntegerv(gl.TEXTURE_BINDING_2D, out)
		texture = gl.Texture(out[0])
		if texture != 0 {
			fmt.Println("unit: ", i, " texture2d: ", texture)
		}

		gl.GetIntegerv(gl.TEXTURE_BINDING_RECTANGLE, out)
		texture = gl.Texture(out[0])
		if texture != 0 {
			fmt.Println("unit: ", i, " textureRect: ", texture)
		}

		gl.GetIntegerv(gl.TEXTURE_BINDING_CUBE_MAP, out)
		texture = gl.Texture(out[0])
		if texture != 0 {
			fmt.Println("unit: ", i, " textureCube: ", texture)
		}

		gl.GetIntegerv(gl.TEXTURE_BINDING_3D, out)
		texture = gl.Texture(out[0])
		if texture != 0 {
			fmt.Println("unit: ", i, " texture3d: ", texture)
		}
	}

	gl.ActiveTexture(gl.TEXTURE0)
}

func (this *Textures) Delete() {
	gl.DeleteTextures(this.Textures)
}
