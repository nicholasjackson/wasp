package engine

import (
	"fmt"
	"testing"

	"github.com/nicholasjackson/wasp/engine/logger"
)

func setupEngine(module string, b *testing.B) *Wasm {
	log := logger.New(nil, nil, nil, nil)
	e := New(log)

	cb := &Callbacks{}
	cb.AddCallback("env", "call_me", callMe)
	conf := &PluginConfig{
		Callbacks: cb,
	}

	err := e.RegisterPlugin("test", module, conf)
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	return e
}

func callIntFunction(e *Wasm, b *testing.B) {
	var outInt int32

	inst, err := e.GetInstance("test", "")
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	err = inst.CallFunction("int_func", &outInt, 3, 2)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

func callStringFunction(e *Wasm, b *testing.B) {
	var outString string

	inst, err := e.GetInstance("test", "")
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	err = inst.CallFunction("string_func", &outString, "Nic")
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
}

func BenchmarkIntFuncGoWASM(b *testing.B) {
	e := setupEngine("../_test_fixtures/go/no_imports/module.wasm", b)

	for n := 0; n < b.N; n++ {
		callIntFunction(e, b)
	}
}

func BenchmarkStringFuncGoWASM(b *testing.B) {
	e := setupEngine("../_test_fixtures/go/no_imports/module.wasm", b)

	for n := 0; n < b.N; n++ {

		callStringFunction(e, b)
	}
}

//func BenchmarkSumRustWASM(b *testing.B) {
//	e := setupEngine("../example/plugins/rust/target/wasm32-wasi/release/module.wasi.wasm", b)
//
//	for n := 0; n < b.N; n++ {
//		callSumFunction(e, b)
//	}
//}
//
//func BenchmarkSumTypeScriptWASM(b *testing.B) {
//	e := setupEngine("../example/plugins/assemblyscript/build/optimized.wasm", b)
//
//	for n := 0; n < b.N; n++ {
//		callSumFunction(e, b)
//	}
//}
//
//func BenchmarkSumCWASM(b *testing.B) {
//	e := setupEngine("../example/plugins/c/a.out.wasm", b)
//
//	for n := 0; n < b.N; n++ {
//		callSumFunction(e, b)
//	}
//}

func BenchmarkIntFuncNative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		intNative(34, n)
	}
}

func BenchmarkStringFuncNative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stringNative("Nic")
	}
}

// callMe is a callback function imported by the wasm module
func callMe(in string) string {
	out := fmt.Sprintf("Hello %s", in)

	return out
}

func intNative(a, b int) int {
	return a * b
}

func stringNative(in string) string {
	return "Hello " + in
}
