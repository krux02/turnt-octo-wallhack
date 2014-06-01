package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

type IRenderer interface {
	Render(renderData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f)
	UseProgram()
	RenderLocations() *RenderLocations
	Update(entiy gamestate.IRenderEntity, etc interface{})
	Delete()
}

type Renderer struct {
	Program              gl.Program
	RenLoc               RenderLocations
	OverrideModeToPoints bool
}

func (this *Renderer) Delete() {
	this.Program.Delete()
	*this = Renderer{}
}

func (this *Renderer) UseProgram() {
	this.Program.Use()
}

func (this *Renderer) RenderLocations() *RenderLocations {
	return &this.RenLoc
}

func (this *Renderer) SetUniform(name string, value interface{}) {
	uniform := this.Program.GetUniformLocation(name)
	switch Value := value.(type) {
	case int:
		uniform.Uniform1i(Value)
	case float32:
		uniform.Uniform1f(Value)
	case mgl.Vec2f:
		uniform.Uniform2f(Value[0], Value[1])
	case mgl.Vec3f:
		uniform.Uniform3f(Value[0], Value[1], Value[2])
	case mgl.Vec4f:
		uniform.Uniform4f(Value[0], Value[1], Value[2], Value[3])
	case mgl.Mat2f:
		uniform.UniformMatrix2f(false, glMat2(&Value))
	case mgl.Mat3f:
		uniform.UniformMatrix3f(false, glMat3(&Value))
	case mgl.Mat4f:
		uniform.UniformMatrix4f(false, glMat4(&Value))
	}
}

func (this *Renderer) Update(entiy gamestate.IRenderEntity, etc interface{}) {
	if etc != nil {
		switch Map := etc.(type) {
		case map[string]interface{}:
			for name, value := range Map {
				this.SetUniform(name, value)
			}
		default:
		}
	}
}

func (this *Renderer) Render(renData *RenderData, Proj, View, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f) {
	this.Program.Use()
	renData.VAO.Bind()

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Model.UniformMatrix4f(false, glMat4(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	Loc.ClippingPlane_ws.Uniform4f(ClippingPlane_ws[0], ClippingPlane_ws[1], ClippingPlane_ws[2], ClippingPlane_ws[3])

	numverts := renData.Numverts

	// simple but dirty way to render with points instead
	if this.OverrideModeToPoints {
		oldMode := renData.Mode
		oldIndices := renData.Indices

		renData.Mode = gl.POINTS
		renData.Indices = 0

		defer func() {
			renData.Mode = oldMode
			renData.Indices = oldIndices
		}()
	}

	if renData.InstanceData == 0 {
		if renData.Indices == 0 {
			gl.DrawArrays(renData.Mode, 0, renData.Numverts)
		} else {
			gl.DrawElements(renData.Mode, numverts, renData.IndexType, uintptr(0))
		}
	} else {
		if renData.Indices == 0 {
			gl.DrawArraysInstanced(renData.Mode, 0, renData.Numverts, renData.NumInstances)
		} else {
			gl.DrawElementsInstanced(renData.Mode, renData.Numverts, renData.IndexType, 0, renData.NumInstances)
		}
	}
}
