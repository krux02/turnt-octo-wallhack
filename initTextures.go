package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

func initTextures() func() {
	textures := make([]gl.Texture, 0, 6)

	gl.ActiveTexture(gl.TEXTURE0)
	colorTexture, err := helpers.LoadTexture1D("textures/gradient.png")
	if err != nil {
		fmt.Println(err)
	} else {
		textures = append(textures, colorTexture)
		gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.ActiveTexture(gl.TEXTURE1)
	}

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
	slopeTexture, err := helpers.LoadTexture2D("textures/GravelCobble0003_2_S.jpg")
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

	// texture 4 is filled by the heightmap
	// texture 5 is used by tw

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

	return func() {
		gl.DeleteTextures(textures)
	}
}
