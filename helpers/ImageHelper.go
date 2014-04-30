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

func LoadTexture1D(name string) error {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(name, err)
		return err
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(name, err)
		return err
	}

	bounds := m.Bounds()

	imageData := image.NewRGBA(m.Bounds())
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	if bounds.Dx() != 1 && bounds.Dy() != 1 {
		panic(fmt.Sprintf("image %s must be one dimensionnal it is %s", name, bounds.String()))
	}
	width := bounds.Dx() * bounds.Dy()

	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA8, width, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

	return nil
}

type myTexture struct {
	texture gl.Texture
	binding gl.GLenum
}

var textureUnitMapping = make(map[string]myTexture)
var filechanges = make(chan string)

func LoadTexture2DWatched(name string) error {
	i, _, _, _ := gl.GetInteger4(gl.TEXTURE_BINDING_2D)
	textureUnitMapping[name] = myTexture{gl.Texture(i), gl.TEXTURE_2D}
	go Listen(name, filechanges)
	return LoadTexture2D(name)
}

func getEnum(enum gl.GLenum) gl.GLenum {
	i, _, _, _ := gl.GetInteger4(gl.TEXTURE_BINDING_2D)
	return gl.GLenum(i)
}

var BindingMap = map[gl.GLenum]gl.GLenum{
	gl.TEXTURE_1D:                   gl.TEXTURE_BINDING_1D,
	gl.TEXTURE_1D_ARRAY:             gl.TEXTURE_BINDING_1D_ARRAY,
	gl.TEXTURE_2D:                   gl.TEXTURE_BINDING_2D,
	gl.TEXTURE_2D_ARRAY:             gl.TEXTURE_BINDING_2D_ARRAY,
	gl.TEXTURE_2D_MULTISAMPLE:       gl.TEXTURE_BINDING_2D_MULTISAMPLE,
	gl.TEXTURE_2D_MULTISAMPLE_ARRAY: gl.TEXTURE_BINDING_2D_MULTISAMPLE_ARRAY,
	gl.TEXTURE_3D:                   gl.TEXTURE_BINDING_3D,
	gl.TEXTURE_CUBE_MAP:             gl.TEXTURE_BINDING_CUBE_MAP,
	gl.TEXTURE_CUBE_MAP_ARRAY:       gl.TEXTURE_BINDING_CUBE_MAP_ARRAY,
	gl.TEXTURE_RECTANGLE:            gl.TEXTURE_BINDING_RECTANGLE,
}

func UpdateTextures() {
	b := true
	for b {
		select {
		case name := <-filechanges:
			texUnit := textureUnitMapping[name]
			outer := gl.Texture(getEnum(BindingMap[texUnit.binding]))
			texUnit.texture.Bind(texUnit.binding)
			LoadTexture2D(name)
			gl.GenerateMipmap(gl.TEXTURE_2D)
			outer.Bind(texUnit.binding)
		default:
			b = false
		}
	}
}

func LoadTexture2D(name string) error {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(name, err)
		return err
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(name, err)
		return err
	}

	bounds := m.Bounds()

	imageData := image.NewRGBA(bounds)
	draw.Draw(imageData, bounds, m, image.ZP, draw.Src)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, bounds.Dx(), bounds.Dy(), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

	return nil
}

func LoadTextureCube(name string) error {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(name, err)
		return err
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(name, err)
		return err
	}

	size := m.Bounds().Max
	W := size.X / 4
	H := size.Y / 3
	top_rect := image.Point{W, 0 * H}
	bottom_rect := image.Point{W, 2 * H}
	left_rect := image.Point{0 * W, H}
	front_rect := image.Point{1 * W, H}
	right_rect := image.Point{2 * W, H}
	back_rect := image.Point{3 * W, H}
	bounds := image.Rect(0, 0, W, H)

	imageData := image.NewRGBA(bounds)
	draw.Draw(imageData, bounds, m, top_rect, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, bounds, m, bottom_rect, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, bounds, m, left_rect, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, bounds, m, right_rect, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, bounds, m, back_rect, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)
	draw.Draw(imageData, bounds, m, front_rect, draw.Src)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pix)

	return nil
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

	img := image.NewRGBA(image.Rect(0, 0, width, height))
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

	go SaveImage(filename, img)
}
