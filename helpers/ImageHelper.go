package helpers

import (
	"bufio"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

const pixelFormat = sdl.PIXELFORMAT_ABGR8888

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

func loadSdlSurface(filename string) *sdl.Surface {
	surface := img.Load(filename)
	if surface == nil {
		panic(sdl.GetError())
	}
	defer surface.Free()
	return surface.ConvertFormat(pixelFormat, 0)
}

func LoadTexture(filename string, target gl.GLenum) {
	switch target {
	case gl.TEXTURE_1D:
		LoadTexture1D(filename)
	case gl.TEXTURE_2D:
		LoadTexture2D(filename)
	case gl.TEXTURE_RECTANGLE:
		LoadTextureRect(filename)
	case gl.TEXTURE_CUBE_MAP:
		LoadTextureCube(filename)
	}
}

func LoadTexture1D(filename string) {
	surface := loadSdlSurface(filename)
	defer surface.Free()

	W, H := int(surface.W), int(surface.H)
	if W != 1 && H != 1 {
		panic(fmt.Sprintf("image %s must be one dimensionnal it is %dx%d", filename, W, H))
	}

	width := W * H
	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA8, width, 0, gl.RGBA, gl.UNSIGNED_BYTE, surface.Pixels())
}

func LoadTexture2D(filename string) {
	surface := loadSdlSurface(filename)
	defer surface.Free()

	W, H := int(surface.W), int(surface.H)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, surface.Pixels())
}

func LoadTextureRect(filename string) {
	surface := loadSdlSurface(filename)
	defer surface.Free()
	W, H := int(surface.W), int(surface.H)
	gl.TexImage2D(gl.TEXTURE_RECTANGLE, 0, gl.RGBA8, W, H, 0, gl.RGBA, gl.UNSIGNED_BYTE, surface.Pixels())
}

func LoadTextureCube(filename string) {
	surface := loadSdlSurface(filename)
	defer surface.Free()

	W := surface.W / 4
	H := surface.H / 3

	top_rect := sdl.Rect{W, 0 * H, W, H}
	bottom_rect := sdl.Rect{W, 2 * H, W, H}
	left_rect := sdl.Rect{0 * W, H, W, H}
	front_rect := sdl.Rect{1 * W, H, W, H}
	right_rect := sdl.Rect{2 * W, H, W, H}
	back_rect := sdl.Rect{3 * W, H, W, H}
	bounds := sdl.Rect{0, 0, W, H}
	imageData := sdl.CreateRGBSurface(0, W, H, 32, surface.Format.Rmask, surface.Format.Gmask, surface.Format.Bmask, surface.Format.Amask)

	surface.Blit(&top_rect, imageData, &bounds)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, gl.RGBA, int(W), int(H), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pixels())
	surface.Blit(&bottom_rect, imageData, &bounds)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, gl.RGBA, int(W), int(H), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pixels())
	surface.Blit(&left_rect, imageData, &bounds)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, gl.RGBA, int(W), int(H), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pixels())
	surface.Blit(&right_rect, imageData, &bounds)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.RGBA, int(W), int(H), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pixels())
	surface.Blit(&back_rect, imageData, &bounds)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, gl.RGBA, int(W), int(H), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pixels())
	surface.Blit(&front_rect, imageData, &bounds)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, gl.RGBA, int(W), int(H), 0, gl.RGBA, gl.UNSIGNED_BYTE, imageData.Pixels())
}

func LoadTextureWatched(filename string, target gl.GLenum) {
	i, _, _, _ := gl.GetInteger4(BindingMap[target])
	Manage(&textureManager{gl.Texture(i), target}, filename)
	LoadTexture(filename, target)
}

func LoadTexture1DWatched(filename string) {
	LoadTextureWatched(filename, gl.TEXTURE_1D)
}

func LoadTexture2DWatched(filename string) {
	LoadTextureWatched(filename, gl.TEXTURE_2D)
}

func LoadTextureRectWatched(filename string) {
	LoadTextureWatched(filename, gl.TEXTURE_RECTANGLE)
}

func LoadTextureCubeWatched(filename string) {
	LoadTextureWatched(filename, gl.TEXTURE_CUBE_MAP)
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
