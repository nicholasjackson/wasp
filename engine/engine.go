package engine

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/wasmerio/wasmer-go/wasmer"
	"golang.org/x/xerrors"
)

type Wasm struct {
	log     *logger.Wrapper
	store   *wasmer.Store
	plugins map[string]*plugin
}

/*
	New creates a new instance of the engine, it takes a single parameter
	logger.Wrapper that is used by the engine to log output.

	To create an engine with logging disabled, pass nil to the New function
*/
func New(log *logger.Wrapper) *Wasm {
	if log == nil {
		// create a nil logger
		log = logger.New(nil, nil, nil, nil)
	}

	w := &Wasm{log: log}

	engine := wasmer.NewEngine()
	w.store = wasmer.NewStore(engine)
	w.plugins = map[string]*plugin{}

	return w
}

type ImportNotFoundError struct {
	Name   string
	Module string
}

func (i ImportNotFoundError) Error() string {
	return fmt.Sprintf(
		"plugin imports the function %s from the namespace %s but no callback is defined for this import",
		i.Name,
		i.Module,
	)
}

/*
	RegisterPlugin registers a plugin with the given parameters with the engine

	Parameters:
		name: The name of the plugin as it will be registered with the engine
		pluginPath: The path to the Wasm module that will be loaded
		callbacks: A collection of functions that can be imported by the Wasm module
		pluginConfig: Additional configuration for the engine such as environment variables and volumes
*/
func (w *Wasm) RegisterPlugin(name, pluginPath string, pluginConfig *PluginConfig) error {
	wasmBytes, err := ioutil.ReadFile(pluginPath)
	if err != nil {
		return xerrors.Errorf("unable to load WASM module: %w", err)
	}

	// always create a config if one does not exist
	if pluginConfig == nil {
		pluginConfig = &PluginConfig{
			Callbacks: &Callbacks{},
		}
	}

	// Compile the module
	module, err := wasmer.NewModule(w.store, wasmBytes)
	if err != nil {
		return xerrors.Errorf("unable to instantiate WASM module: %w", err)
	}

	// validate that there are callbacks for all the imported functions
	for _, i := range module.Imports() {
		// wasi functions that are provided by the system are loaded in the wasi_... namespaces
		if strings.HasPrefix(i.Module(), "wasi_") ||
			(i.Module() == "env" && i.Name() == "raise_error") ||
			(i.Module() == "env" && i.Name() == "abort") {
			// default import
		} else {
			if pluginConfig == nil {
				return ImportNotFoundError{i.Name(), i.Module()}
			}

			if pluginConfig.Callbacks == nil {
				return ImportNotFoundError{i.Name(), i.Module()}
			}

			if m, ok := pluginConfig.Callbacks.callbackFunctions[i.Module()]; ok {
				if _, ok := m[i.Name()]; !ok {
					return ImportNotFoundError{i.Name(), i.Module()}
				}
			} else {
				return ImportNotFoundError{i.Name(), i.Module()}
			}
		}
	}

	p := &plugin{
		module: module,
		config: pluginConfig,
	}

	w.plugins[name] = p

	return nil
}

/*
	GetInstance retrieves an instance of a plugin that can be used for calling functions .The instance
	returned has its own memory and resources.

	Parameters:
		name: The name of the plugin to retrieve an instance for
		workspaceDir: Optional workspace directory to mount for the plugin, workspace directories can be
									used to share filesystem data between groups of plugins.
									this directory is mounted to /workspace inside the Wasm module.
*/
func (w *Wasm) GetInstance(name, workspaceDir string) (*Instance, error) {
	// find the plugin
	p, ok := w.plugins[name]
	if !ok {
		return nil, xerrors.Errorf("plugin %s, not found, ensure all plugins are registered before use", name)
	}

	// Create the Wasi environment
	// we can specify directories,etc for each instance
	wasi := wasmer.NewWasiStateBuilder("wasi-plugins")

	// add the environment variables
	if p.config != nil && p.config.Environment != nil {
		for k, v := range p.config.Environment {
			wasi.Environment(k, v)
		}
	}

	if workspaceDir != "" {
		wasi.MapDirectory("workspace", workspaceDir)
	}
	//.Environment("TESTER", "NIC").MapDirectory("host", "./").Finalize()
	sb, err := wasi.Finalize()
	if err != nil {
		return nil, xerrors.Errorf("unable to create Wasi state: %w", err)
	}

	io, err := sb.GenerateImportObject(w.store, p.module)
	if err != nil {
		return nil, err
	}

	inst := NewInstance()
	inst.importObject = io

	// Add the default imports
	defaultCallbacks := w.getDefaultCallbacks(inst)

	// add the default callbacks to our user defined list
	p.config.Callbacks.merge(defaultCallbacks)

	p.config.Callbacks.addCallbacks(inst, w.store, w.log)

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
