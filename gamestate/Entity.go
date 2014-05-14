package gamestate

import (
	"encoding/json"
	mgl "github.com/Jragonmiris/mathgl"
	"io"
	"math"
)

type Entity struct {
	Position    mgl.Vec4f
	Orientation mgl.Quatf
}

func (e *Entity) Model() mgl.Mat4f {
	pos := e.Position
	return mgl.Translate3D(pos[0], pos[1], pos[2]).Mul4(e.Orientation.Mat4())
}

func (e *Entity) SetModel(m mgl.Mat4f) {
	m00 := m[0]
	m10 := m[1]
	m20 := m[2]

	m01 := m[4]
	m11 := m[5]
	m21 := m[6]

	m02 := m[8]
	m12 := m[9]
	m22 := m[10]

	m03 := m[12]
	m13 := m[13]
	m23 := m[14]

	qw := float32(math.Sqrt(float64(1+m00+m11+m22))) / 2
	qx := (m21 - m12) / (4 * qw)
	qy := (m02 - m20) / (4 * qw)
	qz := (m10 - m01) / (4 * qw)

	e.Orientation = mgl.Quatf{qw, mgl.Vec3f{qx, qy, qz}}
	e.Position = mgl.Vec4f{m03, m13, m23, 1}
}

func (e *Entity) View() mgl.Mat4f {
	return e.Model().Inv()
}

func (e *Entity) SetView(m mgl.Mat4f) {
	e.SetModel(m.Inv())
}

func (e *Entity) Save(writer io.Writer) {
	encoder := json.NewEncoder(writer)
	m := map[string][4]float32{
		"Position":    [4]float32(e.Position),
		"Orientation": [4]float32{e.Orientation.W, e.Orientation.V[0], e.Orientation.V[1], e.Orientation.V[2]},
	}
	encoder.Encode(m)
}

func (e *Entity) Load(reader io.Reader) {
	decoder := json.NewDecoder(reader)
	m := map[string][4]float32{}
	decoder.Decode(m)
	e.Position = mgl.Vec4f(m["Position"])
	q := m["Orientation"]
	e.Orientation.W = q[0]
	e.Orientation.V = mgl.Vec3f{q[1], q[2], q[3]}
}
