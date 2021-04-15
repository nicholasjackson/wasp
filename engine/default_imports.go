package engine

func (w *Wasm) getDefaultCallbacks(i *Instance) *Callbacks {
	cb := &Callbacks{}

	cb.AddCallback(
		"env",
		"abort",
		func(a, b, c, d int32) {
			i.log.Debug("abort called")
		},
	)

	cb.AddCallback(
		"env",
		"raise_error",
		func(err string) {
			i.log.Debug("Error raised by plugin", "error", err)
			i.setError(err)
		},
	)

	return cb
}
