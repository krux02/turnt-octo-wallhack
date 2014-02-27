#version 330 core

uniform sampler2DRect Image;

out vec4 color;

void main() {
	color = texture(Image,gl_FragCoord.xy);
}

