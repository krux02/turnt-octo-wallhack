package gamestate

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/math32"
	"image"
	"image/color"
	"math"
)

// import "fmt"

type HeightMap struct {
	W, H int
	Data []float32
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
	return &HeightMap{w, h, make([]float32, w*h)}
}

func Gauss2f(v mgl.Vec2f) mgl.Vec2f {
	return mgl.Vec2f{math32.Gauss(v[0]), math32.Gauss(v[1])}
}

func Bump(center mgl.Vec2f, height float32) {

}

func (this *HeightMap) InRange(x, y int) bool {
	return 0 <= x && x < this.W && 0 <= y && y < this.H
}

func (this *HeightMap) Get(x, y int) float32 {
	x = x & (this.W - 1)
	y = y & (this.H - 1)
	return this.Data[this.W*y+x]
}

func (this *HeightMap) Get2f(x, y float32) float32 {
	l := math32.Floor(x)
	r := math32.Floor(x + 1)
	b := math32.Floor(y)
	t := math32.Floor(y + 1)

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
	x0 := int(math32.Floor(x))
	x1 := x0 + 1
	y0 := int(math32.Floor(y))
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

func (m *HeightMap) RayCast(pos mgl.Vec3f, dir mgl.Vec3f) (out mgl.Vec3f, hit bool) {
	var visit = func(x, y int) bool {
		fmt.Println("visit ", x, y)
		if !m.InRange(x, y) {
			fmt.Println("out of range")
			return false
		}

		p1 := mgl.Vec3f{float32(x), float32(y), m.Get(x, y)}
		p2 := mgl.Vec3f{float32(x + 1), float32(y), m.Get(x+1, y)}
		p3 := mgl.Vec3f{float32(x + 1), float32(y + 1), m.Get(x+1, y+1)}
		p4 := mgl.Vec3f{float32(x), float32(y + 1), m.Get(x, y+1)}

		var factor float32

		factor, hit = triangle_intersection(p1, p2, p3, pos, dir)
		if hit {
			out = dir.Mul(factor).Add(pos)
			fmt.Println("hit 1 ", out, hit)
			return false
		}
		factor, hit = triangle_intersection(p3, p4, p1, pos, dir)
		if hit {
			out = dir.Mul(factor).Add(pos)
			fmt.Println("hit2 ", out, hit)
			return false
		}

		return true
	}

	x0 := pos[0]
	y0 := pos[1]
	x1 := pos[0] + dir[0]
	y1 := pos[1] + dir[1]

	raytrace(x0, y0, x1, y1, visit)
	return
}

func triangle_intersection(V1, V2, V3, O, D mgl.Vec3f) (out float32, hit bool) {
	const EPSILON = 0.000001
	var e1, e2 mgl.Vec3f //Edge1, Edge2
	var P, Q, T mgl.Vec3f
	var det, inv_det, u, v, t float32

	//Find vectors for two edges sharing V1
	e1 = V2.Sub(V1)
	e2 = V3.Sub(V1)
	//Begin calculating determinant - also used to calculate u parameter
	P = D.Cross(e2)
	//if determinant is near zero, ray lies in plane of triangle
	det = e1.Dot(P)
	//NOT CULLING
	if det > -EPSILON && det < EPSILON {
		return 0, false
	}
	inv_det = 1 / det

	//calculate distance from V1 to ray origin
	T = O.Sub(V1)

	//Calculate u parameter and test bound
	u = T.Dot(P) * inv_det
	//The intersection lies outside of the triangle
	if u < 0 || u > 1 {
		return 0, false
	}

	//Prepare to test v parameter
	Q = T.Cross(e1)

	//Calculate V parameter and test bound
	v = D.Dot(Q) * inv_det
	//The intersection lies outside of the triangle
	if v < 0 || u+v > 1 {
		return 0, false
	}

	t = e2.Dot(Q) * inv_det

	if t > EPSILON { //ray intersection
		return t, true
	}

	// No hit, no win
	return 0, false
}

func raytrace(x0, y0, x1, y1 float32, visit func(x, y int) bool) {

	dx := math32.Abs(x1 - x0)
	dy := math32.Abs(y1 - y0)

	x := int(math32.Floor(x0))
	y := int(math32.Floor(y0))

	var x_inc, y_inc int
	var err float32

	if dx == 0 {
		x_inc = 0
		err = math32.Inf(1)
	} else if x1 > x0 {
		x_inc = 1
		err = (math32.Floor(x0) + 1 - x0) * dy
	} else {
		x_inc = -1
		err = (x0 - math32.Floor(x0)) * dy
	}

	if dy == 0 {
		y_inc = 0
		err -= math32.Inf(1)
	} else if y1 > y0 {
		y_inc = 1
		err -= (math32.Floor(y0) + 1 - y0) * dx
	} else {
		y_inc = -1
		err -= (y0 - math32.Floor(y0)) * dx
	}

	for visit(x, y) {
		if err > 0 {
			y += y_inc
			err -= dx
		} else {
			x += x_inc
			err += dy
		}
	}
}
