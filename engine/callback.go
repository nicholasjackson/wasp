package engine

import (
	"reflect"

	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/niemeyer/pretty"
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

func (c *Callbacks) merge(cb *Callbacks) {
	for cb, fs := range cb.callbackFunctions {
		for name, f := range fs {
			c.AddCallback(cb, name, f)
		}
	}

	pretty.Println(c)
}

func (c *Callbacks) addCallbacks(i *Instance, store *wasmer.Store, log *logger.Wrapper) {
	for ns, fs := range c.callbackFunctions {
		callbacks := map[string]wasmer.IntoExtern{}

		for name, f := range fs {
			ft, ff := createCallback(i, ns, name, f)
			callbacks[name] = wasmer.NewFunction(store, ft, ff)
		}

		i.importObject.Register(ns, callbacks)
	}
}

func createCallback(i *Instance, ns, name string, callFunc interface{}) (*wasmer.FunctionType, func([]wasmer.Value) ([]wasmer.Value, error)) {
	callback := reflect.TypeOf(callFunc)
	callback.In(0)

	inParams := []wasmer.ValueKind{}
	for i := 0; i < callback.NumIn(); i++ {
		inParams = append(inParams, wasmer.I32)
	}

	outParams := []wasmer.ValueKind{}
	for i := 0; i < callback.NumOut(); i++ {
		outParams = append(outParams, wasmer.I32)
	}

	ft := wasmer.NewFunctionType(
		wasmer.NewValueTypes(inParams...),
		wasmer.NewValueTypes(outParams...))

	ff := func(args []wasmer.Value) ([]wasmer.Value, error) {

		i.log.Debug("Callback called", "namespace", ns, "name", name, "args", args)

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
		f := reflect.ValueOf(callFunc)
		out := f.Call(inParams)

		i.log.Debug("Called callback function", "out", out)

		// check returned parameters = expected
		if len(out) != callback.NumOut() {
			return nil, xerrors.Errorf(
				"callback function %s.%s received with incorrect number of return parameters, expected: %d, received: %d, signature: %s",
				ns,
				name,
				callback.NumOut(),
				len(out),
				callback.String(),
			)
		}

		// process the response parameters
		outParams := []wasmer.Value{}
		for n := 0; n < callback.NumOut(); n++ {
			switch callback.Out(n).Kind() {
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
	}

	return ft, ff
}
