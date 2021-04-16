package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/wasp/engine/logger"
	"github.com/stretchr/testify/assert"
)

func testSetupEngine(t *testing.T, module string, conf *PluginConfig) (Instance, string) {
	hl := hclog.Default()
	hl.SetLevel(hclog.Debug)

	log := logger.New(hl.Info, hl.Debug, hl.Error, hl.Trace)
	e := New(log)

	err := e.RegisterPlugin("test", module, conf)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	dir := t.TempDir()

	inst, err := e.GetInstance("test", dir)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	return inst, dir
}

func TestWritesToWorkspaceDirectory(t *testing.T) {
	t.Skip()
	i, d := testSetupEngine(t, "../test_fixtures/rust/no_imports/module.wasm", nil)
	fmt.Println(d)

	ioutil.WriteFile(filepath.Join(d, "/in.txt"), []byte("hello"), os.ModePerm)

	i.CallFunction("workspace_write", nil, "workspace")
	assert.FileExists(t, filepath.Join(d, "hello.txt"))
}
