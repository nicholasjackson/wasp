package engine

import (
	"golang.org/x/xerrors"
)

// getStringSize calls the Wasm module to discover the size of the string
// referenced by the memory address addr
func (i *wasmerInstance) getStringSize(addr int32) (int32, error) {
	//get the size of the string
	stringSize, err := i.instance.Exports.GetFunction("get_string_size")
	if err != nil {
		return 0, err
	}

	r, err := stringSize(addr)
	if err != nil {
		return 0, err
	}

	return r.(int32), nil
}

// allocate memory in the Wasm module
// returns a pointer to the location of allocated memory that can be written to
// using the instances memory collection.
func (i *wasmerInstance) allocate(size int32) (int32, error) {
	allocate, err := i.instance.Exports.GetFunction("allocate")
	if err != nil {
		return 0, xerrors.Errorf("unable to get allocate function from module, ensure the Wasm module implements the default ABI: %w", err)
	}

	r, err := allocate(size)
	if err != nil {
		return 0, xerrors.Errorf("error calling allocate size %d: %w", size, err)
	}

	return r.(int32), nil

}

func (i *wasmerInstance) deallocate(addr int32, size int32) error {
	deallocate, err := i.instance.Exports.GetFunction("deallocate")
	if err != nil {
		return xerrors.Errorf("unable to get deallocate function from module, ensure the Wasm module implements the default ABI: %w", err)
	}

	_, err = deallocate(addr, size)
	if err != nil {
		return xerrors.Errorf("error calling deallocate addr %d: %w", addr, err)
	}

	return nil
}
