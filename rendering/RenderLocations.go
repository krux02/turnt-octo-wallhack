package rendering

import "github.com/go-gl/gl"

type RenderLocations struct {
	Vertex_os, Vertex_ms, Normal_ms, TexCoord gl.AttribLocation
	InstancePosition_ws                       gl.AttribLocation
	Position_ndc                              gl.AttribLocation
	Proj, View, Model, Rot2D                  gl.UniformLocation
	CameraPos_ws                              gl.UniformLocation
	PalmTree, Image, GroundTexture, Skybox    gl.UniformLocation
	ClippingPlane_cs, ClippingPlane_ws        gl.UniformLocation
	HeightMap, LowerBound, UpperBound         gl.UniformLocation
	Time                                      gl.UniformLocation
}
