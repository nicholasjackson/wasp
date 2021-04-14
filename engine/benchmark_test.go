package engine

import (
	"fmt"
	"testing"

	"github.com/nicholasjackson/wasp/engine/logger"
)

func setupEngine(module string, b *testing.B) *Instance {
	log := logger.New(nil, nil, nil, nil)
	e := New(log)

	e.AddCallback("env", "call_me", callMe)

	err := e.RegisterPlugin("test", module, nil)
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	inst, err := e.GetInstance("test", "")
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	return inst
}

func callSumFunction(inst *Instance, b *testing.B) {
	var outInt int32

	err := inst.CallFunction("sum", &outInt, 3, 2)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

func BenchmarkSumGoWASM(b *testing.B) {
	e := setupEngine("../example/plugins/go/module.wasm", b)

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

func BenchmarkSumRustWASM(b *testing.B) {
	e := setupEngine("../example/plugins/rust/target/wasm32-wasi/release/module.wasi.wasm", b)

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

func BenchmarkSumTypeScriptWASM(b *testing.B) {
	e := setupEngine("../example/plugins/assemblyscript/build/optimized.wasm", b)

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

func BenchmarkSumCWASM(b *testing.B) {
	e := setupEngine("../example/plugins/c/a.out.wasm", b)

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

func BenchmarkSumNative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sumNative(34, n)
	}
}

// callMe is a callback function imported by the wasm module
func callMe(in string) string {
	out := fmt.Sprintf("Hello %s", in)

	return out
}

func sumNative(a, b int) int {
	return a * b
}
