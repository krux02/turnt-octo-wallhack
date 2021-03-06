#version 330 core

// Input vertex data, different for all executions of this shader.

uniform mat4 Model;
uniform mat4 View;
uniform mat4 Proj;

uniform float Time;
uniform float WaterHeight;

uniform vec3 waveColor1 = vec3(0,0,1);
uniform vec3 waveColor2 = vec3(0,1,1);
uniform vec2 waveDir = vec2(0.707107);
uniform float waveAmplitudeScale = 0.35;
uniform vec4 ClippingPlane_ws;

in vec3 Vertex_ms;
in vec3 Normal_ms;

out vec4 v_color;
out vec4 pos_ws;
out vec4 pos_cs;
out vec3 Normal_ws;
out vec3 Normal_cs;

void main() {
	float wavePos = dot(Vertex_ms.xy, waveDir)+Time;
	float s = sin(wavePos);
	float c = cos(wavePos);
	float waveHeight = s * waveAmplitudeScale;

	vec3 waveNormal = normalize(vec3(waveDir * vec2(-c*waveAmplitudeScale),1));
	vec3 waveColor = mix(waveColor1, waveColor2, (s+1)/2);
	
	v_color = vec4(waveColor,1);
	Normal_ws = waveNormal;
	Normal_cs = (View * vec4(waveNormal, 0)).xyz;


	vec3 pos = vec3( Vertex_ms.xy, WaterHeight + waveHeight + 10 );

	pos_ws = Model * vec4(pos,1);
	pos_cs = View * pos_ws;

	gl_Position =  Proj * pos_cs;

	gl_ClipDistance[0] = dot(pos_ws, ClippingPlane_ws);
}

