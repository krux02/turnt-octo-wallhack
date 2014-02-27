#version 330 core

in vec4 Vertex_ms;
in vec4 Normal_ms;

out vec4 Normal_cs;
out vec4 Position_cs;

uniform mat4 Proj;
uniform mat4 View;
uniform mat4 Model;

uniform vec4 ClippingPlane_cs;

void main() {
	gl_Position = Proj * View * Model * Vertex_ms;
	Normal_cs = View * Model * Normal_ms;
	Position_cs = View * Model * Vertex_ms;
	gl_ClipDistance[0] = dot(Position_cs, ClippingPlane_cs);
}
