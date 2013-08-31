package main

import (
	"code.google.com/p/freetype-go/freetype"
	"flag"
	"fmt"
	"github.com/go-gl/gl"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"math/rand"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "luxi-fonts/luxisr.ttf", "filename of the ttf font")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

var foregroundTexture gl.Texture

func RandomNoiseRectangle(w, h int) {

	data := make([]uint32, w*h)

	for i, _ := range data {
		data[i] = uint32(rand.Int31()) & 0x3fffffff
	}

	foregroundTexture.Bind(gl.TEXTURE_RECTANGLE)
	gl.TexImage2D(gl.TEXTURE_RECTANGLE, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg, bg := image.Black, image.Transparent

	wt, ht := 256, 32

	rgba := image.NewRGBA(image.Rect(0, 0, wt, ht))

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(font)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Draw the text.
	pt := freetype.Pt(1, 31)

	_, err = c.DrawString("Hallo Welt! XOXOXOXOXO", pt)
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			rgba.SetRGBA(i, j, color.RGBA{0xff, 0, 0, 0xff} )
			rgba.Set(wt-i-1, ht-j-1, color.RGBA{0, 0xff, 0, 0xff} )
		}
	}

	gl.TexSubImage2D(gl.TEXTURE_RECTANGLE, 0, 64, 64, wt, ht, gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)
}

func initTextures() func() {
	gl.ActiveTexture(gl.TEXTURE0)
	colorTexture, err := LoadTexture1D("gradient.png")
	if err != nil {
		fmt.Println(err)
	}

	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.ActiveTexture(gl.TEXTURE1)

	gl.ActiveTexture(gl.TEXTURE1)
	detailTexture, err := LoadTexture2D("GravelCobbleS.jpg")
	if err != nil {
		fmt.Println("cant load GravelCobble.jpg")
		fmt.Println(err)
	}

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)

	gl.ActiveTexture(gl.TEXTURE2)
	slopeTexture, err := LoadTexture2D("GravelCobble0003_2_S.jpg")
	if err != nil {
		fmt.Println("cant load GravelCobble0003_2_S.jpg")
		fmt.Println(err)
	}

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)

	gl.ActiveTexture(gl.TEXTURE3)
	foregroundTexture = gl.GenTexture()
	RandomNoiseRectangle(1024, 768)

	return func() {
		gl.DeleteTextures([]gl.Texture{colorTexture, detailTexture, slopeTexture, foregroundTexture})
	}
}
