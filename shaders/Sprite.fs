#version 330 core

uniform sampler2D PalmTree;

in vec2 v_texCoord;

out vec4 color;

void main() {
	color = texture(PalmTree, v_texCoord);
}

