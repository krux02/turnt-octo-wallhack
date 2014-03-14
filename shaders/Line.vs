#version 330 core

uniform mat4 Proj;
uniform mat4 View;

in vec4 Vertex_ws;
in vec4 Color;

out vec4 vColor;

void main() {
	gl_Position = Proj*View*Vertex_ws;
	vColor = Color;
}
