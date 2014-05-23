#version 330 core

in vec4 Vertex_ms;
in vec4 Normal_ms;

out vec4 Normal_cs;
out vec4 Position_cs;

uniform mat4 Proj;
uniform mat4 View;
uniform mat4 Model;

uniform vec4 ClippingPlane_ws;

void main() {
	
	Normal_cs = View * Model * Normal_ms;
	vec4 Position_ws = Model * Vertex_ms;
	Position_cs = View * Position_ws;
	gl_ClipDistance[0] = dot(Position_ws, ClippingPlane_ws);
	gl_Position = Proj * Position_cs;
}
