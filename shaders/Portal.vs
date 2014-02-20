#version 330 core

in vec4 Vertex_ms;
in vec4 Normal_ms;

out vec4 Normal_cs;
out vec4 pos_ws;
out vec4 pos_cs;

uniform mat4 Proj;
uniform mat4 View;
uniform mat4 Model;

void main() {

	gl_Position = Proj * View * Model * Vertex_ms;
	Normal_cs = View * Model * Normal_ms;
	pos_ws = Model * Vertex_ms;
	pos_cs = View * Model * Vertex_ms;
}
