#version 330 core

// Input vertex data, different for all executions of this shader.
in vec3 Vertex_ms;
in vec3 Normal_ms;

uniform mat4 Model;
uniform mat4 Matrix;
uniform float Time;

uniform vec3 waveColor1 = vec3(0,0,1);
uniform vec3 waveColor2 = vec3(0,1,1);
uniform vec2 waveDir = vec2(0.707107);
uniform float waveAmplitudeScale = 0.35;

out vec4 v_color;
out vec4 pos_ws;
out vec3 normal_ws;

void main() {
	float wavePos = dot(Vertex_ms.xy, waveDir)+Time;
	float s = sin(wavePos);
	float c = cos(wavePos);
	float waveHeight = s * waveAmplitudeScale;

	vec3 waveNormal = normalize(vec3(waveDir * vec2(-c*waveAmplitudeScale),1));
	vec3 waveColor = mix(waveColor1, waveColor2, (s+1)/2);
	
	v_color = vec4(waveColor,1);
	normal_ws = waveNormal;


	vec3 pos = vec3( Vertex_ms.xy, Vertex_ms.z + waveHeight + 10 );

	pos_ws = Model * vec4(pos,1);
    gl_Position = Matrix * vec4(pos,1);
}
