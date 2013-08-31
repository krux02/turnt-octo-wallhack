#version 330 core

uniform float time;
uniform float highlight;

uniform float min_h = -17;
uniform float max_h = 28;
uniform sampler1D u_color;
uniform sampler2D u_texture;
uniform sampler2D u_slope;
uniform sampler2DRect u_screenRect;

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);

in vec4 v_color;
in vec3 pos_ws;
in vec3 normal_ws;
// Ouput data
out vec3 color;

void main()
{
	float light = dot(-lightDir,normal_ws);
	vec3 texColorA;
	if(dot(normal_ws,vec3(0,0,1)) > 0.5)
		texColorA = texture(u_color,(pos_ws.z-min_h)/(max_h-min_h)).rgb;
	else
		texColorA = texture(u_slope, pos_ws.xy).xyz;

	vec3 texColorB = texture(u_texture, pos_ws.xy).xyz;
	vec3 texColor = texColorA*texColorB;
	texColor = max((light * sunColor),ambientColor) * mix(v_color.xyz, texColor, v_color.w);
	vec4 foreGround = texture(u_screenRect, gl_FragCoord.xy);
	//texColor = mix(foreGround.xyz, texColor, 0.9);
	//texColor = foreGround.xyz;

	/*
	if (int(gl_FragCoord.x) % 2 == 0) {
		color = foreGround.rgb * vec3(foreGround.a) + texColor * vec3(1-foreGround.a);
	}
	else {
		color = texColor;
	}
	*/







	color = texColor;
	//color = texture(screenRect, gl_FragCoord.xy).rgb;

	if( gl_PrimitiveID == highlight ) {
		color = vec3(1) - color;
	}
}
