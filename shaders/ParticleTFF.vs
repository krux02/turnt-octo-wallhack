#version 330 core

in vec3 Pos1;
in vec3 Pos2;
in vec3 StartDir;
in float Lifetime;

out vec3 v_Pos1;
out vec3 v_Pos2;
out float v_Lifetime;

uniform float FrictionFactor = 0.99 ;
uniform float Timestep = 0.016666666;
uniform vec3 RandomDirs[64];
uniform vec3 Origin;
uniform float Gravity;
uniform float MaxLifetime;
uniform vec3 Dir = vec3(0);

uniform sampler2D HeightMap;
uniform vec3 LowerBound;
uniform vec3 UpperBound;

float heightAt(vec2 pos) {
	float minh = LowerBound.z;
	float maxh = UpperBound.z;
	return minh + texture(HeightMap,(pos.xy+vec2(0.5)) / UpperBound.xy).r * (maxh-minh);
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
	if(Lifetime > 0) {
		v_Pos2 = Pos1;
		v_Pos1 = Pos1+(Pos1-Pos2)*FrictionFactor+vec3(0,0,Gravity)*Timestep;

		float h = heightAt(v_Pos1.xy);

		if( h > v_Pos1.z ) {
			vec3 dir = v_Pos1 - v_Pos2;
			vec3 n = normalAt(v_Pos1.xy);
			dir = reflect(dir, n);
			v_Pos1.z = h;
			v_Pos2 = v_Pos1-dir * 0.5;
		}
		
		v_Lifetime = Lifetime-Timestep;
	}
	else {
		v_Pos1 = Origin;
		v_Pos2 = Origin + StartDir*Timestep + Dir*Timestep;
		v_Lifetime = MaxLifetime;
	}
}
