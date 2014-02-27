package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"image"
	"image/color"
	"math"
)

// import "fmt"

type HeightMap struct {
	W, H    int
	Data    []float32
	MinTree []float32
	MaxTree []float32
}

func NewHeightMap(w, h int) (out *HeightMap) {
	if (w & (w - 1)) != 0 {
		panic("no pow of 2 size")
	}
	if (h & (h - 1)) != 0 {
		panic("no pow of 2 size")
	}
	if w != h {
		panic("width and height needs to be equal")
	}
	lodSize := 0
	for h2, w2 := h>>1, w>>1; h2 > 0 && w2 > 0; h2, w2 = h2>>1, w2>>1 {
		lodSize += h2 * w2
	}
	out = &HeightMap{w, h, make([]float32, w*h), make([]float32, lodSize), make([]float32, lodSize)}

	return
}

var gauss_factor = 1 / math.Sqrt(2*math.Pi)

func Gauss(x float32) float32 {
	return float32(math.Exp(float64(-x*x/2)) * gauss_factor)
}

func Gauss2f(v mgl.Vec2f) mgl.Vec2f {
	return mgl.Vec2f{Gauss(v[0]), Gauss(v[1])}
}

func Bump(center mgl.Vec2f, height float32) {

}

func (this *HeightMap) Get(x, y int) float32 {
	x = x & (this.W - 1)
	y = y & (this.H - 1)
	return this.Data[this.W*y+x]
}

func (this *HeightMap) Get2f(x, y float32) float32 {
	l := float32(math.Floor(float64(x)))
	r := float32(math.Floor(float64(x + 1)))
	b := float32(math.Floor(float64(y)))
	t := float32(math.Floor(float64(y + 1)))

	bl := this.Get(int(l), int(b))
	br := this.Get(int(r), int(b))
	tl := this.Get(int(l), int(t))
	tr := this.Get(int(r), int(t))

	bh := bl*(r-x) + br*(x-l)
	th := tl*(r-x) + tr*(x-l)

	h := bh*(t-y) + th*(y-b)

	return h
}

func (this *HeightMap) Set(x, y int, v float32) {
	if 0 <= x && x < this.W && 0 <= y && y < this.H {
		this.Data[this.W*y+x] = v
	}
}

func (m *HeightMap) flat(x, y int) int {
	return m.W*y + x
}

func (m *HeightMap) unflat(i int) (int, int) {
	return i % m.W, i / m.W
}

func (m *HeightMap) Normal(x int, y int) mgl.Vec3f {
	l := x - 1
	r := x + 1
	b := y - 1
	t := y + 1

	hi := m.Get(x, y)
	lh := m.Get(l, y) - hi
	rh := m.Get(r, y) - hi
	bh := m.Get(x, b) - hi
	th := m.Get(x, t) - hi

	v1 := mgl.Vec3f{1, 0, rh}.Normalize()
	v2 := mgl.Vec3f{0, 1, th}.Normalize()
	v3 := mgl.Vec3f{-1, 0, lh}.Normalize()
	v4 := mgl.Vec3f{0, -1, bh}.Normalize()

	n1 := v1.Cross(v2).Normalize()
	n2 := v2.Cross(v3).Normalize()
	n3 := v3.Cross(v4).Normalize()
	n4 := v4.Cross(v1).Normalize()

	return n1.Add(n2).Add(n3).Add(n4).Normalize()
}

func (m *HeightMap) Normalf(x float32, y float32) (n mgl.Vec3f) {
	x0 := int(math.Floor(float64(x)))
	x1 := x0 + 1
	y0 := int(math.Floor(float64(y)))
	y1 := y0 + 1

	n00 := m.Normal(x0, y0)
	n10 := m.Normal(x1, y0)
	n01 := m.Normal(x0, y1)
	n11 := m.Normal(x1, y1)

	w := x - float32(x0)
	h := y - float32(y0)

	n0 := n00.Mul(1 - w).Add(n10.Mul(w))
	n1 := n01.Mul(1 - w).Add(n11.Mul(w))

	n = n0.Mul(1 - h).Add(n1.Mul(h))
	return
}

func (m *HeightMap) TexturePixels() (pixels []float32) {
	pixels = make([]float32, m.W*m.H)
	minh, maxh := m.Bounds()
	for i, _ := range pixels {
		pixels[i] = (m.Data[i] - minh) / (maxh - minh)
	}
	return
}

func (m *HeightMap) Bounds() (float32, float32) {

	var min_h, max_h float32 = 10000.0, -10000.0

	for _, v := range m.Data {
		if v < min_h {
			min_h = v
		}
		if v > max_h {
			max_h = v
		}
	}

	return min_h, max_h
}

func (m *HeightMap) ExportImage() image.Image {

	minh, maxh := m.Bounds()
	diff := maxh - minh
	rect := image.Rect(0, 0, m.W, m.H)
	img := image.NewGray16(rect)

	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			h := (m.Get(x, y) - minh) / diff
			c := color.Gray16{uint16(h * math.MaxUint16)}
			img.SetGray16(x, y, c)
		}
	}

	return img
}

func NewHeightMapFromFile(filename string) *HeightMap {
	img, err := helpers.ReadToGray16(filename)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	m := NewHeightMap(w, h)

	src := img.Pix
	dst := m.Data
	for i := 0; i < len(src); i += 2 {
		dst[i>>1] = float32(int(src[i])|(int(src[i+1])<<8)) / math.MaxUint16
	}
	return m
}
