package rendering

import (
	//	"errors"
	"github.com/go-gl/gl"
)

type FrameBuffer struct {
	W, H          int
	Framebuffer   gl.Framebuffer
	RenderTexture gl.Texture
	DepthTexture  gl.Texture
}

const target = gl.TEXTURE_2D

func NewFrameBuffer(width, height int) (this *FrameBuffer) {
	this = &FrameBuffer{width, height, gl.GenFramebuffer(), gl.GenTexture(), gl.GenTexture()}

	this.Resize(width, height)

	var outer gl.Framebuffer
	{
		params := []int32{0}
		gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, params)
		outer = gl.Framebuffer(params[0])
	}

	this.Framebuffer.Bind()
	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, target, this.RenderTexture, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, target, this.DepthTexture, 0)
	outer.Bind()

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		this.Delete()
		panic("framebuffer incomplete")
	}

	return
}

func (this *FrameBuffer) Resize(width, height int) {
	this.W, this.H = width, height
	//this.Framebuffer.Bind()
	// this.RenderTexture := gl.GenTexture()
	this.RenderTexture.Bind(target)
	gl.TexImage2D(target, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	this.RenderTexture.Unbind(target)

	// depthStencilTexture := gl.GenTexture()
	this.DepthTexture.Bind(target)
	gl.TexImage2D(target, 0, gl.DEPTH24_STENCIL8, width, height, 0, gl.DEPTH_STENCIL, gl.UNSIGNED_INT_24_8, nil)
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, gl.CLAMP)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_T, gl.CLAMP)
	this.DepthTexture.Unbind(target)
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
