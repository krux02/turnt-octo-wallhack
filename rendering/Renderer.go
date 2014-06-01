package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
)

type IRenderer interface {
	Render(renderData *RenderData, Proj mgl.Mat4f, View mgl.Mat4f, Model mgl.Mat4f, ClippingPlane_ws mgl.Vec4f)
	UseProgram()
	RenderLocations() *RenderLocations
	Update(entiy gamestate.IRenderEntity, etc interface{})
	Delete()
}

type RenderUpdateFunc func(*RenderLocations, gamestate.IRenderEntity, interface{})
type RenderInitFunc func(*RenderLocations)

type Renderer struct {
	Program              gl.Program
	RenLoc               RenderLocations
	OverrideModeToPoints bool
	UpdateFunc           func(loc *RenderLocations, entity gamestate.IRenderEntity, etc interface{})
}

func NewRenderer(program gl.Program, name string, initFunc RenderInitFunc, updateFunc RenderUpdateFunc) *Renderer {
	this := new(Renderer)

	this.Program = program
	this.Program.Use()

	helpers.BindLocations(name, this.Program, &this.RenLoc)

	if initFunc != nil {
		initFunc(&this.RenLoc)
	}

	this.UpdateFunc = updateFunc
	return this
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

/*
func (this *Renderer) SetUniformByName(name string, value interface{}) {
	location := this.Program.GetUniformLocation(name)
	this.SetUniform(location, value)
}

func (this *Renderer) SetUniform(location gl.UniformLocation, value interface{}) {
	switch Value := value.(type) {
	case int:
		location.Uniform1i(Value)
	case float32:
		location.Uniform1f(Value)
	case mgl.Vec2f:
		location.Uniform2f(Value[0], Value[1])
	case mgl.Vec3f:
		location.Uniform3f(Value[0], Value[1], Value[2])
	case mgl.Vec4f:
		location.Uniform4f(Value[0], Value[1], Value[2], Value[3])
	case mgl.Mat2f:
		location.UniformMatrix2f(false, glMat2(&Value))
	case mgl.Mat3f:
		location.UniformMatrix3f(false, glMat3(&Value))
	case mgl.Mat4f:
		location.UniformMatrix4f(false, glMat4(&Value))
	}
}
*/

func (this *Renderer) Update(entiy gamestate.IRenderEntity, etc interface{}) {
	if this.UpdateFunc != nil {
		this.UpdateFunc(&this.RenLoc, entiy, etc)
	}
	/*
		else if etc != nil {
			switch Map := etc.(type) {
			case map[string]interface{}:
				for name, value := range Map {
					this.SetUniformByName(name, value)
				}
			default:
			}
		}
	*/
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
