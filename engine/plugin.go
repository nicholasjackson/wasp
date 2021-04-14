package engine

import "github.com/wasmerio/wasmer-go/wasmer"

type plugin struct {
	module    *wasmer.Module
	callbacks *Callbacks
}

// PluginConfig defines configuration for the plugin environment
type PluginConfig struct {
	// Environment variables that are available to the module instance
	Environment map[string]string
	// Volumes are globally writable volumes for the module instances
	Volumes map[string]string
}
