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

func (w *Wasm) addCallbacks(namespace string, i *Instance) {
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
				for n := 0; n < callback.NumIn(); n++ {
					switch callback.In(n).Kind() {
					case reflect.String:
						in, err := i.getStringFromMemory(args[n].I32())
						if err != nil {
							panic(err)
						}

						ps := reflect.ValueOf(in)
						inParams = append(inParams, ps)
					case reflect.Int32:
						ps := reflect.ValueOf(args[n].I32())
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
				for n := 0; n < callback.NumOut(); n++ {
					switch callback.In(n).Kind() {
					case reflect.String:
						s, err := i.setStringInMemory(out[n].String())
						if err != nil {
							panic(err)
						}

						outParams = append(outParams, wasmer.NewI32(s))
					case reflect.Int32:
						outParams = append(outParams, wasmer.NewI32(out[n].Int()))

					default:
						return nil, fmt.Errorf("Only String and Int32 parameters can be used for callback functions")
					}
				}

				return outParams, nil
			},
		)
	}

	i.importObject.Register(namespace, callbacks)
}
