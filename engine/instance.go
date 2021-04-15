package engine

import (
	"encoding/binary"
	"time"

	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/wasmerio/wasmer-go/wasmer"
	"golang.org/x/xerrors"
)

// Instance represents an instance of a plugin
type Instance struct {
	instance          *wasmer.Instance
	importObject      *wasmer.ImportObject
	instanceFunctions *instanceFunctions
	log               *logger.Wrapper

	// Volume is the name of the instance specific volume
	Volume string

	// map of the address and size of any memory
	// allocated by this instance
	allocatedMemory map[int32]int32
}

func NewInstance() *Instance {
	// allocatedMemory collects any pointers created by passing or receiving complex
	// types from the function.
	//
	// all the pointers in this collection should be deallocated once the function call has completed to
	// avoid leaking memory in the instance
	am := map[int32]int32{}
	return &Instance{allocatedMemory: am}
}

// Remove the instance and cleanup any volumes
func (i *Instance) Remove() error {
	return nil
}

// CallFunction in the Wasm module with the given parameters
// The response from the function will automatically be cast into the type specified
// by outputParam. In the instance that outputParam is a complex type that is returned
// as a pointer from the WASMFunction CallFunction reads the WasmModule memory and
// sets outputParam
func (i *Instance) CallFunction(name string, outputParam interface{}, inputParams ...interface{}) error {
	f, err := i.instance.Exports.GetFunction(name)
	if err != nil {
		return xerrors.Errorf("unable to find the function %s in the Wasm module: %w", err)
	}

	// ensure the deallocation of memory is always gets called, pass a reference as the slice is not yet populated
	defer i.freeAllocatedMemory()

	// parse the input parameters, if we have a string we need to set that in the Wasm modules
	// memory and pass a pointer to the function instead
	processedParams := make([]interface{}, len(inputParams))
	for n, p := range inputParams {
		switch p.(type) {
		case string:
			// we have a string parameter, let's allocate the memory for this in the wasm host and copy
			// the string to it
			addr, err := i.setStringInMemory(p.(string))
			if err != nil {
				return xerrors.Errorf("unable to set string in module memory: %w", err)
			}

			processedParams[n] = addr

			i.log.Debug(
				"Setting string in instance memory",
				"string",
				p,
				"addr",
				addr,
			)

		case []byte:
			// we have a byte slice parameter, copy this into the Wasm modules
			// memory and replace with the address for the copied structure
			addr, err := i.setBytesInMemory(p.([]byte))
			if err != nil {
				return err
			}

			processedParams[n] = addr

		default:
			processedParams[n] = p
		}
	}

	t := time.Now()

	i.log.Debug(
		"Calling function",
		"name", name,
		"outputParam", outputParam,
		"inputParam", processedParams)

	resp, err := f(processedParams...)
	if err != nil {
		i.log.Error("Calling function failed", "name", name, "error", err)
		return xerrors.Errorf("unable to call function: %w", err)
	}

	i.log.Debug(
		"Called function",
		"name", name,
		"outputParam", outputParam,
		"inputParam", processedParams,
		"response", resp,
		"time taken", time.Now().Sub(t))

	switch outputParam.(type) {
	case *string:
		s, err := i.getStringFromMemory(resp.(int32))
		if err != nil {
			return xerrors.Errorf("unable to get string from instance memory: %w", err)
		}

		*outputParam.(*string) = s

	case *[]byte:
		data, err := i.getBytesFromMemory(resp.(int32))
		if err != nil {
			return err
		}

		*outputParam.(*[]byte) = data

	case *int32:
		*outputParam.(*int32) = resp.(int32)
	default:
		return xerrors.Errorf("output parameters can only be of type *int32 or *string")
	}

	return nil
}

// setStringInMemory copies a Go string to the Wasm modules linear memory
// it first allocates the memory by calling the modules helper function
// allocate and then copies the string.
//
// setStringInMemory allocates memory in the Wasm module, this memory needs to be manually freed
// by calling the Wasm modules deallocate with the ptr returned by this function.
//
// Note: Strings are copied as a null terminating string to give compatibility with
// C strings.
func (i *Instance) setStringInMemory(s string) (int32, error) {
	size := len(s) + 1 // allocate 1 more byte than the string size for the null terminator
	addr, err := i.instanceFunctions.allocate(int32(size))
	if err != nil {
		return 0, xerrors.Errorf("unable to allocate memory in wasm module: %w", err)
	}

	// add the allocated memory to the collection so that we can deallocate it later
	i.allocatedMemory[addr] = int32(size)

	i.log.Debug(
		"Allocated memory in host",
		"size", size,
		"addr", addr)

	// write the string to the memory
	m, err := i.instance.Exports.GetMemory("memory")
	if err != nil {
		return 0, xerrors.Errorf("unable to read Wasm module memory, ensure the Wasm module exports the memory named 'memory': %w", err)
	}

	// check the memory is big enough to store the string we want
	if m.DataSize() < uint(addr)+uint(size) {
		return 0, xerrors.Errorf("unable to write string to memory, memory is not large enough to contain string")
	}

	for n, c := range s {
		m.Data()[int(addr)+n] = byte(c)
	}

	// add the null terminating character
	m.Data()[int(addr)+size] = '\x00'

	return addr, nil
}

// getStringFromMemory returns a the string stored at the Wasm modules
// memory address addr
func (i *Instance) getStringFromMemory(addr int32) (string, error) {
	m, err := i.instance.Exports.GetMemory("memory")
	if err != nil {
		return "", xerrors.Errorf("unable to read Wasm module memory, ensure the Wasm module exports the memory named 'memory': %w", err)
	}

	//get the size of the string
	ss, err := i.instanceFunctions.getStringSize(addr)
	if err != nil {
		return "", xerrors.Errorf("unable to get the size for the string at address: %d, from the Wasm module: %w", addr, err)
	}

	// check the memory is big enough to read the string we want
	if len(m.Data()) < int(addr+ss) {
		return "", xerrors.Errorf("Unable to read string from memory, memory is not large enough to contain string")
	}

	// add the allocated memory to the collection so that we can deallocate it later
	i.allocatedMemory[addr] = int32(ss)

	s := string(m.Data()[addr : addr+ss])

	i.log.Debug(
		"Got string from memory",
		"addr", addr,
		"result", s)

	return s, nil
}

// setBytesInMemory copies the byte slice to the Wasm modules
// memory and returns the address of the data
// The function first allocates memory in the destination Wasm module
// by calling the modules allocate function copying the data.
//
// Note: The array created in the destination Wasm module always has the
// length of the array stored at the first 4 bytes as a uint32
func (i *Instance) setBytesInMemory(data []byte) (int32, error) {
	size := len(data) + 4 // allocate 4 more bytes than the byte size as the size is encoded as a uint32 at the begining of the structure
	addr, err := i.instanceFunctions.allocate(int32(size))
	if err != nil {
		return 0, err
	}

	// add the allocated memory to the collection so that we can deallocate it later
	i.allocatedMemory[addr] = int32(size)

	i.log.Debug(
		"Allocated memory in host",
		"size", size,
		"addr", addr)

	m, err := i.instance.Exports.GetMemory("memory")
	if err != nil {
		panic(err)
	}

	// add the length as a uint32 to the first 4 bytes
	binary.LittleEndian.PutUint32(m.Data()[int(addr):], uint32(len(data)))

	// copy the data
	copy(m.Data()[int(addr)+4:], data)

	// return the address of the new array
	return addr, nil
}

// getBytesFromMemory copies an array from the Wasm modules memory
// into a Go byte slice. The array stored in the Wasm modules memory
// must have the length of the array encoded into the first 4 bytes
// encoded as a little endian uint32.
func (i *Instance) getBytesFromMemory(addr int32) ([]byte, error) {
	m, err := i.instance.Exports.GetMemory("memory")
	if err != nil {
		panic(err)
	}

	//get the size of the data from the first 4 bytes
	byteLen := binary.LittleEndian.Uint32(m.Data()[addr:])

	// add the allocated memory to the collection so that we can deallocate it later
	i.allocatedMemory[addr] = int32(byteLen + 4)

	// copy the data
	data := make([]byte, byteLen)
	copy(data, m.Data()[addr+4:uint32(addr)+4+byteLen])

	i.log.Debug(
		"Got bytes from memory",
		"addr", addr,
		"size", byteLen,
		"result", data)

	return data, nil
}

// freeAllocatedMemory frees any memory that has been created in the instance
// for passing complex types between the host and Wasm module
func (i *Instance) freeAllocatedMemory() {
	for addr, size := range i.allocatedMemory {
		err := i.instanceFunctions.deallocate(addr, size)
		if err != nil {
			i.log.Error(
				"Unable to deallocate memory, potential memory leak",
				"addr", addr,
				"size", size,
				"error", err,
			)
		}

		i.log.Debug(
			"Deallocated module instance memory",
			"addr", addr,
			"size", size,
		)
	}

	// clear the cache
	i.allocatedMemory = map[int32]int32{}
}
