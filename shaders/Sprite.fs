#version 330 core

uniform sampler2D PalmTree;

in vec2 v_texCoord;

out vec4 color;

void main() {
	//color = texture(PalmTree, v_texCoord);
	color = vec4(1,0,1,1);
}

