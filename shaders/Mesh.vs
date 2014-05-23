#version 330 core

uniform mat4 Proj;
uniform mat4 View;
uniform mat4 Model;
uniform vec4 ClippingPlane_ws;

in vec4 Vertex_ms;
in vec4 Normal_ms;

out vec4 Normal_ws;

void main() {
	vec4 Position_ws = Model * Vertex_ms;
	gl_Position = Proj * View * Position_ws;
	gl_ClipDistance[0] = dot(Position_ws, ClippingPlane_ws);
	Normal_ws = Model * Normal_ms;
}
