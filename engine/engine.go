package engine

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type Wasm struct {
	log               hclog.Logger
	instance          *wasmer.Instance
	store             *wasmer.Store
	importObject      *wasmer.ImportObject
	module            *wasmer.Module
	instanceFunctions *instanceFunctions
	callbackFunctions map[string]interface{}
}

func New(log hclog.Logger) *Wasm {
	cbf := map[string]interface{}{}
	w := &Wasm{log: log, callbackFunctions: cbf}

	return w
}

func (w *Wasm) LoadPlugin(path string) error {
	engine := wasmer.NewEngine()
	w.store = wasmer.NewStore(engine)

	wasmBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Unable to load WASM module, error: %s", err)
	}

	// Compile the module
	module, err := wasmer.NewModule(w.store, wasmBytes)
	if err != nil {
		return fmt.Errorf("Unable to instantiate WASM module, error: %s", err)
	}
	w.module = module

	// Add the Wasi environment
	wasi, err := wasmer.NewWasiStateBuilder("wasi-plugins").Environment("TESTER", "NIC").MapDirectory("host", "./").Finalize()
	if err != nil {
		return err
	}

	wi, err := wasi.GenerateImportObject(w.store, w.module)
	if err != nil {
		return err
	}
	w.importObject = wi

	// Add the default imports
	w.addDefaults()
	w.addCallbacks("plugin")

	return nil
}

func (w *Wasm) GetInstance() error {
	// Create the new instance of the module
	instance, err := wasmer.NewInstance(w.module, w.importObject)
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

// setStringInMemory copies a Go string to the Wasm modules linear memory
// it first allocates the memory by calling the modules helper function
// allocate and then copies the string.
//
// Note: Strings are copied as a null terminating string to give compatibility with
// C strings.
func (w *Wasm) setStringInMemory(s string) (int32, error) {
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

// getStringFromMemory returns a the string stored at the Wasm modules
// memory address addr
func (w *Wasm) getStringFromMemory(addr int32) (string, error) {
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

// setBytesInMemory copies the byte slice to the Wasm modules
// memory and returns the address of the data
// The function first allocates memory in the destination Wasm module
// by calling the modules allocate function copying the data.
//
// Note: The array created in the destination Wasm module always has the
// length of the array stored at the first 4 bytes as a uint32
func (w *Wasm) setBytesInMemory(data []byte) (int32, error) {
	size := len(data) + 4 // allocate 4 more bytes than the byte size as the size is encoded as a uint32 at the begining of the structure
	addr, err := w.instanceFunctions.Allocate(int32(size))
	if err != nil {
		return 0, err
	}

	w.log.Debug("Allocated memory in host", "size", size, "addr", addr)

	m, err := w.instance.Exports.GetMemory("memory")
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
func (w *Wasm) getBytesFromMemory(addr int32) ([]byte, error) {
	m, err := w.instance.Exports.GetMemory("memory")
	if err != nil {
		panic(err)
	}

	//get the size of the data from the first 4 bytes
	byteLen := binary.LittleEndian.Uint32(m.Data()[addr:])

	// copy the data
	data := make([]byte, byteLen)
	copy(data, m.Data()[addr+4:uint32(addr)+4+byteLen])

	w.log.Debug("Got bytes from memory", "addr", addr, "size", byteLen, "result", data)

	return data, nil
}
