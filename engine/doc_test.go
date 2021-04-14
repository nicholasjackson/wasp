package engine

import (
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/wasp/engine/logger"
)

func ExampleNew() {
	log := hclog.Default()
	log = log.Named("main")

	wl := logger.New(log.Info, log.Debug, log.Error, log.Trace)

	e := New(wl)

	// Register a plugin called test
	err := e.RegisterPlugin("test", "./path_to_plugin.wasm", nil, nil)
	if err != nil {
		panic(err)
	}

	// Get a new instance of the plugin
	i, err := e.GetInstance("test", "")
	if err != nil {
		panic(err)
	}

	// Remove any temporary files used by the instance
	defer i.Remove()

	// The output parameter from a Wasm function is passed as a reference
	// the engine will automatically convert the parameter into the correct type.
	var outInt int32

	// Call the sum function in the plugin
	err = i.CallFunction("sum", &outInt, 3, 2)
	if err != nil {
		panic(err)
	}

	log.Info("Response from function", "name", "sum", "result", outInt)
}
