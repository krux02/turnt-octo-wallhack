package main

import (
	"github.com/krux02/tw"
	"reflect"
	"unsafe"
)

func MakeGui(bar *tw.Bar, options interface{}) {
	elem := reflect.ValueOf(options).Elem()
	typ := reflect.TypeOf(options).Elem()

	if elem.Kind() != reflect.Struct {
		panic("wrong arguments need struct pointer")
	}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		structField := typ.Field(i)
		if field.Kind() == reflect.Bool {
			dataPtr := unsafe.Pointer(field.Pointer())
			Name := structField.Name
			bar.AddVarRW(Name, tw.TYPE_BOOL8, dataPtr, "")
		}
	}

}
