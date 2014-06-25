#version 330 core

uniform sampler2D Image;
uniform vec2 ViewPortSize = vec2(1280,800);

out vec4 color;

void main() {
	color = texture(Image, gl_FragCoord.xy / ViewPortSize );
}

