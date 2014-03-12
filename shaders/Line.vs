#version 330 core

uniform Mat4 View;
uniform Mat4 Proj;

in Vec4 Vertex_ws;
in Vec4 Color;

out vColor;

void main() {
	out.vColor = in.Color;
}
