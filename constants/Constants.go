package constants

const (
	_ = iota // TEXTURE0 is only used for temporarly bound textures
	TextureGround
	TextureCliffs
	TextureColorBand
	TextureHeightMap
	TextureTree
	TextureFireBall
	TextureSkybox
	TextureFont
)

var Texture = map[string]int{
	"TextureGround":    TextureGround,
	"TextureCliffs":    TextureCliffs,
	"TextureColorBand": TextureColorBand,
	"TextureHeightMap": TextureHeightMap,
	"TextureTree":      TextureTree,
	"TextureFireBall":  TextureFireBall,
	"TextureSkybox":    TextureSkybox,
	"TextureFont":      TextureFont,
}
