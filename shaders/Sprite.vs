#version 330 core

uniform vec4 ClippingPlane_ws;

in vec4 Vertex_ms;
in vec2 TexCoord;
in vec4 InstancePosition_ws;

uniform mat4 Proj;
uniform mat4 View;
uniform mat3 Rot2D;

out vec2 v_texCoord;
out vec4 pos_ws;

void main() {
	v_texCoord = TexCoord;
	pos_ws = vec4(InstancePosition_ws.xyz + Rot2D * Vertex_ms.xyz,1);
	vec4 Position_cs = View * vec4(InstancePosition_ws.xyz + Rot2D * Vertex_ms.xyz,1);
	vec3 sum = Position_cs.xyz;
	gl_Position = Proj * vec4(sum, 1);
	gl_ClipDistance[0] = dot(ClippingPlane_ws, pos_ws);
}
