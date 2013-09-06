package main

import (
	"reflect"
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
	"io/ioutil"
)

import "fmt"


func ReadShaderFile(name string) string {
	name = fmt.Sprintf("shaders/%s", name)

	source, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("can't read file", name)
		panic(err)
	}

	return string(source)
}

func MakeProgram(vertFname, fragFname string) gl.Program {
	vertSource := ReadShaderFile(vertFname)
	fragSource := ReadShaderFile(fragFname)
	
	return glh.NewProgram(glh.Shader{gl.VERTEX_SHADER, vertSource}, glh.Shader{gl.FRAGMENT_SHADER, fragSource})
}

func BindLocations(prog gl.Program, locations interface{})  {
	value := reflect.ValueOf(locations).Elem()
	Type := reflect.TypeOf(locations).Elem()

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := Type.Field(i)
		fieldName := fieldType.Name

		switch fieldValue.Interface().(type) {
		case gl.AttribLocation:
			loc := prog.GetAttribLocation(fieldName)
			fieldValue.SetInt(int64(loc))
		case gl.UniformLocation:
			loc := prog.GetUniformLocation(fieldName)
			fieldValue.SetInt(int64(loc))
		default:
		}
	}
}

func LocationMap(locations interface{}) (map[string]gl.AttribLocation, map[string]gl.UniformLocation) {
	attribs := make(map[string]gl.AttribLocation)
	uniforms := make(map[string]gl.UniformLocation)

	value := reflect.ValueOf(locations).Elem()
	Type := reflect.TypeOf(locations).Elem()

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := Type.Field(i)
		fieldName := fieldType.Name

		switch fieldValue.Interface().(type) {
		case gl.AttribLocation:
			loc := gl.AttribLocation(fieldValue.Int())
			attribs[fieldName] = loc
		case gl.UniformLocation:
			loc := gl.UniformLocation(fieldValue.Int())
			uniforms[fieldName] = loc
		default:
		}
	}

	return attribs,uniforms
}

func SetAttribPointers(locations interface{}, vertexData interface{}) {
	attribs,_ := LocationMap(locations)

	Type := reflect.TypeOf(vertexData).Elem()
	stride := int(Type.Size())

	for i := 0; i < Type.NumField(); i++ {
		structElement := Type.Field(i)
		elementType := structElement.Type

		var size uint
		var typ gl.GLenum
		var kind reflect.Kind

		switch elementType.Kind() {
		case reflect.Array:
			size = uint(elementType.Len())
			kind = elementType.Elem().Kind()
		default:
			size = 1
			kind = elementType.Kind()
		}

		switch kind {
		case reflect.Float32:
			typ = gl.FLOAT
		case reflect.Float64:
			typ = gl.DOUBLE
		default:
			panic("not implemented yet")
		}

		offset := structElement.Offset
		Loc := attribs[structElement.Name]
		Loc.EnableArray()

		// fmt.Printf("Loc: %d, size: %d, type: %d, stride: %d, offset: %d\n", Loc, size, typ, stride, offset)
		Loc.AttribPointer(size, typ, false, stride, offset)
	}
}
