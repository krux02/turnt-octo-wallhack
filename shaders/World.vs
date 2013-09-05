#version 330 core

// Input vertex data, different for all executions of this shader.
in vec3 Vertex_ms;
in vec3 Normal_ms;

uniform mat4 Model;
uniform mat4 Matrix;
uniform float SeaLevel;
uniform float Time;

uniform vec3 waveColor1 = vec3(0,0,1);
uniform vec3 waveColor2 = vec3(0,1,1);
uniform vec2 waveDir = vec2(0.707107);
uniform float waveAmplitudeScale = 0.35;

out vec4 v_color;
out vec3 pos_ws;
out vec3 normal_ws;

void main() {
	float wavePos = dot((Model*vec4(Vertex_ms,1)).xy, waveDir)+Time;
	float s = sin(wavePos);
	float c = cos(wavePos);
	float waveHeight = s * waveAmplitudeScale + SeaLevel;

	vec3 waveNormal = normalize(vec3(waveDir * vec2(-c*waveAmplitudeScale),1));
	vec3 waveColor = mix(waveColor1, waveColor2, (s+1)/2);

	pos_ws = Vertex_ms;
	
	float mixValue = clamp(pos_ws.z-waveHeight, 0, 1);
	v_color = vec4(waveColor,mixValue);
	normal_ws = mix(waveNormal, Normal_ms, mixValue);

	vec3 pos = vec3( Vertex_ms.xy, max(waveHeight, Vertex_ms.z));

    gl_Position = Matrix * vec4(pos,1);
}
