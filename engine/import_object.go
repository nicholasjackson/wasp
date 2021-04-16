package engine

import (
	"github.com/stretchr/testify/mock"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type importObject interface {
	Register(string, map[string]wasmer.IntoExtern)
}

type mockImportObject struct {
	mock.Mock
}

func (i *mockImportObject) Register(namespace string, imports map[string]wasmer.IntoExtern) {
	i.Called(namespace, imports)
}
