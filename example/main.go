package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/wasp/engine"
	"github.com/nicholasjackson/wasp/engine/logger"
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

	wl := wrappedHCLogger(log.Named("engine"))

	e := engine.New(wl)

	// add a function that can be called by wasm
	e.AddCallback("env", "call_me", callMe)

	err := e.RegisterPlugin("test", *plugin, nil)
	if err != nil {
		log.Error("Error loading plugin", "error", err)
		os.Exit(1)
	}

	// Get a new instance of the plugin
	i, err := e.GetInstance("test", "")
	if err != nil {
		log.Error("Error getting plugin instance", "error", err)
		os.Exit(1)
	}

	// cleanup
	defer i.Remove()

	// test calling an int
	var outInt int32
	err = i.CallFunction("sum", &outInt, 3, 2)
	if err != nil {
		log.Error("Error calling function", "name", "sum", "error", err)
		os.Exit(1)
	}
	log.Info("Response from function", "name", "sum", "result", outInt)

	// test calling a string
	var outStr string
	err = i.CallFunction("hello", &outStr, "Nic")
	if err != nil {
		log.Error("Error calling function", "name", "hello", "error", err)
		os.Exit(1)
	}
	log.Info("Response from function", "name", "hello", "result", outStr)

	// test calling bytes
	var outData []byte
	err = i.CallFunction("reverse", &outData, []byte{1, 2, 3})
	if err != nil {
		log.Error("Error calling function", "name", "reverse", "error", err)
		os.Exit(1)
	}
	log.Info("Response from function", "name", "reverse", "result", outData)

	err = i.CallFunction("callback", &outStr)
	if err != nil {
		log.Error("Error calling function", "name", "callback", "error", err)
		os.Exit(1)
	}
	log.Info("Response from function", "name", "callback", "result", outStr)
}

func callMe(in string) string {
	out := fmt.Sprintf("Hello %s", in)
	fmt.Println(out)

	return out
}

func wrappedHCLogger(hl hclog.Logger) *logger.Wrapper {
	return logger.New(hl.Info, hl.Debug, hl.Error, hl.Trace)
}
