package main

/*
#cgo linux LDFLAGS: -lGL -lGLEW
#include "initDebugContext.h"
*/
import "C"

import (
	"fmt"
	"github.com/go-gl/gl"
	"unsafe"
)

func initDebugContext() bool {
	return C.initDebugContext() == 1
}

func PointParameterf(pname gl.GLenum, param float32) {
	C.glPointParameterf(C.GLenum(pname), C.GLfloat(param))
}

func PointParameterfv(pname gl.GLenum, params []float32) {

	if len(params) == 0 {
		panic("Invalid slice length")
	}

	C.glPointParameterfv(C.GLenum(pname), (*C.GLfloat)(&params[0]))
}

const (
	POINT_DISTANCE_ATTENUATION_ARB = C.GL_POINT_DISTANCE_ATTENUATION_ARB
	POINT_SIZE_MAX_ARB             = C.GL_POINT_SIZE_MAX_ARB
	POINT_SIZE_MIN_ARB             = C.GL_POINT_SIZE_MIN_ARB
	POINT_SPRITE_ARB = C.GL_POINT_SPRITE_ARB
)

//export goDebugCallback
func goDebugCallback(source, _type, id, severity C.uint, length C.int, message *C.char, userParam unsafe.Pointer) {
	switch source {
	case 0x8246:
		fmt.Println("DEBUG_SOURCE_API_ARB")
	case 0x8247:
		fmt.Println("DEBUG_SOURCE_WINDOW_SYSTEM_ARB")
	case 0x8248:
		fmt.Println("DEBUG_SOURCE_SHADER_COMPILER_ARB")
	case 0x8249:
		fmt.Println("DEBUG_SOURCE_THIRD_PARTY_ARB")
	case 0x824A:
		fmt.Println("DEBUG_SOURCE_APPLICATION_ARB")
	case 0x824B:
		fmt.Println("DEBUG_SOURCE_OTHER_ARB")
	}

	switch _type {
	case 0x824C:
		fmt.Println("DEBUG_TYPE_ERROR_ARB")
	case 0x824D:
		fmt.Println("DEBUG_TYPE_DEPRECATED_BEHAVIOR_ARB")
	case 0x824E:
		fmt.Println("DEBUG_TYPE_UNDEFINED_BEHAVIOR_ARB")
	case 0x824F:
		fmt.Println("DEBUG_TYPE_PORTABILITY_ARB")
	case 0x8250:
		fmt.Println("DEBUG_TYPE_PERFORMANCE_ARB")
	case 0x8251:
		fmt.Println("DEBUG_TYPE_OTHER_ARB")
	}

	switch severity {
	case 0x9146:
		fmt.Println("DEBUG_SEVERITY_HIGH_ARB")
	case 0x9147:
		fmt.Println("DEBUG_SEVERITY_MEDIUM_ARB")
	case 0x9148:
		fmt.Println("DEBUG_SEVERITY_LOW_ARB")
	}

	fmt.Printf("id: %d\n", id)
	panic(C.GoStringN(message, length))
}

type GLerror gl.GLenum

func (err GLerror) Error() string {
	switch err {
	case gl.NO_ERROR:
		return "NO_ERROR"
	case gl.INVALID_ENUM:
		return "INVALID_ENUM"
	case gl.INVALID_VALUE:
		return "INVALID_VALUE"
	case gl.INVALID_OPERATION:
		return "INVALID_OPERATION"
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		return "INVALID_FRAMEBUFFER_OPERATION"
	case gl.OUT_OF_MEMORY:
		return "OUT_OF_MEMORY"
	case gl.STACK_UNDERFLOW:
		return "STACK_UNDERFLOW"
	case gl.STACK_OVERFLOW:
		return "STACK_OVERFLOW"
	}
	panic(fmt.Sprintf("invalid GLerror: %d", err))
}
