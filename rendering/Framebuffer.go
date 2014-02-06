package rendering

import (
	//	"errors"
	"github.com/go-gl/gl"
)

type FrameBuffer struct {
	Framebuffer   gl.Framebuffer
	RenderTexture gl.Texture
	DepthTexture  gl.Texture
}

func NewFrameBuffer() *FrameBuffer {
	framebuffer := gl.GenFramebuffer()

	var outer gl.Framebuffer
	{
		params := []int32{0}
		gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, params)
		outer = gl.Framebuffer(params[0])
	}

	framebuffer.Bind()
	renderTexture := gl.GenTexture()
	renderTexture.Bind(gl.TEXTURE_RECTANGLE)
	gl.TexImage2D(gl.TEXTURE_RECTANGLE, 0, gl.RGBA, 1024, 768, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_RECTANGLE, renderTexture, 0)

	depthTexture := gl.GenTexture()
	depthTexture.Bind(gl.TEXTURE_RECTANGLE)
	gl.TexImage2D(gl.TEXTURE_RECTANGLE, 0, gl.DEPTH_COMPONENT24, 1024, 768, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_WRAP_S, gl.CLAMP)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_WRAP_T, gl.CLAMP)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_RECTANGLE, depthTexture, 0)

	/*
		gl.TexImage2D(gl.TEXTURE_DEPTH, 0, internalformat, width, height, border, format, typ, pixels)
		depthrenderbuffer.Bind()
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, 1024, 768)
		depthrenderbuffer.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER)
	*/

	//DrawBuffers := []gl.GLenum{gl.COLOR_ATTACHMENT0}
	//gl.DrawBuffers(len(DrawBuffers), DrawBuffers)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("framebuffer incomplete")
	}

	outer.Bind()
	return &FrameBuffer{framebuffer, renderTexture, depthTexture}
}

func (this *FrameBuffer) Delete() {
	this.Framebuffer.Delete()
	this.RenderTexture.Delete()
	this.DepthTexture.Delete()
	*this = FrameBuffer{}
}

func (this *FrameBuffer) Bind() {
	this.Framebuffer.Bind()
}

func (this *FrameBuffer) Unbind() {
	this.Framebuffer.Unbind()
}
