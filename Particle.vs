#version 330 core

in vec3 a_pos1;
in vec3 a_pos2;
in float a_lifetime;

out vec4 vertexColor;

uniform mat4 matrix;
uniform float u_maxLifetime = 500;

void main() {
	gl_Position = matrix * vec4(a_pos1,1);

	vertexColor.r = clamp(2*a_lifetime/u_maxLifetime,0,1);
	vertexColor.g = clamp(2*a_lifetime/u_maxLifetime,1,2)-1;
	vertexColor.b = 0;
	vertexColor.a = 0.5;
}
