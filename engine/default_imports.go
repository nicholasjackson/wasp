package engine

import (
	"github.com/wasmerio/wasmer-go/wasmer"
)

func (w *Wasm) addDefaults(i *Instance) {

	i.importObject.Register(
		"env",
		map[string]wasmer.IntoExtern{
			"abort": wasmer.NewFunction(
				w.store,
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
