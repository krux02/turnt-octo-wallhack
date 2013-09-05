in vec4 vertex_os;
in vec2 texCoord;
in vec4 position_cs;

out vec2 v_texCoord;

uniform mat4 Proj;

void main() {
	v_texCoord = texCoord;
	gl_Position = Proj * vec4( vertex_os.xyz + position_cs.xyz, 1);
}
