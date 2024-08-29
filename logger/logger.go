package logger

import (
	"fmt"
	"log/slog"
)

var (
	l *slog.Logger
)

func init() {
	if l != nil {
		return
	}

	l = slog.Default()
}

func Info(s string) {
	l.Info(s)
}

func Infof(s string, v ...any) {
	l.Info(fmt.Sprintf(s, v...))
}

func Error(e error) {
	l.Error(e.Error())
}

func Errorf(s string, v ...any) {
	l.Error(fmt.Sprintf(s, v...))
}

func SError(s string) {
	l.Error(s)
}
