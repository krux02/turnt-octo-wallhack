package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
)

type Water struct {
	W, H       int
	LowerBound mgl.Vec3f
	UpperBound mgl.Vec3f
}

func (this *Water) GetModel() mgl.Mat4f {
	return mgl.Ident4f()
}

func (this *Water) GetMesh() IMesh {
	return this
}

func (this *Water) Vertices() interface{} {
	return WaterVertices(this.W, this.H)
}

func (this *Water) Indices() interface{} {
	return TriangulationIndices(this.W, this.H)
}

func (this *Water) InstanceData() interface{} {
	return nil
}

func (this *Water) Mode() Mode {
	return Triangles
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
