package actioncable

type NilLogger struct {
}

func (l *NilLogger) Debug(message string) {
}

func (l *NilLogger) Debugf(message string, args ...interface{}) {
}

func (l *NilLogger) Info(message string) {
}

func (l *NilLogger) Infof(message string, args ...interface{}) {
}

func (l *NilLogger) Warn(message string) {
}

func (l *NilLogger) Warnf(message string, args ...interface{}) {
}

func (l *NilLogger) Error(message string) {
}

func (l *NilLogger) Errorf(message string, args ...interface{}) {
}
