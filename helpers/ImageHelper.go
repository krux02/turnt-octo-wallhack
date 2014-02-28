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

	imageData := image.NewRGBA(bounds)
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, bounds.Dx(), bounds.Dy(), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

	return texture, nil
}

func LoadTextureCube(name string) (gl.Texture, error) {
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

	size := m.Bounds().Max
	W := size.X / 4
	H := size.Y / 3
	top_rect := image.Rect(W, 0*H, 2*W, 1*H)
	bottom_rect := image.Rect(W, 2*H, 2*W, 3*H)
	left_rect := image.Rect(0*W, H, 1*W, 2*H)
	front_rect := image.Rect(1*W, H, 2*W, 2*H)
	right_rect := image.Rect(2*W, H, 3*W, 2*H)
	back_rect := image.Rect(3*W, H, 4*W, 2*H)
	bounds := image.Rect(0, 0, W, H)

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_CUBE_MAP)

	imageData := image.NewRGBA(bounds)
	draw.Draw(imageData, top_rect, m, image.ZP, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, bottom_rect, m, image.ZP, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, left_rect, m, image.ZP, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, right_rect, m, image.ZP, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, front_rect, m, image.ZP, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, back_rect, m, image.ZP, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

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
		panic(fmt.Sprintf("cant write to file %s %s", filename, err))
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	err = png.Encode(writer, img)
	if err != nil {
		panic(fmt.Sprintf("cant write to file %s %s", filename, err))
	}

	fmt.Println("file written", filename)
}

func SaveTexture(target gl.GLenum, level int, filename string) {
	params := make([]int32, 1)
	gl.GetTexLevelParameteriv(target, level, gl.TEXTURE_WIDTH, params)
	width := int(params[0])
	gl.GetTexLevelParameteriv(target, level, gl.TEXTURE_HEIGHT, params)
	height := int(params[0])

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	pixels := make([]uint8, width*height*4)
	gl.GetTexImage(target, level, gl.RGBA, gl.UNSIGNED_BYTE, pixels)

	// invert bottom/top
	stride := width * 4
	for i := 0; i < height; i++ {
		j := height - i - 1
		l := img.Pix[i*stride : (i+1)*stride]
		r := pixels[j*stride : (j+1)*stride]
		copy(l, r)
	}

	SaveImage(filename, img)
}
