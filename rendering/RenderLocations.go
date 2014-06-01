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
	ColorBand, Texture, Slope              gl.UniformLocation
	ClippingPlane_ws                       gl.UniformLocation
	HeightMap, LowerBound, UpperBound      gl.UniformLocation
	Time, WaterHeight                      gl.UniformLocation
}

type RenderData struct {
	VAO          gl.VertexArray
	InstanceData gl.Buffer
	NumInstances int
	Indices      gl.Buffer
	IndexType    gl.GLenum
	Vertices     gl.Buffer
	Numverts     int
	Mode         gl.GLenum
}

func (this *RenderData) Delete() {
	this.VAO.Delete()
	this.InstanceData.Delete()
	this.Indices.Delete()
	this.Vertices.Delete()
	*this = RenderData{}
}
