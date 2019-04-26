package actioncable

import "log"

type StandardLibLogger struct {
	logger *log.Logger
}

func NewStandardLibLogger(l *log.Logger) Logger {
	return &StandardLibLogger{logger: l}
}

func (l *StandardLibLogger) Debug(message string) {
}

func (l *StandardLibLogger) Debugf(message string, args ...interface{}) {
}

func (l *StandardLibLogger) Info(message string) {
	l.logger.Println(message)
}

func (l *StandardLibLogger) Infof(message string, args ...interface{}) {
	l.logger.Printf(message, args...)
}

func (l *StandardLibLogger) Warn(message string) {
	l.Info(message)
}

func (l *StandardLibLogger) Warnf(message string, args ...interface{}) {
	l.Infof(message, args...)
}

func (l *StandardLibLogger) Error(message string) {
	l.Info(message)
}

func (l *StandardLibLogger) Errorf(message string, args ...interface{}) {
	l.Infof(message, args...)
}
