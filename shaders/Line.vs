#version 330 core

uniform mat4 View;
uniform mat4 Proj;

in vec4 Vertex_ws;
in vec4 Color;

out vec4 vColor;

void main() {
	vColor = Color;
}
