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

uniform float normal_scale = 0.5; // don't forget: this is the default value! 
uniform mat4 Model;
uniform mat4 View;
uniform mat4 Proj;

in vec4 Normal;

void main()
{
    vec4 v0     = gl_in[0].gl_Position;
    gl_Position = gl_ModelViewProjectionMatrix * v0;
    EmitVertex();

    vec4 v1     = v0 + vec4(vs_normal[0] * normal_scale, 0);
    gl_Position = gl_ModelViewProjectionMatrix * v1;
    EmitVertex();

    EndPrimitive();
}
