#version 330 core

// assuming you have vertex normals, you need to render a vertex
// only a single time. with any other prim type, you may render
// the same normal multiple times

// Geometry shaders can only output points, line strips or triangle
// strips by definition. you output a single line per vertex. therefore, 
// the maximum number of vertices per line_strip is 2. This is effectively
// the same as rendering distinct line segments.
layout (points) in;
layout (line_strip, max_vertices = 2) out;

uniform float normal_scale = 0.5;
uniform vec4 BaseColor = vec4(1,0,0,1);

uniform mat4 Model;
uniform mat4 View;
uniform mat4 Proj;

in vec4 pos_ws[];
in vec3 Normal_ws[];

out vec4 vColor;

void main()
{
	mat4 m = Proj * View;

    gl_Position = m * pos_ws[0];
    vColor = BaseColor;
    EmitVertex();

    gl_Position = m* (pos_ws[0] + vec4(Normal_ws[0] * normal_scale,0));
    vColor = BaseColor;
    EmitVertex();

    EndPrimitive();
}
