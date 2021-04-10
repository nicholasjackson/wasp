package main

// #include <stdio.h>
// #include <string.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type WasmString uintptr

// Create a C string from the Go string and return a pointer to the
// memory location where this is stored in the modules linear memory
// C.CString is not available in tiny go
func cstring(s string) WasmString {
	size := int32(len(s)) + 1 // add one byte for the null terminator
	uptr := allocate(size)
	ptr := unsafe.Pointer(uptr)

	buf := (*[1 << 28]byte)(ptr)[: len(s)+1 : len(s)+1]
	copy(buf, s)
	buf[len(s)] = 0

	return WasmString(uptr)
}

// Convert a C string to a Go string
// C.GoString is not available in tiny go
func gostring(ptr WasmString) string {
	cstr := (*C.char)(unsafe.Pointer(ptr))
	slen := int(C.strlen(cstr))
	sbuf := make([]byte, slen)
	copy(sbuf, (*[1 << 28]byte)(unsafe.Pointer(ptr))[:slen:slen])
	return string(sbuf)
}

// allocate memory that can be written to by the Wasm host
// returns a pointer to this location in the modules linear memory.
//go:export allocate
func allocate(size int32) uintptr {
	ptr := C.malloc(C.size_t(size))
	return uintptr(ptr)
}

// enables the host to determine the size of a string
//go:export get_string_size
func getStringSize(a uintptr) int {
	char := (*C.char)(unsafe.Pointer(uintptr(a)))
	return int(C.strlen(char))
}

func main() {}

//go:export sum
func sum(a, b int) int {
	//get("test")
	return a + b
}

//go:export hello
func hello(in WasmString) WasmString {
	// get the string from the memory pointer
	s := gostring(in)

	return cstring("Hello " + s)
}
