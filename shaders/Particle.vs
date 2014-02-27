#version 330 core

// instance data (divisor = 1)
in vec3 Pos1;
in vec3 Pos2;
in float Lifetime;
// vertex data (divisor = 0)
in vec4 Vertex_os;
in vec2 TexCoord;

out vec4 vertexColor;
out vec2 v_texCoord;

uniform mat4 Proj;
uniform mat4 View;
uniform float MaxLifetime = 500;
uniform vec4 ClippingPlane_ws;

void main() {
	vec4 pos_ws = vec4(Pos1,1);
	v_texCoord = TexCoord;
	vec4 pos_cs = View * pos_ws;
 	gl_Position = Proj * vec4(Vertex_os.xyz + pos_cs.xyz, 1);
 	gl_ClipDistance[0] = dot(pos_ws, ClippingPlane_ws);

	vertexColor.r = clamp(2*Lifetime/MaxLifetime,0,1);
	vertexColor.g = clamp(2*Lifetime/MaxLifetime,1,2)-1;
	vertexColor.b = 0;
	vertexColor.a = 1;
}



