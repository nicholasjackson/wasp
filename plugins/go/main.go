package main

// #include <stdio.h>
// #include <string.h>
// #include <stdlib.h>
import "C"
import "unsafe"

// allocate memory that can be written to
//go:export allocate
func allocate(size int32) uintptr {
	ptr := C.malloc(C.size_t(size))
	return uintptr(ptr)
}

// C.CString is not available in tiny go
func cstring(s string) uintptr {
	size := int32(len(s)) + 1 // add one byte for the null terminator
	uptr := allocate(size)
	ptr := unsafe.Pointer(uptr)

	buf := (*[1 << 28]byte)(ptr)[: len(s)+1 : len(s)+1]
	copy(buf, s)
	buf[len(s)] = 0

	return uptr
}

// C.GoString is not available in tiny go
func gostring(s *C.char) string {
	slen := int(C.strlen(s))
	sbuf := make([]byte, slen)
	copy(sbuf, (*[1 << 28]byte)(unsafe.Pointer(s))[:slen:slen])
	return string(sbuf)
}

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
func hello(in uintptr) uintptr {
	// get the string from memory pointer
	s := gostring((*C.char)(unsafe.Pointer(in)))

	return cstring("Hello " + s)
}
