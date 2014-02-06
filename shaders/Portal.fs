#version 330 core

uniform sampler2DRect U_Image;

out vec4 color;

void main() {
	color = texture(U_Image,gl_FragCoord.xy);
}

