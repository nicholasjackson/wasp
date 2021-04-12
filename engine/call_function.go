package engine

import (
	"fmt"
	"time"
)

// CallFunction in the Wasm module with the given parameters
// The response from the function will automatically be cast into the type specified
// by outputParam. In the instance that outputParam is a complex type that is returned
// as a pointer from the WASMFunction CallFunction reads the WasmModule memory and
// sets outputParam
func (w *Wasm) CallFunction(name string, outputParam interface{}, inputParams ...interface{}) error {
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

		case []byte:
			// we have a byte slice parameter, copy this into the Wasm modules
			// memory and replace with the address for the copied structure
			addr, err := w.setBytesInMemory(p.([]byte))
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
		w.log.Error("Calling function failed", "name", name, "error", err)
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

	case *[]byte:
		data, err := w.getBytesFromMemory(resp.(int32))
		if err != nil {
			return err
		}

		*outputParam.(*[]byte) = data

	case *int32:
		*outputParam.(*int32) = resp.(int32)
	default:
		return fmt.Errorf("output parameters can only be of type *int32 or *string")
	}

	return nil
}
