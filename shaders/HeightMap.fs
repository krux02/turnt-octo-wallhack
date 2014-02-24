#version 330 core

uniform float Min_h = -17;
uniform float Max_h = 28;

uniform sampler1D U_color;
uniform sampler2D U_texture;
uniform sampler2D U_slope;

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);

in vec4 v_color;
in vec4 pos_ws;
in vec4 normal_ws;
// Ouput data
out vec4 color;

void main() {
	float sunIntensity = dot(-lightDir,normal_ws.xyz);
	vec3 light = max((sunIntensity * sunColor),ambientColor);
	
	vec3 colorA = texture(U_color,(pos_ws.z-Min_h)/(Max_h-Min_h)).rgb * texture(U_texture, pos_ws.xy).xyz;;
	vec3 colorB = texture(U_slope, pos_ws.xz).xyz;
	vec3 colorC = texture(U_slope, pos_ws.yz).xyz;
	float fractionA = pow(max(normal_ws.z, 0), 15);
	float fractionB = pow(abs(normal_ws.y), 15);
	float fractionC = pow(abs(normal_ws.x), 15);
	float len = fractionA+fractionB+fractionC;
	color.rgb = colorA*vec3(fractionA/len)+colorB*vec3(fractionB/len)+colorC*vec3(fractionC/len);
	color.rgb =  light * mix(v_color.xyz, color.rgb, v_color.w);
	color.a = 1;
}


