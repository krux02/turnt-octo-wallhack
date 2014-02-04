#version 330 core

uniform sampler2DRect U_screenRect;

out vec4 color;

void main() {
	color = texture(U_screenRect, gl_FragCoord.xy);
}
