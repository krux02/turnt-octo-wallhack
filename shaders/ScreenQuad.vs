#version 330 core

in vec4 Vertex_ndc;

void main() {
	gl_Position = Vertex_ndc;
}
