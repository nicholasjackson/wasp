package engine

import "github.com/nicholasjackson/wasp/engine/logger"

func (w *Wasm) getDefaultCallbacks(i Instance, l *logger.Wrapper) *Callbacks {
	cb := &Callbacks{}

	cb.AddCallback(
		"env",
		"abort",
		func(a, b, c, d int32) {
			l.Debug("abort called")
		},
	)

	cb.AddCallback(
		"env",
		"raise_error",
		func(err string) {
			l.Debug("Error raised by plugin", "error", err)
			i.setError(err)
		},
	)

	return cb
}
