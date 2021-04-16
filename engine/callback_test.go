package engine

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/stretchr/testify/require"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func setupCallbackTests(t *testing.T) (*mockInstance, *logger.Wrapper) {
	l := hclog.NewNullLogger()
	wl := logger.New(l.Info, l.Debug, l.Error, l.Trace)

	mi := &mockInstance{}

	return mi, wl
}

func TestCreateCallbackCreatesFunctionWithCorrectSignatureInt(t *testing.T) {
	i, l := setupCallbackTests(t)

	ft, ff := createCallback(i, l, "testns", "testfunc", testCallbackFuncInt)
	require.NotNil(t, ft)
	require.NotNil(t, ff)

	// should have 1 input parameter of type int32
	require.Len(t, ft.Params(), 1)
	require.Equal(t, wasmer.I32, ft.Params()[0].Kind())

	// should have one outputParam type int32
	require.Len(t, ft.Results(), 1)
	require.Equal(t, wasmer.I32, ft.Results()[0].Kind())
}

func TestCreateCallbackCreatesFunctionWithCorrectSignatureString(t *testing.T) {
	i, l := setupCallbackTests(t)

	ft, ff := createCallback(i, l, "testns", "testfunc", testCallbackFuncString)
	require.NotNil(t, ft)
	require.NotNil(t, ff)

	// should have 1 input parameter of type int32
	require.Len(t, ft.Params(), 1)
	require.Equal(t, wasmer.I32, ft.Params()[0].Kind())

	// should have one outputParam type int32
	require.Len(t, ft.Results(), 1)
	require.Equal(t, wasmer.I32, ft.Results()[0].Kind())
}

func TestCreateCallbackCreatesFunctionWithCorrectSignatureEmpty(t *testing.T) {
	i, l := setupCallbackTests(t)

	ft, ff := createCallback(i, l, "testns", "testfunc", testCallbackFuncEmpty)
	require.NotNil(t, ft)
	require.NotNil(t, ff)

	require.Len(t, ft.Params(), 0)
	require.Len(t, ft.Results(), 0)
}

func TestCallbackFunctionPassesPtrAsString(t *testing.T) {
	i, l := setupCallbackTests(t)
	_, ff := createCallback(i, l, "testns", "testfunc", testCallbackFuncEmpty)

	_, err := ff([]wasmer.Value{wasmer.NewI32(134234)})
	require.NoError(t, err)
}

func testCallbackFuncString(in string) string {
	return in
}

func testCallbackFuncInt(in int) int {
	return in
}

func testCallbackFuncEmpty() {
}
