package helpers

import (
	"fmt"
	"github.com/go-gl/gl"
	"io/ioutil"
	"log"
	"reflect"
)

const (
	VERTEX_INFO_FLAG = 1
)

func ReadShaderFile(name string) string {
	name = fmt.Sprintf("shaders/%s", name)

	source, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("can't read file", name)
		panic(err)
	}

	return string(source)
}

func MakeShader(type_ gl.GLenum, filename string) gl.Shader {
	source := ReadShaderFile(filename)
	shader := gl.CreateShader(type_)
	shader.Source(source)

	shader.Compile()
	log := shader.GetInfoLog()
	if log != "" {
		panic(fmt.Sprint(filename, log))
	}

	return shader
}

func MakeProgram3(vertFname, geomFname, fragFname string) gl.Program {
	vertShader := MakeShader(gl.VERTEX_SHADER, vertFname)
	defer vertShader.Delete()
	geomShader := MakeShader(gl.GEOMETRY_SHADER, geomFname)
	defer geomShader.Delete()
	fragShader := MakeShader(gl.FRAGMENT_SHADER, fragFname)
	defer fragShader.Delete()

	program := gl.CreateProgram()
	program.AttachShader(vertShader)
	program.AttachShader(geomShader)
	program.AttachShader(fragShader)
	program.Link()

	log := program.GetInfoLog()
	if log != "" {
		if program.Get(gl.LINK_STATUS) == gl.FALSE {
			panic(fmt.Sprint("linking ", vertFname, geomFname, fragFname, log))
		} else {
			fmt.Print("linking ", vertFname, geomFname, fragFname, log)
		}
	}

	return program
}

func MakeProgram(vertFname, fragFname string) gl.Program {
	vao := gl.GenVertexArray()
	vao.Bind()

	vertShader := MakeShader(gl.VERTEX_SHADER, vertFname)
	defer vertShader.Delete()
	fragShader := MakeShader(gl.FRAGMENT_SHADER, fragFname)
	defer fragShader.Delete()

	program := gl.CreateProgram()
	program.AttachShader(vertShader)
	program.AttachShader(fragShader)
	program.Link()

	linkstat := program.Get(gl.LINK_STATUS)
	if linkstat != 1 {
		// log.Panic("Program link failed, sources=", vertFname, fragFname, "\nstatus=", linkstat, "\nInfo log: ", program.GetInfoLog())
	}

	program.Validate()
	valstat := program.Get(gl.VALIDATE_STATUS)
	if valstat != 1 {
		// log.Panic("Program validation failed: ", valstat)
	}

	log := program.GetInfoLog()
	if log != "" {
		if program.Get(gl.LINK_STATUS) == gl.FALSE {
			panic(fmt.Sprint("linking ", vertFname, fragFname, log))
		} else {
			fmt.Print("linking ", vertFname, fragFname, log)
		}
	}

	vao.Delete()

	return program
}

func ByteSizeOfSlice(slice interface{}) int {
	value := reflect.ValueOf(slice)
	typ := reflect.TypeOf(slice)
	if typ.Kind() != reflect.Slice {
		panic("slice is not a slice")
	}
	size := value.Len() * int(typ.Elem().Size())
	return size
}

func BindLocations(name string, prog gl.Program, locations interface{}) {
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
			if loc == -1 {
				log.Println("warning ", fieldName, " in ", name, " not valid")
			}
			fieldValue.SetInt(int64(loc))
		default:
		}
	}
}

func PrintLocations(logger log.Logger, locations interface{}) {
	value := reflect.ValueOf(locations).Elem()
	typ := reflect.TypeOf(locations).Elem()
	logger.Printf("%s:\n", typ.Name())
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldName := typ.Field(i).Name
		logger.Printf("\t%s:\t%d\n", fieldName, fieldValue.Int())
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

	return attribs, uniforms
}

func SetAttribPointers(locations interface{}, vertexData interface{}, instanced bool) {
	attribs, _ := LocationMap(locations)
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
		if Loc >= 0 {
			Loc.EnableArray()
			Loc.AttribPointer(size, typ, false, stride, offset)
			if instanced {
				Loc.AttribDivisor(1)
			}
			if log.Flags()&VERTEX_INFO_FLAG != 0 {
				log.Printf("%s: Loc: %d, size: %d, type: %d, stride: %d, offset: %d\n", structElement.Name, Loc, size, typ, stride, offset)
			}
		} else if log.Flags()&VERTEX_INFO_FLAG != 0 {
			log.Printf("%s: Loc: %d, !!!\n", structElement.Name, Loc)
		}
	}
}
