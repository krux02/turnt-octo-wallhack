#version 330 core

in vec3 Pos1;
in vec3 Pos2;
in float Lifetime;

out vec4 vertexColor;

uniform mat4 Matrix;
uniform float MaxLifetime = 500;

void main() {
	gl_Position = Matrix * vec4(Pos1,1);

	vertexColor.r = clamp(2*Lifetime/MaxLifetime,0,1);
	vertexColor.g = clamp(2*Lifetime/MaxLifetime,1,2)-1;
	vertexColor.b = 0;
	vertexColor.a = 0.5;
}
