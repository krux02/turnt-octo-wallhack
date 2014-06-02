package rendering

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/constants"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
	"reflect"
)

type RenderUpdateFunc func(*RenderLocations, renderstuff.IRenderEntity, interface{})
type RenderInitFunc func(*RenderLocations)

type Renderer struct {
	Program              gl.Program
	RenLoc               RenderLocations
	OverrideModeToPoints bool
	UpdateFunc           RenderUpdateFunc
	RenData              map[*renderstuff.Mesh]*RenderData
}

func NewRenderer(program gl.Program, name string, initFunc RenderInitFunc, updateFunc RenderUpdateFunc) *Renderer {
	this := new(Renderer)

	this.Program = program
	this.Program.Use()

	helpers.BindLocations(name, this.Program, &this.RenLoc)

	if initFunc != nil {
		initFunc(&this.RenLoc)
	} else {
		this.setImageUniforms()
	}

	this.UpdateFunc = updateFunc
	this.RenData = map[*renderstuff.Mesh]*RenderData{}
	return this
}

func (this *Renderer) setImageUniforms() {
	locType := reflect.TypeOf(this.RenLoc)
	locVal := reflect.ValueOf(this.RenLoc)

	for i := 0; i < locType.NumField(); i++ {
		fieldName := locType.Field(i).Name
		switch uniform := locVal.Field(i).Interface().(type) {
		case gl.UniformLocation:
			if val, ok := constants.Texture[fieldName]; ok {
				uniform.Uniform1i(val)
			}
		}
	}
}

func (this *Renderer) Delete() {
	this.Program.Delete()
	for _, value := range this.RenData {
		value.Delete()
	}
	*this = Renderer{}
}

func (this *Renderer) Render(entity renderstuff.IRenderEntity, Proj, View mgl.Mat4f, ClippingPlane_ws mgl.Vec4f, additionalUniforms interface{}) {
	this.Program.Use()
	mesh := entity.Mesh()
	renData := this.RenData[mesh]
	if renData == nil {
		md := LoadMeshToGpu(mesh, &this.RenLoc)
		renData = &md
		this.RenData[mesh] = &md
	}
	Model := entity.Model()

	renData.VAO.Bind()

	if this.UpdateFunc != nil {
		this.UpdateFunc(&this.RenLoc, entity, additionalUniforms)
	}

	Loc := this.RenLoc
	Loc.View.UniformMatrix4f(false, glMat4(&View))
	Loc.Model.UniformMatrix4f(false, glMat4(&Model))
	Loc.Proj.UniformMatrix4f(false, glMat4(&Proj))
	Loc.ClippingPlane_ws.Uniform4f(ClippingPlane_ws[0], ClippingPlane_ws[1], ClippingPlane_ws[2], ClippingPlane_ws[3])

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

	Draw(renData)
}

func Draw(renData *RenderData) {
	numverts := renData.Numverts
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
