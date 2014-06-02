#version 330 core

uniform sampler2D TextureTree;

in vec2 v_texCoord;
in vec4 pos_ws;

out vec4 color;

void main() {

	color = texture(TextureTree, v_texCoord);
	if(color.a < 0.5)
		discard;
}		
