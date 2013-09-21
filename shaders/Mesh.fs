#version 330 core

in vec4 Normal_ws;

out vec4 color;

uniform vec4 BaseColor = vec4(1,0,0,1);

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);

void main() {
	float sunIntensity = dot(-lightDir,Normal_ws.xyz);
	vec3 light = max((sunIntensity * sunColor),ambientColor);
	color = BaseColor * vec4(light,1);
}
