#version 330 core

uniform float Highlight;

uniform float Min_h = -17;
uniform float Max_h = 28;
uniform sampler1D U_color;
uniform sampler2D U_texture;
uniform sampler2D U_slope;

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);
uniform vec4 U_clippingPlane;

in vec4 v_color;
in vec4 pos_ws;
in vec3 normal_ws;
// Ouput data
out vec3 color;

void main()
{
	if( dot(pos_ws, U_clippingPlane) < 0 ) {
		discard;
	}
	
	
	
	float sunIntensity = dot(-lightDir,normal_ws);
	vec3 light = max((sunIntensity * sunColor),ambientColor);
	
	vec3 colorA = texture(U_color,(pos_ws.z-Min_h)/(Max_h-Min_h)).rgb * texture(U_texture, pos_ws.xy).xyz;;
	vec3 colorB = texture(U_slope, pos_ws.xz).xyz;
	vec3 colorC = texture(U_slope, pos_ws.yz).xyz;
	float fractionA = pow(max(dot(normal_ws,vec3(0,0,1)),0), 15);
	float fractionB = pow(abs(dot(normal_ws,vec3(0,1,0))),15);
	float fractionC = pow(abs(dot(normal_ws,vec3(1,0,0))),15);
	float len = fractionA+fractionB+fractionC;
	color = colorA*vec3(fractionA/len)+colorB*vec3(fractionB/len)+colorC*vec3(fractionC/len);
	
	color =  light * mix(v_color.xyz, color, v_color.w);
	
	

	//if( gl_PrimitiveID == Highlight ) {
	//	color = vec3(1) - color;
	//}
}
