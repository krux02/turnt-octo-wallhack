package main

import (
	"github.com/Jragonmiris/mathgl/examples/opengl-tutorial/helper"
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
	"github.com/go-gl/glfw"
	"github.com/krux02/mathgl"
	"io/ioutil"
	"math/rand"
	"unsafe"
	"math"
)

type ParticleVertex struct {
	Pos1     mathgl.Vec3f
	Pos2     mathgl.Vec3f
	Lifetime float32
}

type ProgramLocations struct {
	pos1        gl.AttribLocation
	pos2        gl.AttribLocation
	lifetime    gl.AttribLocation
	origin      gl.UniformLocation
	gravity     gl.UniformLocation
	maxLifetime gl.UniformLocation
}

type RenderProgramLocations struct {
	pos1        gl.AttribLocation
	lifetime    gl.AttribLocation
	matrix      gl.UniformLocation
	maxLifetime gl.UniformLocation
}

type ParticleSystem struct {
	Program         gl.Program
	Locations       ProgramLocations
	RenderProgram   gl.Program
	RenderLocations RenderProgramLocations
	Data1           gl.Buffer
	Data2           gl.Buffer
	NumParticles    int
	Origin          mathgl.Vec3f
	Gravity         float32
	InitialSpeed    float32
	MaxLifetime     float32
}

func NewParticleSystem(numParticles int, origin mathgl.Vec3f, initialSpeed, maxLifetime float32) *ParticleSystem {
	vertices := make([]ParticleVertex, numParticles)
	buffer1, buffer2 := gl.GenBuffer(), gl.GenBuffer()

	for i, _ := range vertices {
		dir := mathgl.Vec3f{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()*2 - 1}
		for dir.Len() > 1 {
			dir = mathgl.Vec3f{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()*2 - 1}
		}
		dir = dir.Mul(initialSpeed)
		vertices[i] = ParticleVertex{origin, origin.Sub(dir), rand.Float32() * maxLifetime}
	}

	buffer1.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, numParticles*int(unsafe.Sizeof(ParticleVertex{})), vertices, gl.STREAM_DRAW)

	buffer2.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, numParticles*int(unsafe.Sizeof(ParticleVertex{})), uintptr(0), gl.STREAM_DRAW)

	program := gl.CreateProgram()

	content, err := ioutil.ReadFile("ParticleTFF.vs")
	if err != nil {
		panic(err)
	}

	shader := glh.MakeShader(gl.VERTEX_SHADER, string(content))

	program.AttachShader(shader)
	program.TransformFeedbackVaryings([]string{"v_pos1", "v_pos2", "v_lifetime"}, gl.INTERLEAVED_ATTRIBS)

	program.Link()

	defer shader.Delete()

	locations := ProgramLocations{
		program.GetAttribLocation("a_pos1"),
		program.GetAttribLocation("a_pos2"),
		program.GetAttribLocation("a_lifetime"),
		program.GetUniformLocation("u_origin"),
		program.GetUniformLocation("u_gravity"),
		program.GetUniformLocation("u_maxLifetime"),
	}

	renderProgram := helper.MakeProgram("Particle.vs", "Particle.fs")
	renderLocations := RenderProgramLocations{
		renderProgram.GetAttribLocation("a_pos1"),
		renderProgram.GetAttribLocation("a_lifetime"),
		renderProgram.GetUniformLocation("matrix"),
		renderProgram.GetUniformLocation("u_maxLifetime"),
	}

	return &ParticleSystem{program, locations, renderProgram, renderLocations, buffer1, buffer2, numParticles, origin, -9.81 / 200, initialSpeed, maxLifetime}
}

func (ps *ParticleSystem) SetUniformsAndProgram() {
	ps.Program.Use()
	ps.Locations.origin.Uniform3f(ps.Origin[0], ps.Origin[1], ps.Origin[2])
	ps.Locations.gravity.Uniform1f(ps.Gravity)
	ps.Locations.maxLifetime.Uniform1f(ps.MaxLifetime)
	ps.Locations.pos1.EnableArray()
	ps.Locations.pos2.EnableArray()
	ps.Locations.lifetime.EnableArray()

	dirs := make([]float32, 64*3)
	for i := 0; i < 64; i++ {
		dir := mathgl.Vec3f{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()*2 - 1}
		for dir.Len() > 1 {
			dir = mathgl.Vec3f{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()*2 - 1}
		}
		dirs[i*3+0] = dir[0]
		dirs[i*3+1] = dir[1]
		dirs[i*3+2] = dir[2]
	}

	ps.Program.GetUniformLocation("randomDirs").Uniform3fv(64, dirs)
}

func (ps *ParticleSystem) DoStep() {
	ps.SetUniformsAndProgram()

	gl.Enable(gl.RASTERIZER_DISCARD)
	defer gl.Disable(gl.RASTERIZER_DISCARD)

	ps.Data1.Bind(gl.ARRAY_BUFFER)

	ps.Locations.pos1.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Pos1))
	ps.Locations.pos2.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Pos2))
	ps.Locations.lifetime.AttribPointer(1, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Lifetime))

	time := glfw.Time()
	ps.Locations.origin.Uniform3f(100*float32(math.Sin(time)),100*float32(math.Cos(time)),100)

	ps.Data2.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	gl.BeginTransformFeedback(gl.POINTS)
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
	gl.EndTransformFeedback()

	tmp := ps.Data1
	ps.Data1 = ps.Data2
	ps.Data2 = tmp

}

func (ps *ParticleSystem) Render(matrix *mathgl.Mat4f) {
	ps.RenderProgram.Use()
	ps.RenderLocations.pos1.EnableArray()
	ps.RenderLocations.lifetime.EnableArray()
	ps.RenderLocations.matrix.UniformMatrix4f(false, (*[16]float32)(matrix))
	ps.RenderLocations.maxLifetime.Uniform1f(ps.MaxLifetime)

	ps.Data1.Bind(gl.ARRAY_BUFFER)
	ps.RenderLocations.pos1.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Pos1))
	ps.RenderLocations.lifetime.AttribPointer(1, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Lifetime))
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
}
