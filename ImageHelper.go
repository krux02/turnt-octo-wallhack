package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
)

func LoadTexture1D(name string) (gl.Texture, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer file.Close()

	bounds := m.Bounds()

	imageData := image.NewRGBA(m.Bounds())
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_1D)

	if bounds.Dx() != 1 && bounds.Dy() != 1 {
		panic(fmt.Sprintf("image %s must be one dimensionnal it is %s", name, bounds.String()))
	}
	width := bounds.Dx() * bounds.Dy()

	glError := gl.GetError()
	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA8, width, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	glError = gl.GetError()

	if gl.NO_ERROR != glError {
		texture.Unbind(gl.TEXTURE_1D)
		fmt.Println("foo")
		return 0, GLerror(glError)
	}

	return texture, nil
}

func LoadTexture2D(name string) (gl.Texture, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer file.Close()

	bounds := m.Bounds()

	imageData := image.NewRGBA(m.Bounds())
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, bounds.Dx(), bounds.Dy(), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	glError := gl.GetError()

	if gl.NO_ERROR != glError {
		texture.Unbind(gl.TEXTURE_1D)
		return 0, GLerror(glError)
	}

	return texture, nil
}
