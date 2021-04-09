package engine

import (
	"testing"
)

func BenchmarkSumGoWASM(b *testing.B) {
	e := New()
	e.LoadPlugin("../plugins/go/module.wasm")
	sum, _ := e.GetFunction("sum")

	for n := 0; n < b.N; n++ {
		sum(34, n)
	}
}

func BenchmarkSumRustWASM(b *testing.B) {
	e := New()
	e.LoadPlugin("../plugins/rust/pkg/sum_bg.wasm")
	sum, _ := e.GetFunction("sum")

	for n := 0; n < b.N; n++ {
		sum(34, n)
	}
}

func BenchmarkSumTypeScriptWASM(b *testing.B) {
	e := New()
	e.LoadPlugin("../plugins/assemblyscript/build/optimized.wasm")
	sum, _ := e.GetFunction("sum")

	for n := 0; n < b.N; n++ {
		sum(34, n)
	}
}

func BenchmarkSumJavaWASM(b *testing.B) {
	e := New()
	e.LoadPlugin("../plugins/java/target/generated/wasm/classes.wasm")
	sum, _ := e.GetFunction("sum")

	for n := 0; n < b.N; n++ {
		sum(34, n)
	}
}

func BenchmarkSumCWASM(b *testing.B) {
	e := New()
	e.LoadPlugin("../plugins/c/a.out.wasm")
	sum, err := e.GetFunction("sum")
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		sum(34, n)
	}
}

func BenchmarkSumNative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sumNative(34, n)
	}
}

func sumNative(a, b int) int {
	return a * b
}
