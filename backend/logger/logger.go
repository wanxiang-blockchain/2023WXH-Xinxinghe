package logger

import (
	rotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
	"triple_star/util/log_formatter"
)

const (
	// LogRotationTime
	// The log will generate a new file to store every 24 hours.
	LogRotationTime = 24
	// LogMaxAge the log can store 48 hours, then will be cleaned.
	LogMaxAge = 48
	// LogMaxSize the log can store 2MB data for one file.
	LogMaxSize = 2 * 1024 * 1024
	// LogPath log file path and name
	LogPath = "logs/triple-star"
	LogName = "log"
)

// rotateLog divide log to different file
func rotateLog() {
	baseLogName := path.Join(LogPath, LogName)
	writer, err := rotator.New(
		baseLogName+".%Y%m%d",
		rotator.WithClock(rotator.UTC),
		rotator.WithMaxAge(LogMaxAge*time.Hour),
		rotator.WithRotationTime(LogRotationTime*time.Hour),
		rotator.WithRotationSize(LogMaxSize),
	)
	if err != nil {
		logrus.WithField("error-msg", err).
			Errorln("config local file system for logger error")
		return
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &log_formatter.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceQuote:      true,
	})
	logrus.AddHook(lfHook)
}

// logLevel
// Level: panic < fatal < error < warn < info < debug < trace
// if we set level = info , then the "debug" and "trace" will not be showed
func getLogLevel(level string) logrus.Level {
	var logLevel logrus.Level
	switch level {
	case "panic":
		logLevel = logrus.PanicLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "trace":
		logLevel = logrus.TraceLevel
	default:
		logLevel = logrus.InfoLevel
	}
	return logLevel
}

func Init(level string) {
	log_formatter.SetTruncateLength(2000)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(getLogLevel(level))
	// output file info and method info
	logrus.SetReportCaller(true)
	// log format
	logrus.SetFormatter(&log_formatter.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		FullTimestamp:   true,
	})
	rotateLog()
}
