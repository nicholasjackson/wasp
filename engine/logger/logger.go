package logger

type LogFunc func(message string, params ...interface{})

type Wrapper struct {
	info  LogFunc
	debug LogFunc
	err   LogFunc
	trace LogFunc
}

func New(info, debug, err, trace LogFunc) *Wrapper {
	return &Wrapper{info, debug, err, trace}
}

func (w *Wrapper) Info(message string, params ...interface{}) {
	if w.info != nil {
		w.info(message, params)
	}
}

func (w *Wrapper) Debug(message string, params ...interface{}) {
	if w.debug != nil {
		w.debug(message, params)
	}
}

func (w *Wrapper) Error(message string, params ...interface{}) {
	if w.err != nil {
		w.err(message, params)
	}
}

func (w *Wrapper) Trace(message string, params ...interface{}) {
	if w.trace != nil {
		w.trace(message, params)
	}
}
