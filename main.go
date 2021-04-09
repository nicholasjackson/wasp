package main

import (
	"flag"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/go-wasm-plugins/engine"
)

var plugin = flag.String("plugin", "", "Path to the Wasm module to load")
var verbose = flag.Bool("v", false, "Verbose output")

func main() {
	flag.Parse()

	log := hclog.Default()
	log = log.Named("main")

	if *verbose {
		log.SetLevel(hclog.Debug)
	}

	e := engine.New(log.Named("engine"))
	err := e.LoadPlugin(*plugin)
	if err != nil {
		log.Error("Error loading plugin", "error", err)
		os.Exit(1)
	}

	// test calling an int
	var outInt int32
	err = e.CallFunction("sum", &outInt, 3, 2)
	if err != nil {
		log.Error("Error calling function", "name", "sum", "error", err)
		os.Exit(1)
	}
	log.Info("Response from function", "name", "sum", "result", outInt)

	// test calling a string
	var outStr string
	err = e.CallFunction("hello", &outStr, "Nic")
	if err != nil {
		log.Error("Error calling function", "name", "hello", "error", err)
		os.Exit(1)
	}
	log.Info("Response from function", "name", "hello", "result", outStr)
}
