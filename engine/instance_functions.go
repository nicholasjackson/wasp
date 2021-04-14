package engine

import (
	"golang.org/x/xerrors"
)

// instanceFunctions are functions that must be exported in the destination
// Wasm module to satisfy Wasps ABI
type instanceFunctions struct {
	inst *Instance
}

func NewInstanceFunctions(inst *Instance) *instanceFunctions {
	return &instanceFunctions{inst}
}

// GetStringSize calls the Wasm module to discover the size of the string
// referenced by the memory address addr
func (i *instanceFunctions) getStringSize(addr int32) (int32, error) {
	//get the size of the string
	stringSize, err := i.inst.instance.Exports.GetFunction("get_string_size")
	if err != nil {
		return 0, err
	}

	r, err := stringSize(addr)
	if err != nil {
		return 0, err
	}

	return r.(int32), nil
}

// Allocate
func (i *instanceFunctions) allocate(size int32) (int32, error) {
	allocate, err := i.inst.instance.Exports.GetFunction("allocate")
	if err != nil {
		return 0, xerrors.Errorf("unable to get allocate function from module, ensure the Wasm module implements the default ABI: %w", err)
	}

	r, err := allocate(size)
	if err != nil {
		return 0, xerrors.Errorf("error calling allocate size %d: %w", size, err)
	}

	return r.(int32), nil

}
