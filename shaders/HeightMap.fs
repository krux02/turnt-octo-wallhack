#version 330 core

uniform vec3 LowerBound;
uniform vec3 UpperBound;

uniform sampler1D ColorBand;
uniform sampler2D Texture;
uniform sampler2D Slope;

uniform vec3 LightDir = vec3(-0.57735);
uniform vec3 AmbientColor = vec3(0.5);
uniform vec3 SunColor = vec3(1);

in vec4 v_color;
in vec4 pos_ws;
in vec4 normal_ws;
// Ouput data
out vec4 color;

void main() {
	float sunIntensity = dot(-LightDir,normal_ws.xyz);
	vec3 light = max((sunIntensity * SunColor),AmbientColor);
	
	vec3 colorA = texture(ColorBand,(pos_ws.z-LowerBound.z)/(UpperBound.z-LowerBound.z)).rgb * texture(Texture, pos_ws.xy).xyz;;
	vec3 colorB = texture(Slope, pos_ws.xz).xyz;
	vec3 colorC = texture(Slope, pos_ws.yz).xyz;
	float fractionA = pow(max(normal_ws.z, 0), 15);
	float fractionB = pow(abs(normal_ws.y), 15);
	float fractionC = pow(abs(normal_ws.x), 15);
	float len = fractionA+fractionB+fractionC;
	color.rgb = colorA*vec3(fractionA/len)+colorB*vec3(fractionB/len)+colorC*vec3(fractionC/len);
	//color.rgb =  light * mix(v_color.xyz, color.rgb, v_color.w);
	color.a = 1;
}


