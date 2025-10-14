package truenas

import "log"

type Logger interface {
	Println(msg string)
	Printf(format string, v ...any)
}

type defaultLogger struct{}

func (l *defaultLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func (l *defaultLogger) Println(msg string) {
	log.Println(msg)
}
