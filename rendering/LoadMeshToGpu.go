package rendering

import (
	"fmt"
	//	mgl "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/renderstuff"
	"reflect"
)

func LoadMeshToGpu(mesh renderstuff.IMesh, renLoc *RenderLocations) (rd RenderData) {
	rd.VAO = gl.GenVertexArray()
	rd.VAO.Bind()

	{
		vertices := mesh.Vertices()
		verticesType := reflect.TypeOf(vertices)
		if verticesType.Kind() != reflect.Slice {
			panic("Vertices is not a slice")
		}
		rd.Vertices = gl.GenBuffer()
		rd.Vertices.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, helpers.ByteSizeOfSlice(vertices), vertices, gl.STATIC_DRAW)
		rd.Numverts = reflect.ValueOf(vertices).Len()
		helpers.SetAttribPointers(renLoc, reflect.ValueOf(vertices).Index(0).Addr().Interface(), false)
		switch mesh.Mode() {
		case renderstuff.Points:
			rd.Mode = gl.POINTS
		case renderstuff.LineStrip:
			rd.Mode = gl.LINE_STRIP
		case renderstuff.LineLoop:
			rd.Mode = gl.LINE_LOOP
		case renderstuff.Lines:
			rd.Mode = gl.LINES
		case renderstuff.TriangleStrip:
			rd.Mode = gl.TRIANGLE_STRIP
		case renderstuff.TriangleFan:
			rd.Mode = gl.TRIANGLE_FAN
		case renderstuff.Triangles:
			rd.Mode = gl.TRIANGLES
		}
	}

	if indices := mesh.Indices(); indices != nil {
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

	if instanceData := mesh.InstanceData(); instanceData != nil {
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
