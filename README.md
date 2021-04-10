# Web Assembly System Plugins (Wasp)

 `Wasp` is a plugin system that leverages Web Assembly (Wasm) modules for Go. Wasp allows you to extend your applictaions by allowing dynamically loaded plugins that can be authored in any language that can compile to the Wasm format. Due to the limitations and sandboxed nature of Wasm not every capability of a language, as it was originally designed to run in the browser. For example it does not natively support the ability to make network connections via sockets, read and write files, etc. Support for these features is currently being worked on as part of the Wasi standard [https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/](https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface), however until this standard is widely adopted capability across languages may vary.
 
In addition Wasm does not have the rich type system that Go has, currently Wasm only supports the types `int32, int64, float32, float64`. This means that any functions that are exposed by a Wasm module can only contain these types for data interchange. For example, the following Go function could be compiled to Wasm using TinyGo or the experimental Go Wasm target and could be called successfully by the interpreter.

```go
//go:export sum
func sum(a, b int) int {
	//get("test")
	return a + b
}
```

But it is not possible to pass a string in and out of the function as `string` is not defined by Wasms type structure:

```go
//go:export hello
func hello(name string) string {
	return "Hello " + s
}
```

To work round this limitation pointers can be used as shown in the rewritten example below.

```go
//go:export hello
func hello(in uintptr) uintptr {
	// get the string from memory pointer
	s := gostring(in)

	return cstring("Hello " + s)
}
```

And the same example written in AssemblyScript:

```TypeScript
export function hello(name: ArrayBuffer): ArrayBuffer {
  let inParam = String.UTF8.decode(name,true)

  return String.UTF8.encode("Hello " + inParam, true)
}
```

Wasp provides a data interchange ABI and helper functions for your Wasm modules that simplifies this process and automatically manages the process of copying Go `string` and `[]byte` to the Wasm modules memory space. Note: the limitation on parameters only affects functions that interface with the plugin host internal functions and methods can use the full type system of the language used to author the Wasm module.

Current examples in the repo show how plugins can be written in:
* Go (TinyGo)
* C
* Rust
* Java
* AssemblyScript

## Basic Usage:

The following example shows how Wasp can be used to call the method `hello` that was exported from a Wasm module. First you create an instance of the engine and load the plugin.

```go
// Create a logger
log := hclog.Default()
log = log.Named("main")

// Create the plugin engine 
e := engine.New(log.Named("engine"))
 
// Load and compile the wasm module
err := e.LoadPlugin("./plugins/go/module.wasm")
if err != nil {
	log.Error("Error loading plugin", "error", err)
	os.Exit(1)
}
```

Then you can use the `CallFunction` method on the engine to call the `hello` function exported from the Wasm module, Wasp automatically converts Go types into the simple types understood by the Wasm module. In the following example Wasp would take the input string "hello", allocate the required memory inside the Wasm module, copy the string data to this memory before calling the destination function with a pointer to this string. Responses work exactly the same way in reverse. 

```go
// Call the function hello that is exported by the module
var outString string
err = e.CallFunction("hello", &outString, 3, 2)
if err != nil {
	log.Error("Error calling function", "name", "hello", "error", err)
	os.Exit(1)
}
log.Info("Response from function", "name", "hello", "result", outString)
```

## Features:
**Done:**  
[x] Basic plugin interface that can load and execute functions in Wasm modules  
[x] Call Go functions from Wasm modules   
[x] Receive and send int32, float32, and string types to the Wasm modules  

**Todo:**  
[ ] Receive and send slices of bytes []byte  
[ ] Ability to define custom ABIs for plugins, currently this is hard coded  
[ ] Tests, lots and lots of tests  
[ ] Support Wasi standard
[ ] Define more robust helper packages for managing complex types
