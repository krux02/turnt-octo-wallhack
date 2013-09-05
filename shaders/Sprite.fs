uniform sampler2D image;

in vec2 v_texCoord;

out vec4 color;

void main() {
	color = texture(image)
}

