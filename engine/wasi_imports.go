package engine

import (
	"encoding/binary"

	"github.com/wasmerio/wasmer-go/wasmer"
)

func (w *Wasm) fdWriteFunc() *wasmer.Function {
	return wasmer.NewFunction(
		w.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			// 0: file descriptor
			// 1: pointer to data in wasm memory
			// 2: length
			// 3: characters written pointer

			//fmt.Printf("fd:%d iovs_ptr:%d iovs_len:%d nwritten:%d\n", args[0].I32(), args[1].I32(), args[2].I32(), args[3].I32())

			mem, _ := w.instance.Exports.GetMemory("memory")
			ptr := args[1].I32() + 8

			offset := binary.LittleEndian.Uint32(mem.Data()[ptr:])
			l := binary.LittleEndian.Uint32(mem.Data()[ptr+4:])
			data := mem.Data()[offset : offset+l]

			switch args[0].I32() {
			case 1:
				w.log.Info("StdOut.Write from Wasm module", "message", string(data))
			case 2:
				w.log.Error("StdErr.Write from Wasm module", "message", string(data))
			default:
				panic("File writing not implemented")
			}
			//p := unsafe.Pointer(uintptr(args[1].I32()))
			//C.fd_write(args[0].I32(), args[1].I32(), 1)

			return []wasmer.Value{wasmer.NewI32(0)}, nil
		})
}

func (w *Wasm) addWasi(importObject *wasmer.ImportObject) {
	// tinygo
	importObject.Register(
		"wasi_unstable",
		map[string]wasmer.IntoExtern{
			"fd_write": w.fdWriteFunc(),
		},
	)

	// assemblyscript
	importObject.Register(
		"wasi_snapshot_preview1",
		map[string]wasmer.IntoExtern{
			"fd_write": w.fdWriteFunc(),
		},
	)
}
