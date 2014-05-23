package rendering

import "github.com/go-gl/gl"

type RenderLocations struct {
	Vertex_ws, Vertex_ms, Vertex_ndc gl.AttribLocation
	Normal_ms, TexCoord, Color       gl.AttribLocation
	InstancePosition_ws              gl.AttribLocation
	// Matrices
	Proj, View, Model, Rot2D gl.UniformLocation
	CameraPos_ws             gl.UniformLocation
	// Textures
	PalmTree, Image, GroundTexture, Skybox gl.UniformLocation
	ClippingPlane_ws                       gl.UniformLocation
	HeightMap, LowerBound, UpperBound      gl.UniformLocation
	Time                                   gl.UniformLocation
}
