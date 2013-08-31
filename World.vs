#version 330 core

// Input vertex data, different for all executions of this shader.
in vec3 vertexPosition_modelspace;
in vec3 vertexNormal_modelspace;

uniform mat4 model;
uniform mat4 matrix;
uniform float seaLevel;
uniform float time;

uniform vec3 waveColor1 = vec3(0,0,1);
uniform vec3 waveColor2 = vec3(0,1,1);
uniform vec2 waveDir = vec2(0.707107);
uniform float waveAmplitudeScale = 0.35;

out vec4 v_color;
out vec3 pos_ws;
out vec3 normal_ws;

void main() {
	float wavePos = dot((model*vec4(vertexPosition_modelspace,1)).xy, waveDir)+time;
	float s = sin(wavePos);
	float c = cos(wavePos);
	float waveHeight = s * waveAmplitudeScale + seaLevel;

	vec3 waveNormal = normalize(vec3(waveDir * vec2(-c*waveAmplitudeScale),1));
	vec3 waveColor = mix(waveColor1, waveColor2, (s+1)/2);

	pos_ws = vertexPosition_modelspace;
	
	float mixValue = clamp(pos_ws.z-waveHeight, 0, 1);
	v_color = vec4(waveColor,mixValue);
	normal_ws = mix(waveNormal, vertexNormal_modelspace, mixValue);

	vec3 pos = vec3( vertexPosition_modelspace.xy, max(waveHeight, vertexPosition_modelspace.z));

    gl_Position = matrix * vec4(pos,1);
}
