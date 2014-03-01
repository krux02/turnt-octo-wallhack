#version 330

uniform samplerCube Skybox;

in vec3 TexCoord;

out vec4 color;

void main() {
	color = texture(Skybox,TexCoord);
	//color.xyz = TexCoord;
}