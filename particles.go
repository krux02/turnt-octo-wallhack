package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
	"github.com/krux02/mathgl"
	"io/ioutil"
	"math"
	"math/rand"
)

type ParticleVertex struct {
	Pos1     mathgl.Vec3f
	Pos2     mathgl.Vec3f
	Lifetime float32
}

type ProgramLocations struct {
	Pos1, Pos2, Lifetime            gl.AttribLocation
	Origin, Gravity, MaxLifetime    gl.UniformLocation
	Heights, LowerBound, UpperBound gl.UniformLocation
	RandomDirs                      gl.UniformLocation
}

type RenderProgramLocations struct {
	Pos1, Pos2, Lifetime gl.AttribLocation
	Matrix, MaxLifetime  gl.UniformLocation
}

type ParticleSystem struct {
	TransformProg                            gl.Program
	TransformLoc                             ProgramLocations
	RenderProg                               gl.Program
	RenderLoc                                RenderProgramLocations
	VaoTff1, VaoTff2, VaoRender1, VaoRender2 gl.VertexArray
	Data1                                    gl.Buffer
	Data2                                    gl.Buffer
	NumParticles                             int
	Origin                                   mathgl.Vec3f
	Gravity                                  float32
	InitialSpeed                             float32
	MaxLifetime                              float32
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

	buffer1.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, ByteSizeOfSlice(vertices), vertices, gl.STREAM_DRAW)

	buffer2.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, ByteSizeOfSlice(vertices), uintptr(0), gl.STREAM_DRAW)

	TransformProg := gl.CreateProgram()

	content, err := ioutil.ReadFile("shaders/ParticleTFF.vs")
	if err != nil {
		panic(err)
	}

	shader := glh.MakeShader(gl.VERTEX_SHADER, string(content))
	TransformProg.AttachShader(shader)
	TransformProg.TransformFeedbackVaryings([]string{"v_Pos1", "v_Pos2", "v_Lifetime"}, gl.INTERLEAVED_ATTRIBS)
	TransformProg.Link()
	shader.Delete()

	TransformProg.Use()

	TransformLoc := ProgramLocations{}
	BindLocations(TransformProg, &TransformLoc)

	renderProgram := MakeProgram("Particle.vs", "Particle.fs")
	renderProgram.Use()
	RenderLoc := RenderProgramLocations{}
	BindLocations(renderProgram, &RenderLoc)

	vaoTff1 := gl.GenVertexArray()
	vaoTff2 := gl.GenVertexArray()
	vaoRender1 := gl.GenVertexArray()
	vaoRender2 := gl.GenVertexArray()

	ps := &ParticleSystem{TransformProg, TransformLoc, renderProgram, RenderLoc, vaoTff1, vaoTff2, vaoRender1, vaoRender2, buffer1, buffer2, numParticles, Origin, -9.81 / 200, initialSpeed, MaxLifetime}

	TransformProg.Use()
	ps.SetUniforms()

	return ps
}

func (ps *ParticleSystem) SetVaos() {
	ps.TransformProg.Use()

	ps.VaoTff1.Bind()
	ps.Data1.Bind(gl.ARRAY_BUFFER)
	SetAttribPointers(&ps.TransformLoc, &ParticleVertex{}, false)
	ps.Data2.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	ps.VaoTff2.Bind()
	ps.Data2.Bind(gl.ARRAY_BUFFER)
	SetAttribPointers(&ps.TransformLoc, &ParticleVertex{}, false)
	ps.Data1.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	ps.RenderProg.Use()

	ps.VaoRender1.Bind()
	ps.Data1.Bind(gl.ARRAY_BUFFER)
	SetAttribPointers(&ps.RenderLoc, &ParticleVertex{}, false)

	ps.VaoRender2.Bind()
	ps.Data2.Bind(gl.ARRAY_BUFFER)
	SetAttribPointers(&ps.RenderLoc, &ParticleVertex{}, false)
}

func (ps *ParticleSystem) SetUniforms() {
	ps.TransformLoc.Origin.Uniform3f(ps.Origin[0], ps.Origin[1], ps.Origin[2])
	ps.TransformLoc.Gravity.Uniform1f(ps.Gravity)
	ps.TransformLoc.MaxLifetime.Uniform1f(ps.MaxLifetime)

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

	ps.TransformLoc.RandomDirs.Uniform3fv(64, dirs)
}

func (ps *ParticleSystem) DoStep(time float64) {
	ps.TransformProg.Use()
	ps.VaoTff1.Bind()

	gl.Enable(gl.RASTERIZER_DISCARD)
	defer gl.Disable(gl.RASTERIZER_DISCARD)

	ps.Data1.Bind(gl.ARRAY_BUFFER)

	SetAttribPointers(&ps.TransformLoc, &ParticleVertex{}, false)

	ps.TransformLoc.Origin.Uniform3f(100*float32(math.Sin(time)), 100*float32(math.Cos(time)), 100)
	//ps.TransformLoc.Origin.Uniform3f(0, 0, 100)

	ps.Data2.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	gl.BeginTransformFeedback(gl.POINTS)
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
	gl.EndTransformFeedback()

	tmp := ps.Data1
	ps.Data1 = ps.Data2
	ps.Data2 = tmp
}

func (ps *ParticleSystem) Render(matrix *mathgl.Mat4f) {
	ps.VaoTff1.Bind()

	ps.RenderProg.Use()
	Loc := ps.RenderLoc

	Loc.Matrix.UniformMatrix4f(false, (*[16]float32)(matrix))
	Loc.MaxLifetime.Uniform1f(ps.MaxLifetime)

	
	ps.Data1.Bind(gl.ARRAY_BUFFER)
	SetAttribPointers(&ps.RenderLoc, &ParticleVertex{}, false)
	gl.DepthMask(false)
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
	gl.DepthMask(true)
}

func (ps *ParticleSystem) Delete() {
	ps.VaoTff1.Delete()
	ps.VaoTff2.Delete()
	ps.VaoRender1.Delete()
	ps.VaoRender2.Delete()
	ps.Data1.Delete()
	ps.Data2.Delete()
	ps.TransformProg.Delete()
	ps.RenderProg.Delete()
}
