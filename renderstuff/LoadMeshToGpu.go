package renderstuff

import (
	"fmt"
	//	mgl "github.com/krux02/mathgl/mgl32"
	"github.com/go-gl-legacy/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"reflect"
)

func LoadMeshToGpu(mesh *Mesh, renLoc *RenderLocations) (rd RenderData) {
	rd.VAO = gl.GenVertexArray()
	rd.VAO.Bind()

	{
		vertices := mesh.Vertices
		verticesType := reflect.TypeOf(vertices)
		if verticesType.Kind() != reflect.Slice {
			panic("Vertices is not a slice")
		}
		rd.Vertices = gl.GenBuffer()
		rd.Vertices.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)
		rd.Numverts = reflect.ValueOf(vertices).Len()
		helpers.SetAttribPointers(renLoc, reflect.ValueOf(vertices).Index(0).Addr().Interface(), false)
		switch mesh.Mode {
		case Points:
			rd.Mode = gl.POINTS
		case LineStrip:
			rd.Mode = gl.LINE_STRIP
		case LineLoop:
			rd.Mode = gl.LINE_LOOP
		case Lines:
			rd.Mode = gl.LINES
		case TriangleStrip:
			rd.Mode = gl.TRIANGLE_STRIP
		case TriangleFan:
			rd.Mode = gl.TRIANGLE_FAN
		case Triangles:
			rd.Mode = gl.TRIANGLES
		default:
			panic("unsupported mode")
		}
	}

	if indices := mesh.Indices; indices != nil {
		indicesType := reflect.TypeOf(indices)
		if indicesType.Kind() != reflect.Slice {
			panic("Indices is not a slice")
		}
		rd.Indices = gl.GenBuffer()
		rd.Indices.Bind(gl.ELEMENT_ARRAY_BUFFER)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, helpers.ByteSizeOfSlice(indices), indices, gl.STATIC_DRAW)
		rd.Numverts = reflect.ValueOf(indices).Len()
		switch indicesType.Elem().Kind() {
		case reflect.Uint8, reflect.Int8:
			rd.IndexType = gl.UNSIGNED_BYTE
		case reflect.Uint16, reflect.Int16:
			rd.IndexType = gl.UNSIGNED_SHORT
		case reflect.Uint32, reflect.Int32:
			rd.IndexType = gl.UNSIGNED_INT
		default:
			panic(fmt.Sprint("unsupported index type", indicesType.Elem().Kind()))
		}
	}

	if instanceData := mesh.InstanceData; instanceData != nil {
		Type := reflect.TypeOf(instanceData)
		if Type.Kind() != reflect.Slice {
			panic("InstanceData is not a slice")
		}
		rd.InstanceData = gl.GenBuffer()
		rd.InstanceData.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(instanceData), instanceData, gl.STATIC_DRAW)
		helpers.SetAttribPointers(renLoc, reflect.ValueOf(instanceData).Index(0).Addr().Interface(), true)

		rd.NumInstances = reflect.ValueOf(instanceData).Len()
	}
	return
}
