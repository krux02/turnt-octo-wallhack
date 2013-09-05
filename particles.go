package main

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
	"github.com/krux02/mathgl"
	"io/ioutil"
	"math"
	"math/rand"
	"unsafe"
)

import "fmt"

type ParticleVertex struct {
	Pos1     mathgl.Vec3f
	Pos2     mathgl.Vec3f
	Lifetime float32
}

type ProgramLocations struct {
	Pos1, Pos2, Lifetime            gl.AttribLocation
	Origin, Gravity, MaxLifetime    gl.UniformLocation
	Heights, LowerBound, UpperBound gl.UniformLocation
}

type RenderProgramLocations struct {
	Pos1        gl.AttribLocation
	Lifetime    gl.AttribLocation
	Matrix      gl.UniformLocation
	MaxLifetime gl.UniformLocation
}

type ParticleSystem struct {
	Program         gl.Program
	Locations       ProgramLocations
	RenderProgram   gl.Program
	RenderLocations RenderProgramLocations
	VertexArray     gl.VertexArray
	Data1           gl.Buffer
	Data2           gl.Buffer
	NumParticles    int
	Origin          mathgl.Vec3f
	Gravity         float32
	InitialSpeed    float32
	MaxLifetime     float32
}

func NewParticleSystem(numParticles int, Origin mathgl.Vec3f, initialSpeed, MaxLifetime float32) *ParticleSystem {
	vertices := make([]ParticleVertex, numParticles)
	buffer1, buffer2 := gl.GenBuffer(), gl.GenBuffer()

	for i, _ := range vertices {
		dir := mathgl.Vec3f{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()*2 - 1}
		for dir.Len() > 1 {
			dir = mathgl.Vec3f{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()*2 - 1}
		}
		dir = dir.Mul(initialSpeed)
		vertices[i] = ParticleVertex{Origin, Origin.Sub(dir), rand.Float32() * MaxLifetime}
	}

	vao := gl.GenVertexArray()
	vao.Bind()

	buffer1.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, numParticles*int(unsafe.Sizeof(ParticleVertex{})), vertices, gl.STREAM_DRAW)

	buffer2.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, numParticles*int(unsafe.Sizeof(ParticleVertex{})), uintptr(0), gl.STREAM_DRAW)

	program := gl.CreateProgram()

	content, err := ioutil.ReadFile("shaders/ParticleTFF.vs")
	if err != nil {
		panic(err)
	}

	shader := glh.MakeShader(gl.VERTEX_SHADER, string(content))

	program.AttachShader(shader)
	program.TransformFeedbackVaryings([]string{"v_Pos1", "v_Pos2", "v_Lifetime"}, gl.INTERLEAVED_ATTRIBS)
	program.Link()

	fmt.Println(program.GetInfoLog())

	defer shader.Delete()

	ProgLoc := ProgramLocations{}
	BindLocations(program, &ProgLoc)

	renderProgram := MakeProgram("Particle.vs", "Particle.fs")
	RenderLoc := RenderProgramLocations{}
	BindLocations(renderProgram, &RenderLoc)

	return &ParticleSystem{program, ProgLoc, renderProgram, RenderLoc, vao, buffer1, buffer2, numParticles, Origin, -9.81 / 200, initialSpeed, MaxLifetime}
}

func (ps *ParticleSystem) SetUniformsAndProgram() {
	ps.Program.Use()
	ps.Locations.Origin.Uniform3f(ps.Origin[0], ps.Origin[1], ps.Origin[2])
	ps.Locations.Gravity.Uniform1f(ps.Gravity)
	ps.Locations.MaxLifetime.Uniform1f(ps.MaxLifetime)
	ps.Locations.Pos1.EnableArray()
	ps.Locations.Pos2.EnableArray()
	ps.Locations.Lifetime.EnableArray()

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

	ps.Program.GetUniformLocation("RandomDirs").Uniform3fv(64, dirs)
}

func (ps *ParticleSystem) DoStep() {
	ps.VertexArray.Bind()
	ps.SetUniformsAndProgram()

	gl.Enable(gl.RASTERIZER_DISCARD)
	defer gl.Disable(gl.RASTERIZER_DISCARD)

	ps.Data1.Bind(gl.ARRAY_BUFFER)

	ps.Locations.Pos1.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Pos1))
	ps.Locations.Pos2.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Pos2))
	ps.Locations.Lifetime.AttribPointer(1, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Lifetime))

	time := glfw.GetTime()
	ps.Locations.Origin.Uniform3f(100*float32(math.Sin(time)), 100*float32(math.Cos(time)), 100)

	ps.Data2.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	gl.BeginTransformFeedback(gl.POINTS)
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
	gl.EndTransformFeedback()

	tmp := ps.Data1
	ps.Data1 = ps.Data2
	ps.Data2 = tmp

}

func (ps *ParticleSystem) Render(matrix *mathgl.Mat4f) {
	ps.VertexArray.Bind()

	ps.RenderProgram.Use()
	ps.RenderLocations.Pos1.EnableArray()
	ps.RenderLocations.Lifetime.EnableArray()
	ps.RenderLocations.Matrix.UniformMatrix4f(false, (*[16]float32)(matrix))
	ps.RenderLocations.MaxLifetime.Uniform1f(ps.MaxLifetime)

	ps.Data1.Bind(gl.ARRAY_BUFFER)
	ps.RenderLocations.Pos1.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Pos1))
	ps.RenderLocations.Lifetime.AttribPointer(1, gl.FLOAT, false, int(unsafe.Sizeof(ParticleVertex{})), unsafe.Offsetof(ParticleVertex{}.Lifetime))
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
}

func (ps *ParticleSystem) Delete() {
	ps.VertexArray.Delete()
	ps.Data1.Delete()
	ps.Data2.Delete()
	ps.Program.Delete()
	ps.RenderProgram.Delete()
}
