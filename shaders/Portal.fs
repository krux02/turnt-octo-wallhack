#version 330 core

uniform sampler2DRect U_Image;

out vec4 color;

void main() {
	color = 0.7 * texture(U_Image,gl_FragCoord.xy);
}

