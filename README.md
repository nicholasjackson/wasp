# Web Assembly System Plugins (Wasp)

 `Wasp` is a plugin system that leverages Web Assembly (Wasm) modules for Go. Wasp allows you to extend your applictaions by allowing dynamically loaded plugins that can be authored in any language that can compile to the Wasm format. Due to the limitations and sandboxed nature of Wasm not every capability of a language, as it was originally designed to run in the browser. For example it does not natively support the ability to make network connections via sockets, read and write files, etc. Support for these features is currently being worked on as part of the Wasi standard [https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/](https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface), however until this standard is widely adopted capability across languages may vary.

Current examples in the repo show how plugins can be written in:
* Go (TinyGo)
* C
* Rust
* Java
* AssemblyScript

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
