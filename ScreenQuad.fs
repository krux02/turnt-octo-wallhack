#version 330 core

uniform sampler2DRect u_screenRect;

out vec4 color;

void main() {
	vec4 texValue = texture(u_screenRect, gl_FragCoord.xy);
	color = texValue;
}
