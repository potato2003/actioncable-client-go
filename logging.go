package actioncable

type Logger interface {
	Debug(string)
	Debugf(string, ...interface{})

	Info(string)
	Infof(string, ...interface{})

	Warn(string)
	Warnf(string, ...interface{})

	Error(string)
	Errorf(string, ...interface{})
}

var logger Logger = defaultLogger()

func defaultLogger() Logger {
	return &NilLogger{}
}

func SetLogger(l Logger) {
	logger = l
}
