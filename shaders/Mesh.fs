#version 330 core

in vec4 Normal_ws;

out vec4 color;


uniform vec4 BaseColor
uniform vec4 lightDir

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);

void main() {
	float sunIntensity = dot(-lightDir,normal_ws);
	vec3 light = max((sunIntensity * sunColor),ambientColor);
	color = BaseColor * vec4(light,1);
}
