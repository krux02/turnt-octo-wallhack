#version 330 core

in vec4 Position_ndc;

void main() {
	gl_Position = Position_ndc;
}
