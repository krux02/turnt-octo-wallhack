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

vec2 textureCoords(vec2 worldPos) {
	return (worldPos.xy+vec2(0.5)) / UpperBound.xy;
}

// translates normalized height to real world height
float worldHeight(float nHeight) {
	return LowerBound.z + nHeight * (UpperBound.z - LowerBound.z);
}

float heightAt(vec2 pos) {
	return worldHeight(texture(HeightMap, textureCoords(pos)).r);
}

vec3 normalAt(vec2 pos) {
	float x = pos.x;
	float y = pos.y;
	float l = pos.x - 1;
	float r = pos.x + 1;
	float b = pos.y - 1;
	float t = pos.y + 1;

	float hi = heightAt(pos);
	float lh = heightAt(vec2(l, y)) - hi;
	float rh = heightAt(vec2(r, y)) - hi;
	float bh = heightAt(vec2(x, b)) - hi;
	float th = heightAt(vec2(x, t)) - hi;

	vec3 v1 = normalize(vec3(1, 0, rh));
	vec3 v2 = normalize(vec3(0, 1, th));
	vec3 v3 = normalize(vec3(-1, 0, lh));
	vec3 v4 = normalize(vec3(0, -1, bh));

	vec3 n1 = normalize(cross(v1, v2));
	vec3 n2 = normalize(cross(v2, v3));
	vec3 n3 = normalize(cross(v3, v4));
	vec3 n4 = normalize(cross(v4, v1));

	return normalize(n1 + n2 + n3 + n4);
}

void main() {
	pos_ws = Model*vec4( Vertex_ms.xy, heightAt(Vertex_ms.xy), 1);
	normal_ws = Model*vec4( normalAt(Vertex_ms.xy), 0);

	gl_Position = Matrix * vec4( Vertex_ms.xy, heightAt(Vertex_ms.xy), 1);
	gl_ClipDistance[0] = dot(pos_ws, ClippingPlane_ws);
}
