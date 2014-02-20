#version 330 core

uniform sampler2D PalmTree;
uniform vec4 U_clippingPlane;

in vec2 v_texCoord;
in vec4 pos_ws;

out vec4 color;

void main() {
	if( dot(pos_ws, U_clippingPlane) < 0 ) {
		discard;
	}

	color = texture(PalmTree, v_texCoord);
	if(color.a < 0.5)
		discard;
}		
