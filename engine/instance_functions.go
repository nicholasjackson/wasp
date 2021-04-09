package engine

type instanceFunctions struct {
	getStringSize func(...interface{}) (interface{}, error)
	allocate      func(...interface{}) (interface{}, error)
	deallocate    func(int32) (int32, error)
}

func NewInstanceFunctions(w *WASM) (*instanceFunctions, error) {
	i := &instanceFunctions{}

	//get the size of the string
	stringSize, err := w.instance.Exports.GetFunction("get_string_size")
	if err != nil {
		return nil, err
	}

	i.getStringSize = stringSize

	allocate, err := w.instance.Exports.GetFunction("allocate")
	if err != nil {
		return nil, err
	}

	i.allocate = allocate

	return i, nil
}

// GetStringSize calls the Wasm module to discover the size of the string
// referenced by the memory address addr
func (i *instanceFunctions) GetStringSize(addr int32) (int32, error) {
	r, err := i.getStringSize(addr)
	if err != nil {
		return 0, err
	}

	return r.(int32), nil
}

// Allocate
func (i *instanceFunctions) Allocate(size int32) (int32, error) {
	r, err := i.allocate(size)
	if err != nil {
		return 0, err
	}

	return r.(int32), nil

}
