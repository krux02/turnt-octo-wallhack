in vec4 position_os;
out vec4 position_cs;

uniform mat4 ViewModel;

void main() {
	position_cs = ViewProj * position_os;
}

