#version 330 core

uniform sampler2D Image;

in vec4 vertexColor;
in vec2 v_texCoord;

out vec4 color;

void main() {
	color = vertexColor * texture(Image,v_texCoord);
}
