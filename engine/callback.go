package engine

import (
	"fmt"
	"reflect"

	"github.com/wasmerio/wasmer-go/wasmer"
)

// AddCallback exposes add a function that can be called from the Wasm module
func (w *Wasm) AddCallback(name string, f interface{}) {
	w.callbackFunctions[name] = f
}

func (w *Wasm) addCallbacks(namespace string) {
	callbacks := map[string]wasmer.IntoExtern{}

	for name, f := range w.callbackFunctions {
		callback := reflect.TypeOf(f)
		callback.In(0)

		inParams := []wasmer.ValueKind{}
		for i := 0; i < callback.NumIn(); i++ {
			inParams = append(inParams, wasmer.I32)
		}

		outParams := []wasmer.ValueKind{}
		for i := 0; i < callback.NumOut(); i++ {
			outParams = append(outParams, wasmer.I32)
		}

		callbacks[name] = wasmer.NewFunction(
			w.store,
			wasmer.NewFunctionType(wasmer.NewValueTypes(inParams...), wasmer.NewValueTypes(outParams...)),
			func(args []wasmer.Value) ([]wasmer.Value, error) {

				w.log.Info("Callback called")

				// build the parameter list
				inParams := []reflect.Value{}
				for i := 0; i < callback.NumIn(); i++ {
					switch callback.In(i).Kind() {
					case reflect.String:
						in, err := w.getStringFromMemory(args[i].I32())
						if err != nil {
							panic(err)
						}

						ps := reflect.ValueOf(in)
						inParams = append(inParams, ps)
					case reflect.Int32:
						ps := reflect.ValueOf(args[i].I32())
						inParams = append(inParams, ps)

					default:
						return nil, fmt.Errorf("Only String and Int32 parameters can be used for callback functions")
					}
				}

				// call the function
				f := reflect.ValueOf(f)
				out := f.Call(inParams)
				w.log.Debug("Called callback function", "out", out)

				// process the response parameters
				outParams := []wasmer.Value{}
				for i := 0; i < callback.NumOut(); i++ {
					switch callback.In(i).Kind() {
					case reflect.String:
						s, err := w.setStringInMemory(out[i].String())
						if err != nil {
							panic(err)
						}

						outParams = append(outParams, wasmer.NewI32(s))
					case reflect.Int32:
						outParams = append(outParams, wasmer.NewI32(out[i].Int()))

					default:
						return nil, fmt.Errorf("Only String and Int32 parameters can be used for callback functions")
					}
				}

				return outParams, nil
			},
		)
	}

	w.importObject.Register(namespace, callbacks)
}
