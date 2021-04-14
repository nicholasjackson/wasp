package engine

import (
	"fmt"
	"io/ioutil"

	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type plugin struct {
	module *wasmer.Module
}

type Wasm struct {
	log               *logger.Wrapper
	store             *wasmer.Store
	callbackFunctions map[string]interface{}
	plugins           map[string]*plugin
}

func New(log *logger.Wrapper) *Wasm {
	cbf := map[string]interface{}{}
	w := &Wasm{log: log, callbackFunctions: cbf}

	engine := wasmer.NewEngine()
	w.store = wasmer.NewStore(engine)
	w.plugins = map[string]*plugin{}

	return w
}

func (w *Wasm) RegisterPlugin(name, path string) error {
	wasmBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Unable to load WASM module, error: %s", err)
	}

	// Compile the module
	module, err := wasmer.NewModule(w.store, wasmBytes)
	if err != nil {
		return fmt.Errorf("Unable to instantiate WASM module, error: %s", err)
	}

	p := &plugin{
		module: module,
	}

	w.plugins[name] = p

	return nil
}

func (w *Wasm) GetInstance(name string) (*Instance, error) {
	// find the plugin
	p, ok := w.plugins[name]
	if !ok {
		return nil, fmt.Errorf("Plugin %s, not found", name)
	}

	// Create the Wasi environment
	// we can specify directories,etc for each instance
	wasi, err := wasmer.NewWasiStateBuilder("wasi-plugins").Environment("TESTER", "NIC").MapDirectory("host", "./").Finalize()
	if err != nil {
		return nil, err
	}

	io, err := wasi.GenerateImportObject(w.store, p.module)
	if err != nil {
		return nil, err
	}

	inst := &Instance{}
	inst.importObject = io

	// Add the default imports
	w.addDefaults(inst)
	w.addCallbacks("plugin", inst)

	// Create the new instance of the module
	instance, err := wasmer.NewInstance(p.module, io)
	if err != nil {
		return nil, fmt.Errorf("Unable to create a new instance of the plugin, error: %s", err)
	}

	// Setup the default functions that are required for memory manipulation operations
	wi := NewInstanceFunctions(inst)
	if err != nil {
		return nil, fmt.Errorf("Unable to import default functions, ensure that the Wasm module correctly imports the base ABI, error: %s", err)
	}

	inst.instanceFunctions = wi
	inst.instance = instance
	inst.log = w.log

	return inst, nil
}
