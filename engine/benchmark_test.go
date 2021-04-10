package engine

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

func setupEngine(module string) *Wasm {
	log := hclog.NewNullLogger()
	e := New(log)

	e.LoadPlugin(module)
	return e
}

func callSumFunction(e *Wasm, b *testing.B) {
	var outInt int32
	err := e.CallFunction("sum", &outInt, 3, 2)

	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

func BenchmarkSumGoWASM(b *testing.B) {
	e := setupEngine("../plugins/go/module.wasm")

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

func BenchmarkSumRustWASM(b *testing.B) {
	e := setupEngine("../plugins/rust/pkg/sum_bg.wasm")

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

func BenchmarkSumTypeScriptWASM(b *testing.B) {
	e := setupEngine("../plugins/assemblyscript/build/optimized.wasm")

	for n := 0; n < b.N; n++ {
		callSumFunction(e, b)
	}
}

//func BenchmarkSumJavaWASM(b *testing.B) {
//	e := New()
//	e.LoadPlugin("../plugins/java/target/generated/wasm/classes.wasm")
//	sum, _ := e.GetFunction("sum")
//
//	for n := 0; n < b.N; n++ {
//		sum(34, n)
//	}
//}
//
//func BenchmarkSumCWASM(b *testing.B) {
//	e := New()
//	e.LoadPlugin("../plugins/c/a.out.wasm")
//	sum, err := e.GetFunction("sum")
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	for n := 0; n < b.N; n++ {
//		sum(34, n)
//	}
//}

func BenchmarkSumNative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sumNative(34, n)
	}
}

func sumNative(a, b int) int {
	return a * b
}
