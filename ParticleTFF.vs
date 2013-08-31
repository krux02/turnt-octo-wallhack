#version 330 core

in vec3 a_pos1;
in vec3 a_pos2;
in float a_lifetime;

out vec3 v_pos1;
out vec3 v_pos2;
out float v_lifetime;

uniform vec3 randomDirs[64];
uniform vec3 u_origin;
uniform float u_gravity;
uniform float u_maxLifetime;

uniform sampler2D heights;
uniform vec3 lowerBound;
uniform vec3 upperBound;

float heightAt(vec2 pos) {
	float minh = lowerBound.z;
	float maxh = upperBound.z;
	return minh + texture(heights,(pos.xy+vec2(0.5)) / upperBound.xy).r * (maxh-minh);
}

vec3 normalAt(vec2 pos) {
	vec2 p1 = vec2(pos.x-0.5,pos.y);
	vec2 p2 = vec2(pos.x+0.5,pos.y);
	vec2 p3 = vec2(pos.x,pos.y-0.5);
	vec2 p4 = vec2(pos.x,pos.y+0.5);

	vec3 dir1 = vec3(p2-p1, heightAt(p2) - heightAt(p1));
	vec3 dir2 = vec3(p4-p3, heightAt(p4) - heightAt(p3));
	
	return normalize(cross(dir1,dir2));
}

void main() {
	

	if(a_lifetime > 0) {
		v_pos2 = a_pos1;
		v_pos1 = a_pos1+(a_pos1-a_pos2)+vec3(0,0,u_gravity);

		float h = heightAt(v_pos1.xy);

		if( h > v_pos1.z ) {
			vec3 dir = v_pos1 - v_pos2;
			vec3 n = normalAt(v_pos1.xy);
			dir = reflect(dir, n);
			v_pos1.z = h;
			v_pos2 = v_pos1-dir * 0.5;
		}
		
		v_lifetime = a_lifetime-1;
	}
	else {
		v_pos1 = u_origin;
		v_pos2 = u_origin+randomDirs[gl_VertexID % 64]+randomDirs[(gl_VertexID/64) % 64];
		v_lifetime = u_maxLifetime;
	}
}
