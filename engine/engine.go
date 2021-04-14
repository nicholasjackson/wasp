package engine

import (
	"fmt"
	"io/ioutil"

	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/wasmerio/wasmer-go/wasmer"
	"golang.org/x/xerrors"
)

type Wasm struct {
	log               *logger.Wrapper
	store             *wasmer.Store
	callbackFunctions map[string]map[string]interface{}
	plugins           map[string]*plugin
}

func New(log *logger.Wrapper) *Wasm {
	cbf := map[string]map[string]interface{}{}
	w := &Wasm{log: log, callbackFunctions: cbf}

	engine := wasmer.NewEngine()
	w.store = wasmer.NewStore(engine)
	w.plugins = map[string]*plugin{}

	return w
}

func (w *Wasm) RegisterPlugin(name, path string, pluginConfig *PluginConfig) error {
	wasmBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("unable to load WASM module: %w", err)
	}

	// Compile the module
	module, err := wasmer.NewModule(w.store, wasmBytes)
	if err != nil {
		return xerrors.Errorf("unable to instantiate WASM module: %w", err)
	}

	p := &plugin{
		module: module,
	}

	w.plugins[name] = p

	return nil
}

func (w *Wasm) GetInstance(name, workspaceDir string) (*Instance, error) {
	// find the plugin
	p, ok := w.plugins[name]
	if !ok {
		return nil, xerrors.Errorf("plugin %s, not found, ensure all plugins are registered before use", name)
	}

	// Create the Wasi environment
	// we can specify directories,etc for each instance
	wasi, err := wasmer.NewWasiStateBuilder("wasi-plugins").Environment("TESTER", "NIC").MapDirectory("host", "./").Finalize()
	if err != nil {
		return nil, xerrors.Errorf("unable to create Wasi state: %w", err)
	}

	io, err := wasi.GenerateImportObject(w.store, p.module)
	if err != nil {
		return nil, err
	}

	inst := &Instance{}
	inst.importObject = io

	// Add the default imports
	w.addDefaults(inst)
	w.addCallbacks(inst)

	// Create the new instance of the module
	instance, err := wasmer.NewInstance(p.module, io)
	if err != nil {
		return nil, xerrors.Errorf("unable to create a new instance of the plugin: %w", err)
	}

	// Setup the default functions that are required for memory manipulation operations
	wi := NewInstanceFunctions(inst)
	if err != nil {
		return nil, fmt.Errorf("unable to import default functions, ensure that the Wasm module correctly imports the base ABI: %w", err)
	}

	inst.instanceFunctions = wi
	inst.instance = instance
	inst.log = w.log

	return inst, nil
}
