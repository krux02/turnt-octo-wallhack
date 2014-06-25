package rendering

import (
	"fmt"
	"github.com/krux02/libovr"
	mgl "github.com/krux02/mathgl/mgl32"
)

type OvrStuff struct {
	Hmd                  *ovr.Hmd
	HmdDesc              ovr.HmdDesc
	Proj                 [2]mgl.Mat4
	EyeRenderDesc        [2]ovr.EyeRenderDesc
	ViewportsFramebuffer [2]Viewport
	ViewportsScreen      [2]Viewport
	Textures             [2]ovr.GLTexture
}

func (this *OvrStuff) Init(w, h int, fb *FrameBuffer) *OvrStuff {
	this.Hmd = ovr.HmdCreate(0)
	if this.Hmd == nil {
		fmt.Println("cant create Hmd device")
		this.Hmd = ovr.HmdCreateDebug(ovr.Hmd_DK1)
	}
	this.HmdDesc = this.Hmd.GetDesc()
	fmt.Printf("%+v\n", this.HmdDesc)
	eyeFovIn := this.HmdDesc.DefaultEyeFov

	var apiConfig ovr.GLConfig
	apiConfig.OGL().Header.API = ovr.RenderAPI_OpenGL
	apiConfig.OGL().Header.Multisample = 1
	apiConfig.OGL().Header.RTSize = ovr.Sizei{int32(w), int32(h)}
	distortionCaps := ovr.DistortionCap_Chromatic
	var ok bool
	this.EyeRenderDesc, ok = this.Hmd.ConfigureRendering(apiConfig.Config(), distortionCaps, eyeFovIn)
	if !ok {
		panic("configure rendering failed")
	} else {
		fmt.Printf("%+v\n", this.EyeRenderDesc)
	}

	// ovr is row major
	this.Proj[0] = mgl.Mat4(ovr.MatrixProjection(eyeFovIn[0], 0.3, 1000, true).FlatArray()).Transpose()
	this.Proj[1] = mgl.Mat4(ovr.MatrixProjection(eyeFovIn[1], 0.3, 1000, true).FlatArray()).Transpose()

	this.ViewportsFramebuffer[0] = Viewport{0, 0, 960, 1080}
	this.ViewportsFramebuffer[1] = Viewport{960, 0, 960, 1080}
	this.ViewportsScreen[0] = Viewport{0, 0, w / 2, h}
	this.ViewportsScreen[1] = Viewport{w / 2, 0, w / 2, h}

	for eye := ovr.Eye_Left; eye < ovr.Eye_Count; eye++ {
		textureData := this.Textures[eye].OGL()
		textureData.Header.RenderViewport = this.ViewportsFramebuffer[eye].ToOvrRecti()
		textureData.Header.API = ovr.RenderAPI_OpenGL
		textureData.Header.TextureSize = ovr.Sizei{1920, 1080}
		textureData.TexId = uint32(fb.RenderTexture)
	}

	return this
}
