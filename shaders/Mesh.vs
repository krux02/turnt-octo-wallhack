#version 330 core

in vec4 Vertex_ms;
in vec4 Normal_ms;

out vec4 Normal_ws;

uniform mat4 Model;
uniform mat4 View;
uniform mat4 Proj;

void main() {
	gl_Position = Proj * View * Model * Vertex_ms;
	Normal_ws = Model * Normal_ms;
}
