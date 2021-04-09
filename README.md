# Wasm Plugin Interface for Golang

This project shows how a plugin system for Go can be built that leverages Wasm modules.

Plugins can be written and executed in any language that can be compile to Wasm. Current examples in the 
repo show how plugins can be written in:
* Go (TinyGo)
* C
* Rust
* Java
* AssemblyScript

Done:  
[x] Basic plugin interface that can load and execute functions in Wasm modules  
[x] Call Go functions from Wasm modules   
[x] Receive and send int32, float32, and string types to the Wasm modules  

Todo:  
[x] Receive and send slices of bytes []byte  
[x] Ability to define cusomt ABI for plugins, currently this is hard coded  
[x] Tests, lots and lots of tests  
