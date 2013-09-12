package helpers

import (
	"bufio"
	"fmt"
	"github.com/go-gl/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

func LoadTexture1D(name string) (gl.Texture, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(name, err)
		return 0, err
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(name, err)
		return 0, err
	}

	bounds := m.Bounds()

	imageData := image.NewRGBA(m.Bounds())
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	texture := gl.GenTexture()
	

	if bounds.Dx() != 1 && bounds.Dy() != 1 {
		panic(fmt.Sprintf("image %s must be one dimensionnal it is %s", name, bounds.String()))
	}
	width := bounds.Dx() * bounds.Dy()

	texture.Bind(gl.TEXTURE_1D)
	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA8, width, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

	return texture, nil
}

func LoadTexture2D(name string) (gl.Texture, error) {
	file, err := os.Open(name)
	if err != nil {

		fmt.Println(name, err)
		return 0, err
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(name, err)
		return 0, err
	}

	bounds := m.Bounds()

	imageData := image.NewRGBA(m.Bounds())
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, bounds.Dx(), bounds.Dy(), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

	return texture, nil
}

func ReadToGray16(filename string) (*image.Gray16, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("cant open file")
		fmt.Println(filename, err)
		return nil, err
	}
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("cant decode file")
		fmt.Println(filename, err)
		return nil, err
	}

	bounds := m.Bounds()

	imageData := image.NewGray16(m.Bounds())
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	return imageData, nil
}

func SaveImage(filename string, img image.Image) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("cant write to file", filename)
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	err = png.Encode(writer, img)
	if err != nil {
		fmt.Println("cant write to file", filename)
		fmt.Println(err)
		return
	}

	fmt.Println("file written", filename)
}
