package rendering

import mgl "github.com/Jragonmiris/mathgl"

type HeightMapVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}

type TreeVertex struct {
	Vertex_ms mgl.Vec4f
	TexCoord  mgl.Vec2f
}

type WaterVertex struct {
	Vertex_ms, Normal_ms mgl.Vec3f
}
