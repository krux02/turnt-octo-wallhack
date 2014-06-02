#version 330 core

uniform vec3 lightDir = vec3(-0.57735);
uniform vec3 ambientColor = vec3(0.5);
uniform vec3 sunColor = vec3(1);

uniform samplerCube TextureSkybox;
uniform sampler2D TextureGround;

uniform sampler2D TextureHeightMap;
uniform vec3 LowerBound;
uniform vec3 UpperBound;

uniform vec4 CameraPos_ws;

float heightAt(vec2 pos) {
	float minh = LowerBound.z;
	float maxh = UpperBound.z;
	return minh + texture(TextureHeightMap,(pos.xy+vec2(0.5)) / UpperBound.xy).r * (maxh-minh);
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
	vec3 ref = refract(dir_ws, Normal_ws, 1/1.337);
	float x = -depth / ref.z;

	vec4 color_refract = texture(TextureGround, pos_ws.xy + x * (ref.xy));
	vec4 color_reflect = texture(TextureSkybox, reflect(dir_ws, Normal_ws));
	color = 0.45 * color_refract + 0.45 * color_reflect + 0.1 * v_color;
	color = max(color, vec4(1-0.5*depth));
}


