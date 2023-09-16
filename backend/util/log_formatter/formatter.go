package log_formatter

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

var truncateLength = 1000

func SetTruncateLength(n int) {
	truncateLength = n

}

func TruncateLogStr(s string) string {
	if len(s) <= truncateLength {
		return s
	}
	str := s[:truncateLength]
	return str
}

// LogPanic print panic message and transmit panic
func LogPanic(msg string) {
	if err := recover(); err != nil {
		time.Sleep(3 * time.Second)
		logrus.WithFields(logrus.Fields{
			"panic-msg":   err,
			"stack-trace": string(debug.Stack()),
		}).Panicln(msg)
	}
}

// BlankPanic just catch the panic, and stop transmitting it
func BlankPanic() {
	if recover() != nil {
	}
}
