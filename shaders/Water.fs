#version 330 core

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);

uniform samplerCube Skybox;
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
in vec3 Normal_ws;
in vec3 Normal_cs;

// Ouput data
out vec4 color;

void main()
{

	float height = heightAt(pos_ws.xy);
	float depth = pos_ws.z - height;
	if( depth < 0 ) {
		discard;
	}
	
	float sunIntensity = dot(-lightDir, Normal_ws);
	vec3 light = max((sunIntensity * sunColor), ambientColor);
	vec3 dir_ws = (pos_ws.xyz - CameraPos_ws.xyz);
	dir_ws = refract(dir_ws, Normal_ws, 1/1.337);
	float x = -depth / dir_ws.z;

	vec4 color_refract = texture(GroundTexture, pos_ws.xy + x * (dir_ws.xy));
	vec4 color_reflect = texture(Skybox, reflect(pos_cs.xyz, Normal_cs));
	color = 0.33 * color_refract + 0.33 * color_reflect + 0.33 * v_color;
	color = max(color, vec4(1-0.5*depth));
}


