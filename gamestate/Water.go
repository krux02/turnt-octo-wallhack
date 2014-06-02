package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
)

type Water struct {
	mesh       *renderstuff.Mesh
	W, H       int
	LowerBound mgl.Vec3f
	UpperBound mgl.Vec3f
	Height     float32
}

func (this *Water) Model() mgl.Mat4f {
	return mgl.Ident4f()
}

func (this *Water) Mesh() *renderstuff.Mesh {
	if this.mesh == nil {
		this.mesh = new(renderstuff.Mesh)
		this.mesh.Vertices = WaterVertices(this.W, this.H)
		this.mesh.Indices = TriangulationIndices(this.W, this.H)
		this.mesh.Mode = renderstuff.Triangles
	}
	return this.mesh
}

type WaterVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

func WaterVertices(W, H int) []WaterVertex {
	vertices := make([]WaterVertex, (W+1)*(H+1))
	i := 0
	for y := 0; y <= H; y++ {
		for x := 0; x <= W; x++ {
			pos := mgl.Vec3f{float32(x), float32(y), 0}
			nor := mgl.Vec3f{0, 0, 1}
			vertices[i] = WaterVertex{pos, nor}
			i += 1
		}
	}
	return vertices
}
