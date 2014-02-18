package settings

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"os"
	"reflect"
)

type BoolOptions struct {
	NoParticleRender,
	NoParticlePhysics,
	NoWorldRender,
	NoTreeRender,
	NoPlayerPhysics,
	DepthClamp,
	Wireframe bool
	StartPosition mgl.Vec4f
}

func (this *BoolOptions) Load() {
	file, err := os.Open("settings.txt")
	if err != nil {
		panic(err)
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
		}
	}
}

func (this *BoolOptions) Save() {
	file, _ := os.Create("settings.txt")
	defer file.Close()

	//fmt.Fprintf(file, "%v\n", this)

	v := reflect.ValueOf(this).Elem()
	t := v.Type()
	N := v.NumField()
	for i := 0; i < N; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Type.Kind() {
		case reflect.Bool:
			fmt.Fprintf(file, "%s %t\n", field.Name, fieldValue.Bool())
		}
	}

}
