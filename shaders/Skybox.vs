#version 330

uniform vec3 vertices[ 8 ] = { 
	vec3( -1,-1,-1 ),
	vec3(  1,-1,-1 ),
	vec3( -1, 1,-1 ),
	vec3(  1, 1,-1 ),
	vec3( -1,-1, 1 ),
	vec3(  1,-1, 1 ),
	vec3( -1, 1, 1 ),
	vec3(  1, 1, 1 )
};

uniform mat4 View;
uniform mat4 Proj;

out vec3 TexCoord;

void main() {
	TexCoord = vertices[gl_VertexID];
	gl_Position = Proj * View * vec4(TexCoord,0);
}