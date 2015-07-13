package rendering

import (
	//	"errors"
	"github.com/go-gl-legacy/gl"
)

type FrameBuffer struct {
	Framebuffer   gl.Framebuffer
	RenderTexture gl.Texture
	DepthTexture  gl.Texture
}

func NewFrameBuffer(width, height int) (this *FrameBuffer) {
	this = &FrameBuffer{gl.GenFramebuffer(), gl.GenTexture(), gl.GenTexture()}

	this.Resize(width, height)

	var outer gl.Framebuffer
	{
		params := []int32{0}
		gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, params)
		outer = gl.Framebuffer(params[0])
	}

	this.Framebuffer.Bind()
	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_RECTANGLE, this.RenderTexture, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.TEXTURE_RECTANGLE, this.DepthTexture, 0)
	outer.Bind()

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		this.Delete()
		panic("framebuffer incomplete")
	}

	return
}

func (this *FrameBuffer) Resize(width, height int) {
	//this.Framebuffer.Bind()
	// this.RenderTexture := gl.GenTexture()
	this.RenderTexture.Bind(gl.TEXTURE_RECTANGLE)
	gl.TexImage2D(gl.TEXTURE_RECTANGLE, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	this.RenderTexture.Unbind(gl.TEXTURE_RECTANGLE)

	// depthStencilTexture := gl.GenTexture()
	this.DepthTexture.Bind(gl.TEXTURE_RECTANGLE)
	gl.TexImage2D(gl.TEXTURE_RECTANGLE, 0, gl.DEPTH24_STENCIL8, width, height, 0, gl.DEPTH_STENCIL, gl.UNSIGNED_INT_24_8, nil)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_WRAP_S, gl.CLAMP)
	gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_WRAP_T, gl.CLAMP)
	this.DepthTexture.Unbind(gl.TEXTURE_RECTANGLE)
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
