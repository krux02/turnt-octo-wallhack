package settings

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/tw"
	"os"
	"reflect"
	"unsafe"
)

type BoolOptions struct {
	NoParticleRender,
	NoParticlePhysics,
	NoWorldRender,
	NoTreeRender,
	NoPlayerPhysics,
	Wireframe bool
	WaterHeight   float32
	StartPosition mgl.Vec4f
}

func (this *BoolOptions) Load() {
	file, err := os.Open("settings.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	v := reflect.ValueOf(this).Elem()
	var n int

	for err == nil {
		var key, value string
		n, err = fmt.Fscanf(file, "%s %s\n", &key, &value)
		if n == 0 {
			return
		}

		fieldValue := v.FieldByName(key)
		switch fieldValue.Kind() {
		case reflect.Bool:
			var b bool
			fmt.Sscan(value, &b)
			fieldValue.SetBool(b)
		case reflect.Float32, reflect.Float64:
			var f float64
			fmt.Sscan(value, &f)
			fieldValue.SetFloat(f)
		}
	}
}

func (this *BoolOptions) Save() {
	file, _ := os.Create("settings.txt")
	defer file.Close()
	v := reflect.ValueOf(this).Elem()
	t := v.Type()
	N := v.NumField()
	for i := 0; i < N; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Type.Kind() {
		case reflect.Bool:
			fmt.Fprintf(file, "%s %t\n", field.Name, fieldValue.Bool())
		case reflect.Float32, reflect.Float64:
			fmt.Fprintf(file, "%s %f\n", field.Name, fieldValue.Float())
		}
	}

}

func (this *BoolOptions) CreateGui(bar *tw.Bar) {
	v := reflect.ValueOf(this).Elem()
	t := v.Type()
	N := v.NumField()
	for i := 0; i < N; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Type.Kind() {
		case reflect.Bool:
			bar.AddVarRW(field.Name, tw.TYPE_BOOL8, unsafe.Pointer(fieldValue.Addr().Pointer()), "")
		case reflect.Float32:
			bar.AddVarRW(field.Name, tw.TYPE_FLOAT, unsafe.Pointer(fieldValue.Addr().Pointer()), "")
		case reflect.Float64:
			bar.AddVarRW(field.Name, tw.TYPE_DOUBLE, unsafe.Pointer(fieldValue.Addr().Pointer()), "")
		}
	}
}
