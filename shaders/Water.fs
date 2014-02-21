#version 330 core


uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);
uniform vec4 U_clippingPlane;

uniform sampler2D GroundTexture;

uniform sampler2D HeightMap;
uniform vec3 LowerBound;
uniform vec3 UpperBound;

uniform vec4 CameraPos_ws;

float heightAt(vec2 pos) {
	float minh = LowerBound.z;
	float maxh = UpperBound.z;
	return minh + texture(HeightMap,(pos.xy+vec2(0.5)) / UpperBound.xy).r * (maxh-minh);
}

in vec4 v_color;
in vec4 pos_ws;
in vec4 pos_cs;
in vec3 normal_ws;

// Ouput data
out vec4 color;

void main()
{
	if( dot(pos_ws, U_clippingPlane) < 0 ) {
		discard;
	}

	float height = heightAt(pos_ws.xy);
	float depth = pos_ws.z - height;
	if( depth < 0 ) {
		discard;
	}
	
	float sunIntensity = dot(-lightDir, normal_ws);
	vec3 light = max((sunIntensity * sunColor), ambientColor);
	vec3 dir_ws = (pos_ws.xyz - CameraPos_ws.xyz);

	dir_ws = refract(dir_ws, normal_ws, 1/1.337);
	float x = -depth / dir_ws.z;

	color = texture(GroundTexture, pos_ws.xy + x * (dir_ws.xy));
	color = mix( color, v_color, depth/4 );
}


