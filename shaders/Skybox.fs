#version 330

uniform samplerCube TextureSkybox;

in vec3 TexCoord;

out vec4 color;

void main() {
	color = texture(TextureSkybox,TexCoord);
	//color.xyz = TexCoord;
}