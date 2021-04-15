package main

import "C"

import (
	"fmt"

	// import the Go ABI functions and helpers
	abi "github.com/nicholasjackson/wasp/go-abi"
)

func main() {}

//go:export sum
func sum(a, b int) int {
	//fmt.Println("Hello")
	//get("test")
	return a + b
}

//go:export hello
func hello(in abi.WasmString) abi.WasmString {
	// get the string from the memory pointer
	s := in.String()

	out := abi.WasmString(0)
	out.Copy("Hello " + s)

	return out
}

//go:export reverse
func reverse(inRaw abi.WasmBytes) abi.WasmBytes {
	inData := inRaw.Bytes()
	outData := []byte{}

	// reverse the array
	for i := len(inData) - 1; i >= 0; i-- {
		outData = append(outData, inData[i])
	}

	outRaw := abi.WasmBytes(0)
	outRaw.Copy(outData)

	return outRaw
}

// Default modules can be changed with the following comment go:wasm-module plugin

//export call_me
func callMe(in abi.WasmString) abi.WasmString

//go:export callback
func callback() abi.WasmString {
	fmt.Println("Running Function")

	// get the string from the memory pointer
	name := abi.WasmString(0)
	name.Copy("Nic")

	s := callMe(name)

	return s
}
