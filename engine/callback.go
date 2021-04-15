package engine

import (
	"reflect"

	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/wasmerio/wasmer-go/wasmer"
	"golang.org/x/xerrors"
)

type Callbacks struct {
	callbackFunctions map[string]map[string]interface{}
}

// AddCallback exposes add a function that can be called from the Wasm module
func (c *Callbacks) AddCallback(namespace string, name string, f interface{}) {
	// ensure the collection is instantiated
	if c.callbackFunctions == nil {
		c.callbackFunctions = make(map[string]map[string]interface{})
	}

	// does the namespace exist
	if _, ok := c.callbackFunctions[namespace]; !ok {
		c.callbackFunctions[namespace] = make(map[string]interface{})
	}

	c.callbackFunctions[namespace][name] = f
}

func (w *Callbacks) addCallbacks(i *Instance, store *wasmer.Store, log *logger.Wrapper) {
	for ns, fs := range w.callbackFunctions {
		callbacks := map[string]wasmer.IntoExtern{}

		for name, f := range fs {
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
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(inParams...), wasmer.NewValueTypes(outParams...)),
				func(args []wasmer.Value) ([]wasmer.Value, error) {

					log.Debug("Callback called", "namespace", ns, "name", name)

					// ensure the deallocation of memory is always gets called, pass a reference as the slice is not yet populated
					defer i.freeAllocatedMemory()

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
							return nil, xerrors.Errorf("only String and Int32 parameters can currently be used for callback functions")
						}
					}

					// call the function
					f := reflect.ValueOf(f)
					out := f.Call(inParams)

					log.Debug("Called callback function", "out", out)

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
							return nil, xerrors.Errorf("only String and Int32 parameters can be used for callback functions")
						}
					}

					return outParams, nil
				},
			)
		}

		i.importObject.Register(ns, callbacks)
	}
}
