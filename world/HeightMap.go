package world

import "math/rand"
import "github.com/krux02/mathgl"
import "github.com/krux02/turnt-octo-wallhack/helpers"
import "github.com/go-gl/gl"
import "math"
import "image"
import "image/color"

// import "fmt"

type HeightMap struct {
	W, H int
	Data []float32
}

func NewHeightMap(w, h int) *HeightMap {
	if (w & (w - 1)) != 0 {
		panic("no pow of 2 size")
	}
	if (h & (h - 1)) != 0 {
		panic("no pow of 2 size")
	}
	return &HeightMap{w, h, make([]float32, w*h, w*h)}
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

func (m *HeightMap) DiamondSquare(factor float32) {
	w, h := m.W, m.H

	stepSize := w

	squares := func() {
		for i := 0; i < w; i += stepSize {
			for j := 0; j < h; j += stepSize {
				sum := m.Get(i, j)
				sum += m.Get(i+stepSize, j)
				sum += m.Get(i, j+stepSize)
				sum += m.Get(i+stepSize, j+stepSize)

				h := (sum / 4.0) + (rand.Float32()-0.5)*factor
				x := i + stepSize/2
				y := j + stepSize/2
				m.Set(x, y, h)
			}
		}
	}

	diamonds := func() {
		for i := 0; i <= w; i += stepSize {
			for j := 0; j <= h; j += stepSize {
				if ((i+j)/stepSize)%2 == 1 {
					sum := float32(0)
					count := float32(0.0)
					if i != 0 {
						sum += m.Get(i-stepSize, j)
						count += 1
					}
					if i < w-1 {
						sum += m.Get(i+stepSize, j)
						count += 1
					}
					if j != 0 {
						sum += m.Get(i, j-stepSize)
						count += 1
					}
					if j < h-1 {
						sum += m.Get(i, j+stepSize)
						count += 1
					}

					h := (sum / count) + (rand.Float32()-0.5)*factor
					m.Set(i, j, h)
				}
			}
		}
	}

	for stepSize > 0 {
		squares()

		stepSize /= 2
		if stepSize == 0 {
			break
		}
		diamonds()

		factor *= 0.5
	}
}

func (m *HeightMap) Normal(x int, y int) mathgl.Vec3f {
	l := x - 1
	r := x + 1
	b := y - 1
	t := y + 1

	hi := m.Get(x, y)
	lh := m.Get(l, y) - hi
	rh := m.Get(r, y) - hi
	bh := m.Get(x, b) - hi
	th := m.Get(x, t) - hi

	v1 := mathgl.Vec3f{1, 0, rh}.Normalize()
	v2 := mathgl.Vec3f{0, 1, th}.Normalize()
	v3 := mathgl.Vec3f{-1, 0, lh}.Normalize()
	v4 := mathgl.Vec3f{0, -1, bh}.Normalize()

	n1 := v1.Cross(v2).Normalize()
	n2 := v2.Cross(v3).Normalize()
	n3 := v3.Cross(v4).Normalize()
	n4 := v4.Cross(v1).Normalize()

	return n1.Add(n2).Add(n3).Add(n4).Normalize()
}

func (m *HeightMap) Triangulate() []int32 {
	w, h := m.W, m.H

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

		// v1h := m.Vertices()[v1].Vertex_ms[2]
		// v2h := m.Vertices()[v2].Vertex_ms[2]
		// v3h := m.Vertices()[v3].Vertex_ms[2]
		// v4h := m.Vertices()[v4].Vertex_ms[2]

		// vec1 := mathgl.Vec3f{0, v1h, 0}
		// vec2 := mathgl.Vec3f{1, 0.5*v2h + 0.5*v3h, 0}
		// vec3 := mathgl.Vec3f{2, v4h, 0}

		// cross1 := vec1.Sub(vec2).Cross(vec3.Sub(vec2))

		// vec1 = mathgl.Vec3f{0, v2h, 0}
		// vec2 = mathgl.Vec3f{1, 0.5*v1h + 0.5*v4h, 0}
		// vec3 = mathgl.Vec3f{2, v3h, 0}

		// cross2 := vec1.Sub(vec2).Cross(vec3.Sub(vec2))

		// fmt.Println(cross1, cross2)

		// if math.Abs(float64(cross1[2])) > math.Abs(float64(cross2[2])) {
		put(v1)
		put(v2)
		put(v3)

		put(v3)
		put(v2)
		put(v4)
		// } else {
		// put(v1)
		// put(v2)
		// put(v4)

		// put(v4)
		// put(v3)
		// put(v1)
		// }

	}

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			quad(i, j)
		}
	}

	return indices
}

func (m *HeightMap) Texture() gl.Texture {
	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	minh, maxh := m.Bounds()

	pixels := make([]float32, m.W*m.H)
	for i, _ := range pixels {
		pixels[i] = (m.Data[i] - minh) / (maxh - minh)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R16, m.W, m.H, 0, gl.RED, gl.FLOAT, pixels)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	return texture
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

const MaxShort = 65535

func (m *HeightMap) ExportImage() image.Image {

	minh, maxh := m.Bounds()
	diff := maxh - minh
	rect := image.Rect(0, 0, m.W, m.H)
	img := image.NewGray16(rect)

	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			h := (m.Get(x, y) - minh) / diff
			c := color.Gray16{uint16(h * MaxShort)}
			img.SetGray16(x, y, c)
		}
	}

	return img
}

func NewHeightMapFramFile(filename string) *HeightMap {
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
		dst[i>>1] = float32(int(src[i])|(int(src[i+1])<<8)) / MaxShort
	}
	return m
}
