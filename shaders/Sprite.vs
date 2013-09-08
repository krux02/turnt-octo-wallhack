#version 330 core

in vec4 Vertex_os;
in vec2 TexCoord;
in vec4 Position_ws;

uniform mat4 Proj;
uniform mat4 View;

out vec2 v_texCoord;

void main() {
	
	v_texCoord = TexCoord;
	vec4 Position_cs = View * Position_ws;
	vec3 sum = Vertex_os.xyz + Position_cs.xyz;
	gl_Position = Proj * vec4(sum, 1);
}


