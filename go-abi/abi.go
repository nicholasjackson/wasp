package abi

// #include <stdio.h>
// #include <string.h>
// #include <stdlib.h>
import "C"
import (
	"encoding/binary"
	"unsafe"
)

type WasmString uintptr
type WasmBytes uintptr

func (w *WasmBytes) Copy(data []byte) {
	// add 4 bytes to store the length of the data as uint32
	size := int32(len(data)) + 4

	*w = WasmBytes(allocate(size))
	ptr := unsafe.Pointer(*w)

	buf := (*[1 << 28]byte)(ptr)[:size:size]

	// copy the data
	copy(buf[4:], data)

	// store the length of the data in the first 4 bytes as a uint32
	binary.LittleEndian.PutUint32(buf, uint32(len(data)))
}

func (w *WasmBytes) Bytes() []byte {
	data := (*[1 << 28]byte)(unsafe.Pointer(*w))

	// get the length of the data from the first 4 bytes
	len := binary.LittleEndian.Uint32(data[:4])
	buf := make([]byte, len)

	copy(buf, data[4:])

	return buf
}

// Create a C string from the Go string and return a pointer to the
// memory location where this is stored in the modules linear memory
// C.CString is not available in tiny go
func (w *WasmString) Copy(s string) {
	size := int32(len(s)) + 1 // add one byte for the null terminator

	*w = WasmString(allocate(size))
	ptr := unsafe.Pointer(*w)

	buf := (*[1 << 28]byte)(ptr)[:size:size]

	copy(buf, s)

	// add the null terminator
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

// String is a helper that returns a WasmString from a string
func String(in string) WasmString {
	ws := WasmString(0)
	ws.Copy(in)

	return ws
}

/* DEFAULT ABI */

// allocate memory that can be written to by the Wasm host
// returns a pointer to this location in the modules linear memory.
//
// Note: It is the hosts responsibility to deallocate any memory that have been
// reserved for complex types passed as references to functions or complex types
// created by this module that have been returned by this module.
//
//go:export allocate
func allocate(size int32) uintptr {
	ptr := C.malloc(C.size_t(size))
	return uintptr(ptr)
}

// deallocate memory that has previously been allocated using the
// allocate function.
//
//go:export deallocate
func deallocate(ptr uintptr, size int32) {
	C.free(unsafe.Pointer(ptr))
}

// enables the host to determine the size of a string
//
//go:export get_string_size
func getStringSize(a uintptr) int {
	char := (*C.char)(unsafe.Pointer(uintptr(a)))
	return int(C.strlen(char))
}

// RaiseError can be called to pass an error back to the host
//
//export raise_error
func raiseError(in WasmString)

func Error(in string) {
	// something bad happened
	err := WasmString(0)
	err.Copy(in)

	raiseError(err)
}

// Default workspace directory if available
const DirWorkspace = "/workspace"

/* END DEFAULT ABI */
