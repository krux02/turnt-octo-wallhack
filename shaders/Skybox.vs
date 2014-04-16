#version 330

/*
uniform vec3[8] vertices = vec3[]( 
	vec3( -1,-1,-1 ),
	vec3(  1,-1,-1 ),
	vec3( -1, 1,-1 ),
	vec3(  1, 1,-1 ),
	vec3( -1,-1, 1 ),
	vec3(  1,-1, 1 ),
	vec3( -1, 1, 1 ),
	vec3(  1, 1, 1 )
);
*/

uniform mat4 View;
uniform mat4 Proj;

in vec3 InTexCoord;

out vec3 TexCoord;

void main() {
	TexCoord = InTexCoord; //vertices[gl_VertexID];
	gl_Position = Proj * View * vec4(TexCoord, 0);
}