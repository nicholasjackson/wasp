# Web Assembly System Plugins (Wasp)

 `Wasp` is a plugin system that leverages Web Assembly (Wasm) modules for Go. Wasp allows you to extend your applictaions by allowing dynamically loaded plugins that can be authored in any language that can compile to the Wasm format. Due to the limitations and sandboxed nature of Wasm not every capability of a language, as it was originally designed to run in the browser. For example it does not natively support the ability to make network connections via sockets, read and write files, etc. Support for these features is currently being worked on as part of the Wasi standard [https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/](https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface), however until this standard is widely adopted capability across languages may vary.

To load and execute WebAssembly modules, Wasp uses the OSS [Wasmer](https://wasmer.io/) library.
 
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

To work round this limitation pointers can be used as shown in the rewritten example below. `WasmString` is not actually a Go `string` but an alias for a pointer to the 
location of a C string that is converted to a Go string with the helper function `gostring`, and exported using the `cstring` function.

```go
//go:export hello
func hello(in WasmString) WasmString {
	// get the string from the memory pointer
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

## Benchmarks:

Calling functions in Wasm modules will never be as fast as native Go functions as the Wasm function is running in a virtual environment. However the intention of Wasp is that it does not replace every function in your application but allows extension points. The following benchmarks only show a simple string calculation where most of the performance is lost through executing the plugin not the speed of the code executing in the plugin. For example, if this function was called in the context of a HTTP handler that makes a database query, adding 3000 nano seconds to a call that original took 200 milliseconds would only add 0.003 milliseconds to the total response. Wasm will always be slower than native code execution however depending on the context this may be an irrelivant and all benchmarks should be taken with a pinch of salt.

```shell
➜ go test -bench=. ./...
?       github.com/nicholasjackson/go-wasm-plugins      [no test files]
goos: linux
goarch: amd64
pkg: github.com/nicholasjackson/go-wasm-plugins/engine
cpu: AMD Ryzen 9 3950X 16-Core Processor            
BenchmarkSumGoWASM-32                     402475              3003 ns/op
BenchmarkSumRustWASM-32                   336662              3048 ns/op
BenchmarkSumTypeScriptWASM-32             363423              3068 ns/op
BenchmarkSumNative-32                   1000000000               0.2330 ns/op
PASS
ok      github.com/nicholasjackson/go-wasm-plugins/engine       4.737s
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
