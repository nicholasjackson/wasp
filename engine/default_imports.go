package engine

// #include <stdlib.h>
/*

extern void debug(void *context, int32_t a);
*/
import "C"

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/wasmerio/wasmer-go/wasmer"
)

//export debug
func debug(_ unsafe.Pointer, sp int32) {
	log.Println(sp)
}

func addDefaults(importObject *wasmer.ImportObject, store *wasmer.Store) {
	importObject.Register(
		"wasi_unstable",
		map[string]wasmer.IntoExtern{
			"fd_write": wasmer.NewFunction(
				store,
				wasmer.NewFunctionType(
					wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes(wasmer.I32),
				),
				func(args []wasmer.Value) ([]wasmer.Value, error) {
					fmt.Printf("%d %d %d %d\n", args[0].I32(), args[1].I32(), args[2].I32(), args[3].I32())
					p := unsafe.Pointer(uintptr(args[1].I32()))
					C.debug(p, C.int(args[3].I32()))
					return []wasmer.Value{wasmer.NewI32(0)}, nil
				},
			),
		},
	)

	importObject.Register(
		"env",
		map[string]wasmer.IntoExtern{
			"abort": wasmer.NewFunction(
				store,
				wasmer.NewFunctionType(
					wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes(),
				),
				func(args []wasmer.Value) ([]wasmer.Value, error) {
					return []wasmer.Value{}, nil
				},
			),
		},
	)
}
