#version 330 core

in vec4 Vertex_os;
in vec2 TexCoord;
in vec4 Position_ws;

uniform mat4 Proj;
uniform mat4 View;
uniform mat3 Rot2D;

out vec2 v_texCoord;
out vec4 pos_ws;

void main() {
	v_texCoord = TexCoord;
	pos_ws = vec4(Position_ws.xyz + Rot2D * Vertex_os.xyz,1);
	vec4 Position_cs = View * vec4(Position_ws.xyz + Rot2D * Vertex_os.xyz,1);
	vec3 sum = Position_cs.xyz;
	gl_Position = Proj * vec4(sum, 1);
}
