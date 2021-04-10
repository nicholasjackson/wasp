package main

// #include <stdio.h>
// #include <string.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type WasmString uintptr
type WasmBytes uintptr

func (w *WasmBytes) Copy(data []byte) {
	size := int32(len(data)) + 1 // add one byte for the null terminator
	*w = WasmBytes(allocate(size))
	ptr := unsafe.Pointer(*w)

	buf := (*[1 << 28]byte)(ptr)[: len(data)+1 : len(data)+1]

	// copy the data
	copy(buf[1:len(data)+1], data)
	buf[0] = byte(len(data))
}

func (w *WasmBytes) GetBytes() []byte {
	blen := (*[1 << 28]byte)(unsafe.Pointer(*w))[0]
	buf := make([]byte, blen)

	copy(buf, (*[1 << 28]byte)(unsafe.Pointer(*w))[1:blen+1])

	return buf
}

// Create a C string from the Go string and return a pointer to the
// memory location where this is stored in the modules linear memory
// C.CString is not available in tiny go
func (w *WasmString) Copy(s string) {
	size := int32(len(s)) + 1 // add one byte for the null terminator
	*w = WasmString(allocate(size))
	ptr := unsafe.Pointer(*w)

	buf := (*[1 << 28]byte)(ptr)[: len(s)+1 : len(s)+1]
	copy(buf, s)
	buf[len(s)] = 0
}

// Convert a C string to a Go string
// C.GoString is not available in tiny go
func (w *WasmString) String() string {
	cstr := (*C.char)(unsafe.Pointer(*w))
	slen := int(C.strlen(cstr))
	sbuf := make([]byte, slen)

	copy(sbuf, (*[1 << 28]byte)(unsafe.Pointer(*w))[:slen:slen])

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
	s := in.String()

	out := WasmString(0)
	out.Copy("Hello " + s)

	return out
}

//go:export reverse
func reverse(inRaw WasmBytes) WasmBytes {
	inData := inRaw.GetBytes()
	outData := []byte{}

	// reverse the array
	for i := len(inData) - 1; i >= 0; i-- {
		outData = append(outData, inData[i])
	}

	outRaw := WasmBytes(0)
	outRaw.Copy(outData)

	return outRaw
}
