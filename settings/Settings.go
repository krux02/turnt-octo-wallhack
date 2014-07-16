package settings

import (
	"fmt"
	"github.com/krux02/tw"
	"os"
	"reflect"
	"unsafe"
)

type BoolOptions struct {
	ParticleRender      bool
	ParticlePhysics     bool
	WorldRender         bool
	WaterRender         bool
	TreeRender          bool
	PlayerPhysics       bool
	Skybox              bool
	Wireframe           bool
	WaterNormals        bool
	DebugLines          bool
	DepthTestDebugLines bool
	PersistentPlayerPos bool
	ShowFramebuffers    bool
	WaterHeight         float32
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
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var f int64
			fmt.Sscan(value, &f)
			fieldValue.SetInt(f)
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
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Fprintf(file, "%s %d\n", field.Name, fieldValue.Int())
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
		case reflect.Int, reflect.Int32, reflect.Int64:
			bar.AddVarRW(field.Name, tw.TYPE_INT32, unsafe.Pointer(fieldValue.Addr().Pointer()), "")
		case reflect.Int16:
			bar.AddVarRW(field.Name, tw.TYPE_INT16, unsafe.Pointer(fieldValue.Addr().Pointer()), "")
		case reflect.Int8:
			bar.AddVarRW(field.Name, tw.TYPE_INT8, unsafe.Pointer(fieldValue.Addr().Pointer()), "")
		}
	}
}
