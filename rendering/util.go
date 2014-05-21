package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
)

func glMat4(mat *mgl.Mat4f) *[16]float32 {
	return (*[16]float32)(mat)
}

func glMat3(mat *mgl.Mat3f) *[9]float32 {
	return (*[9]float32)(mat)
}

type RenderData struct {
	VAO                gl.VertexArray
	InstanceDataBuffer gl.Buffer
	NumInstances       int
	Indices            gl.Buffer
	Vertices           gl.Buffer
	Numverts           int
}

func (this *RenderData) Delete() {
	this.VAO.Delete()
	this.Indices.Delete()
	this.Vertices.Delete()
	*this = RenderData{}
}
