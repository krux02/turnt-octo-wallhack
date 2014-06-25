package renderstuff

import "github.com/go-gl/gl"

type RenderLocations struct {
	Vertex_ws, Vertex_ms, Vertex_ndc gl.AttribLocation
	Normal_ms, TexCoord, Color       gl.AttribLocation
	InstancePosition_ws              gl.AttribLocation
	// Matrices
	Proj, View, Model, Rot2D gl.UniformLocation
	CameraPos_ws             gl.UniformLocation
	// Textures
	TextureTree, Image, TextureGround, TextureSkybox gl.UniformLocation
	TextureColorBand, TextureCliffs                  gl.UniformLocation
	ClippingPlane_ws                                 gl.UniformLocation
	TextureHeightMap, LowerBound, UpperBound         gl.UniformLocation
	Time, WaterHeight                                gl.UniformLocation
	ViewPortSize                                     gl.UniformLocation
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
