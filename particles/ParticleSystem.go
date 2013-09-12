package particles

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
	"github.com/krux02/mathgl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"io/ioutil"
	"math"
	"math/rand"
)

type ParticleVertex struct {
	Pos1     mathgl.Vec3f
	Pos2     mathgl.Vec3f
	Lifetime float32
}

type ParticleShapeVertex struct {
	Vertex_os mathgl.Vec4f
	TexCoord  mathgl.Vec2f
}

type ProgramLocations struct {
	Pos1, Pos2, Lifetime            gl.AttribLocation
	Origin, Gravity, MaxLifetime    gl.UniformLocation
	Heights, LowerBound, UpperBound gl.UniformLocation
	RandomDirs                      gl.UniformLocation
}

type RenderProgramLocations struct {
	Pos1, Pos2, Lifetime, TexCoord, Vertex_os gl.AttribLocation
	Proj, View, MaxLifetime, Image            gl.UniformLocation
}

type ParticleSystem struct {
	TransformProg                            gl.Program
	TransformLoc                             ProgramLocations
	RenderProg                               gl.Program
	RenderLoc                                RenderProgramLocations
	VaoTff1, VaoTff2, VaoRender1, VaoRender2 gl.VertexArray
	Data1                                    gl.Buffer
	Data2                                    gl.Buffer
	ShapeData                                gl.Buffer
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
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STREAM_DRAW)

	buffer2.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), uintptr(0), gl.STREAM_DRAW)

	shapeData := CreateShapeDataBuffer()

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
	helpers.BindLocations(TransformProg, &TransformLoc)

	renderProgram := helpers.MakeProgram("Particle.vs", "Particle.fs")
	renderProgram.Use()
	RenderLoc := RenderProgramLocations{}
	helpers.BindLocations(renderProgram, &RenderLoc)

	vaoTff1 := gl.GenVertexArray()
	vaoTff2 := gl.GenVertexArray()
	vaoRender1 := gl.GenVertexArray()
	vaoRender2 := gl.GenVertexArray()

	ps := new(ParticleSystem)
	ps.TransformProg = TransformProg
	ps.TransformLoc = TransformLoc
	ps.RenderProg = renderProgram
	ps.RenderLoc = RenderLoc
	ps.VaoTff1 = vaoTff1
	ps.VaoTff2 = vaoTff2
	ps.VaoRender1 = vaoRender1
	ps.VaoRender2 = vaoRender2
	ps.Data1 = buffer1
	ps.Data2 = buffer2
	ps.ShapeData = shapeData
	ps.NumParticles = numParticles
	ps.Origin = Origin
	ps.Gravity = -9.81 / 200
	ps.InitialSpeed = initialSpeed
	ps.MaxLifetime = MaxLifetime

	TransformProg.Use()
	ps.SetUniforms()
	ps.SetVaos()

	return ps
}

func (ps *ParticleSystem) SetVaos() {
	ps.TransformProg.Use()

	ps.VaoTff1.Bind()
	ps.Data1.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&ps.TransformLoc, &ParticleVertex{}, true)
	ps.Data2.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	ps.VaoTff2.Bind()
	ps.Data2.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&ps.TransformLoc, &ParticleVertex{}, true)
	ps.Data1.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	ps.RenderProg.Use()

	ps.VaoRender1.Bind()
	ps.Data1.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&ps.RenderLoc, &ParticleVertex{}, true)
	ps.RenderLoc.Pos1.AttribDivisor(1)
	ps.RenderLoc.Lifetime.AttribDivisor(1)
	ps.ShapeData.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&ps.RenderLoc, &ParticleShapeVertex{}, true)

	ps.VaoRender2.Bind()
	ps.Data2.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&ps.RenderLoc, &ParticleVertex{}, true)
	ps.RenderLoc.Pos1.AttribDivisor(1)
	ps.RenderLoc.Lifetime.AttribDivisor(1)
	ps.ShapeData.Bind(gl.ARRAY_BUFFER)
	helpers.SetAttribPointers(&ps.RenderLoc, &ParticleShapeVertex{}, true)
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

	// ps.Data1.Bind(gl.ARRAY_BUFFER)
	// SetAttribPointers(&ps.TransformLoc, &ParticleVertex{}, false)

	ps.TransformLoc.Origin.Uniform3f(100*float32(math.Sin(time)), 100*float32(math.Cos(time)), 100)
	ps.Data2.BindBufferBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)

	gl.BeginTransformFeedback(gl.POINTS)
	gl.DrawArrays(gl.POINTS, 0, ps.NumParticles)
	gl.EndTransformFeedback()

	ps.Data1, ps.Data2 = ps.Data2, ps.Data1
	ps.VaoRender1, ps.VaoRender2 = ps.VaoRender2, ps.VaoRender1
	ps.VaoTff1, ps.VaoTff2 = ps.VaoTff2, ps.VaoTff1
}

func CreateShapeDataBuffer() gl.Buffer {
	fmt.Println("CreateShapeDataBuffer:")

	particleShape := []ParticleShapeVertex{
		ParticleShapeVertex{mathgl.Vec4f{-1, -1, 0, 1}, mathgl.Vec2f{0, 1}},
		ParticleShapeVertex{mathgl.Vec4f{-1, 1, 0, 1}, mathgl.Vec2f{0, 0}},
		ParticleShapeVertex{mathgl.Vec4f{1, 1, 0, 1}, mathgl.Vec2f{1, 0}},
		ParticleShapeVertex{mathgl.Vec4f{1, -1, 0, 1}, mathgl.Vec2f{1, 1}},
	}

	particleShapeBuffer := gl.GenBuffer()
	particleShapeBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(particleShape), particleShape, gl.STATIC_DRAW)

	return particleShapeBuffer
}

func (ps *ParticleSystem) Render(Proj mathgl.Mat4f, View mathgl.Mat4f) {
	gl.PointSize(64)

	ps.VaoRender1.Bind()

	ps.RenderProg.Use()
	Loc := ps.RenderLoc

	Loc.Proj.UniformMatrix4f(false, (*[16]float32)(&Proj))
	Loc.View.UniformMatrix4f(false, (*[16]float32)(&View))
	Loc.MaxLifetime.Uniform1f(ps.MaxLifetime)
	Loc.Image.Uniform1i(6)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.DepthMask(false)
	gl.DrawArraysInstanced(gl.TRIANGLE_FAN, 0, 4, ps.NumParticles)
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
