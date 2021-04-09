package engine

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type WASM struct {
	log               hclog.Logger
	instance          *wasmer.Instance
	instanceFunctions *instanceFunctions
}

func New(log hclog.Logger) *WASM {
	return &WASM{log: log}
}

func (w *WASM) LoadPlugin(path string) error {
	wasmBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Unable to load WASM module, error: %s", err)
	}

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	// Add the default imports
	importObject := wasmer.NewImportObject()
	addDefaults(importObject, store)

	// Compile the module
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		return fmt.Errorf("Unable to instantiate WASM module, error: %s", err)
	}

	// Create the new instance of the module
	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return fmt.Errorf("Unable to create a new instance of the WASM module, error: %s", err)
	}
	w.instance = instance

	// Setup the default functions that are required for memory manipulation operations
	wi, err := NewInstanceFunctions(w)
	if err != nil {
		return fmt.Errorf("Unable to import default functions, ensure that the WASM module correctly imports the base ABI, error: %s", err)
	}

	w.instanceFunctions = wi

	return nil
}

// CallFunction in the Wasm module with the given parameters
// The response from the function will automatically be cast into the type specified
// by outputParam. In the instance that outputParam is a complex type that is returned
// as a pointer from the WASMFunction CallFunction reads the WasmModule memory and
// sets outputParam
func (w *WASM) CallFunction(name string, outputParam interface{}, inputParams ...interface{}) error {
	f, err := w.instance.Exports.GetFunction(name)
	if err != nil {
		return fmt.Errorf("Unable to export the WASM function, error: %s", err)
	}

	// parse the input parameters, if we have a string we need to set that in the Wasm modules
	// memory and pass a pointer to the function instead
	processedParams := make([]interface{}, len(inputParams))
	for i, p := range inputParams {
		switch p.(type) {
		case string:
			// we have a string parameter, let's allocate the memory for this in the wasm host and copy
			// the string to it
			addr, err := w.setStringInMemory(p.(string))
			if err != nil {
				return err
			}
			processedParams[i] = addr

		default:
			processedParams[i] = p
		}

	}

	t := time.Now()
	w.log.Debug("Calling function", "name", name, "outputParam", outputParam, "inputParam", processedParams)

	resp, err := f(processedParams...)
	if err != nil {
		return err
	}

	w.log.Debug(
		"Called function",
		"name", name,
		"outputParam", outputParam,
		"inputParam", processedParams,
		"response", resp,
		"time taken", time.Now().Sub(t),
	)

	switch outputParam.(type) {
	case *string:
		s, err := w.getStringFromMemory(resp.(int32))
		if err != nil {
			return err
		}

		*outputParam.(*string) = s

	case *int32:
		*outputParam.(*int32) = resp.(int32)
	default:
		return fmt.Errorf("output parameters can only be of type *int32 or *string")
	}

	return nil
}

func (w *WASM) setStringInMemory(s string) (int32, error) {
	size := len(s) + 1 // allocate 1 more byte than the string size for the null terminator
	addr, err := w.instanceFunctions.Allocate(int32(size))
	if err != nil {
		return 0, err
	}

	w.log.Debug("Allocated memory in host", "size", size, "addr", addr)

	// write the string to the memory
	m, err := w.instance.Exports.GetMemory("memory")
	if err != nil {
		panic(err)
	}

	for i, c := range s {
		m.Data()[int(addr)+i] = byte(c)
	}

	// add the null terminating character
	m.Data()[int(addr)+size] = '\x00'

	return addr, nil
}

//	m, err := w.instance.Exports.GetMemory("memory")
//	if err != nil {
//		panic(err)
//	}
//
//	// copy the string into memory
//	//f, err := w.instance.Exports.GetFunction("allocate")
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	//// allocate the memory
//	//inPtr, err := f(len(s))
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	//w.instance.Exports.GetMemory()
//
//	return 0
//}

// getStringFromMemory returns a the string stored at the Wasm modules
// memory address addr
func (w *WASM) getStringFromMemory(addr int32) (string, error) {
	m, err := w.instance.Exports.GetMemory("memory")
	if err != nil {
		panic(err)
	}

	//get the size of the string
	ss, err := w.instanceFunctions.getStringSize(addr)
	if err != nil {
		return "", err
	}

	s := string(m.Data()[addr : addr+ss.(int32)])
	w.log.Debug("Got string from memory", "addr", addr, "result", s)

	return s, nil
}
