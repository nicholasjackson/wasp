package main

import (
	"fmt"
	"os"
	"path"

	// import the Go ABI functions and helpers
	abi "github.com/nicholasjackson/wasp/go-abi"
)

func main() {}

//go:export string_func
func stringFunc(in abi.WasmString) abi.WasmString {
	// get the string from the memory pointer
	s := in.String()

	out := abi.WasmString(0)
	out.Copy("Hello " + s)

	return out
}

//go:export int_func
func intFunc(a, b int) int {
	//fmt.Println("Hello")
	//get("test")
	return a + b
}

//go:export workspace_write
func workspaceWrite(dir abi.WasmString) {
	filePath := dir.String()
	fmt.Println("Writing, ", filePath)

	f, err := os.Open(path.Join(filePath, "in.txt"))
	if err != nil {
		fmt.Println("Unable to write file", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 100)
	f.Read(buf)
	if err != nil {
		fmt.Println("Unable to write file", err)
		return
	}

	fout, err := os.Create(path.Join(filePath, "hello.txt."))
	if err != nil {
		fmt.Println("Unable to write file", err)
	}
	defer fout.Close()

	fout.Write([]byte("blah"))
}
