package gamestate

/*
import (
	"fmt"
	"reflect"
)

var meshNames = []string{"Portal", "Monkey"}


func (this *Meshes) Load() *Meshes {
	v := reflect.ValueOf(this).Elem()
	t := v.Type()
	N := v.NumField()
	for i := 0; i < N; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		name := fmt.Sprintf("meshes/%s.obj", field.Name)
		fieldValue.Set(reflect.ValueOf(LoadMesh(name)))
	}
	return this
}


func LoadMeshes() (meshes map[string]*Mesh) {
	for name := range meshNames {
		meshes[name] = LoadMesh(fmt.Sprintf("meshes/%s.obj", field.Name))
	}
	return meshes
}
*/
