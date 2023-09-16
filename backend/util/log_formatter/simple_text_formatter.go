package log_formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

type SimpleTextFormatter struct {
	TimestampFormat  string
	CallerPrettyfier func(*runtime.Frame) (function string, file string)
}

// Format renders a single log entry
func (f *SimpleTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// set timestamp format
	if f.TimestampFormat == "" {
		f.TimestampFormat = "2006/01/02 15:04:05"
	}
	// set func name and file name
	var funcVal, fileVal string
	if entry.HasCaller() {
		if f.CallerPrettyfier == nil {
			var dir string
			dir, funcVal = filepath.Split(entry.Caller.Function)

			fpStr := strings.Split(entry.Caller.File, "/")
			file := getFilePath(dir, fpStr)
			fileVal = fmt.Sprintf("%s:%d", file, entry.Caller.Line)
		} else {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
	}

	// write
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString(fmt.Sprintf("[%s]",
		entry.Time.Format(f.TimestampFormat)))
	f.writeKeyVal(b, "msg", entry.Message)
	f.writeKeyVal(b, "func", funcVal)
	f.writeKeyVal(b, "file", fileVal)
	for k, v := range entry.Data {
		strVal, ok := v.(string)
		if !ok {
			strVal = fmt.Sprint(v)
		}
		f.writeKeyVal(b, k, strVal)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *SimpleTextFormatter) writeKeyVal(b *bytes.Buffer, key string, val string) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(fmt.Sprintf("%q", val))
}
