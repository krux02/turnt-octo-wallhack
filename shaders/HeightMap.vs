#version 330 core

uniform vec3 LowerBound;
uniform vec3 UpperBound;

uniform mat4 Model;
uniform mat4 Matrix;
uniform sampler2D HeightMap;
uniform vec4 ClippingPlane_ws;

// Input vertex data, different for all executions of this shader.
in vec3 Vertex_ms;
in vec3 Normal_ms;

out vec4 pos_ws;
out vec4 normal_ws;

float heightAt(vec2 pos) {
	float minh = LowerBound.z;
	float maxh = UpperBound.z;
	return minh + texture(HeightMap,(pos.xy+vec2(0.5)) / UpperBound.xy).r * (maxh-minh);
}

void main() {
	pos_ws = Model*vec4(Vertex_ms,1);
	normal_ws = Model*vec4(Normal_ms,0);
	gl_Position = Matrix * vec4( Vertex_ms.xy, heightAt(Vertex_ms.xy), 1);
	gl_ClipDistance[0] = dot(pos_ws, ClippingPlane_ws);
}
