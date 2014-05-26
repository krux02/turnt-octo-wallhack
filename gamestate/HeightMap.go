package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/math32"
	"github.com/krux02/turnt-octo-wallhack/mathint"
	"image"
	"image/color"
	"math"
)

// import "fmt"

type HeightMap struct {
	W, H       int
	Data       []float32
	HasChanges bool
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
	return &HeightMap{
		W:          w,
		H:          h,
		Data:       make([]float32, w*h),
		HasChanges: true,
	}
}

func Gauss2fv(v mgl.Vec2f) float32 {
	return math32.Gauss(v.Len())
}

func (this *HeightMap) Bump(center mgl.Vec2f, height float32) {
	minX := mathint.Max(math32.RoundInt(center[0])-5, 0)
	maxX := mathint.Min(math32.RoundInt(center[0])+5, this.W)
	minY := mathint.Max(math32.RoundInt(center[1])-5, 0)
	maxY := mathint.Min(math32.RoundInt(center[1])+5, this.H)

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			v := mgl.Vec2f{float32(x), float32(y)}
			v = v.Sub(center)
			bump := math32.Gauss(v.Len())
			h := this.Get(x, y)
			this.Set(x, y, h+bump)
		}
	}
	this.HasChanges = true
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

func (m *HeightMap) Normal2f(x float32, y float32) (n mgl.Vec3f) {
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
	var min_h, max_h float32 = math.MaxFloat32, -math.MaxFloat32
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

func (m *HeightMap) RayCast(pos mgl.Vec3f, dir mgl.Vec3f) (factor float32, hit bool) {

	var visit = func(x, y int) bool {
		if !m.InRange(x, y) {
			return false
		}

		p1 := mgl.Vec3f{float32(x), float32(y), m.Get(x, y)}
		p2 := mgl.Vec3f{float32(x + 1), float32(y), m.Get(x+1, y)}
		p3 := mgl.Vec3f{float32(x + 1), float32(y + 1), m.Get(x+1, y+1)}
		p4 := mgl.Vec3f{float32(x), float32(y + 1), m.Get(x, y+1)}

		factor, hit = helpers.TriangleIntersection(p4, p1, p2, pos, dir)
		if hit {
			return false
		}
		factor, hit = helpers.TriangleIntersection(p2, p3, p4, pos, dir)
		if hit {
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

type HeightMapVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

func (this *HeightMap) CreateVertexArray() (vertices interface{}, indices interface{}) {
	vertices = Vertices(this)
	indices = TriangulationIndices(this.W, this.H)
	return vertices, indices
}

func (this *HeightMap) GetMesh() IMesh {
	return this
}

func (this *HeightMap) GetModel() mgl.Mat4f {
	return mgl.Ident4f()
}

func TriangulationIndices(w, h int) []int32 {
	indexCount := 6 * w * h
	indices := make([]int32, indexCount, indexCount)

	i := 0
	put := func(v int) {
		indices[i] = int32(v)
		i += 1
	}

	flat := func(x, y int) int {
		return (w+1)*y + x
	}

	quad := func(x, y int) {
		v1 := flat(x, y)
		v2 := flat(x+1, y)
		v3 := flat(x, y+1)
		v4 := flat(x+1, y+1)

		put(v1)
		put(v2)
		put(v3)

		put(v3)
		put(v2)
		put(v4)
	}

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			quad(i, j)
		}
	}

	return indices
}

func Vertices(m *HeightMap) []HeightMapVertex {
	vertices := make([]HeightMapVertex, (m.W+1)*(m.H+1))
	i := 0
	for y := 0; y <= m.H; y++ {
		for x := 0; x <= m.W; x++ {
			h := m.Get(x, y)
			pos := mgl.Vec3f{float32(x), float32(y), h}
			nor := m.Normal(x, y)
			vertices[i] = HeightMapVertex{pos, nor}
			i += 1
		}
	}
	return vertices
}
