package helpers

import "github.com/go-gl/gl"

type textureManager struct {
	texture gl.Texture
	target  gl.GLenum
}

func (this *textureManager) Update(filename string) {
	outer := gl.Texture(getEnum(BindingMap[this.target]))
	this.texture.Bind(this.target)
	LoadTexture(filename, this.target)
	gl.GenerateMipmap(this.target)
	outer.Bind(this.target)
}

func getEnum(enum gl.GLenum) gl.GLenum {
	i, _, _, _ := gl.GetInteger4(gl.TEXTURE_BINDING_2D)
	return gl.GLenum(i)
}
