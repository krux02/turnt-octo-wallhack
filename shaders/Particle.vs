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

void main() {
	vec4 Position_ws = vec4(Pos1,1);
	gl_Position = Proj * View * Position_ws;
	v_texCoord = TexCoord;

	vec4 Position_cs = View * Position_ws;
 	vec3 sum = Vertex_os.xyz + Position_cs.xyz;
 	gl_Position = Proj * vec4(sum, 1);

	vertexColor.r = clamp(2*Lifetime/MaxLifetime,0,1);
	vertexColor.g = clamp(2*Lifetime/MaxLifetime,1,2)-1;
	vertexColor.b = 0;
	vertexColor.a = 0.5;
}



