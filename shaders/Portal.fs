#version 330 core
#define M_PI 3.1415926535897932384626433832795

uniform sampler2DRect U_Image;
uniform vec4 U_clippingPlane;


in vec4 Normal_cs;
in vec4 pos_ws;
in vec4 pos_cs;

out vec4 color;

vec4 mymix(vec4 color, float alpha) {
	float a = alpha * 6 / M_PI;

	float x = 1 - min(1, min(a, 3-a));
	float y = 1 - min(1, abs(a - 1));
	float z = 1 - min(1, abs(a - 2));

	float r = dot(vec4(x,y,z,0), color);
	float g = dot(vec4(y,z,x,0), color);
	float b = dot(vec4(z,x,y,0), color);

	return vec4(r,g,b, color.a);
	
	/*
	if (a < 1) {
		return color.rgba;
	} else if (a < 2) {
		return color.gbra;
	} else {
		return color.brga;
	}
	*/
}

void main() {
	if( dot(pos_ws, U_clippingPlane) < 0 ) {
		discard;
	}

	vec4 t = texture(U_Image,gl_FragCoord.xy);
	float alpha = acos(abs(dot(normalize(pos_cs.xyz), Normal_cs.xyz)));

	float dist = length(pos_cs.xyz);
	dist = min(dist / 3, 1);
	color = mix(t, mymix(t, alpha), dist);
}

