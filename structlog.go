// Package structlog provides a simple structured logging facade.
package structlog

import (
	"bytes"
	"fmt"
	"log"

	"github.com/go-logfmt/logfmt"
	"github.com/jjeffery/kv"
)

// Logger is the fundamental interface for all log operations. Log creates a
// log event from keyvals, a variadic sequence of alternating keys and values.
// Implementations must be safe for concurrent use by multiple goroutines. In
// particular, any implementation of Logger that appends to keyvals or
// modifies any of its elements must make a copy first.
//
// This interface (and its description) has been copied from go-kit
// (https://github.com/go-kit/kit/blob/master/log/log.go). Note that
// application logging callers are not expected to check for errors,
// see https://github.com/go-kit/kit/issues/164.
type Logger interface {
	Log(keyvals ...interface{}) error
}

// LoggerFunc is an adapter to allow use of ordinary functions as Loggers. If
// f is a function with the appropriate signature, LoggerFunc(f) is a Logger
// object that calls f.
type LoggerFunc func(...interface{}) error

// Log implements Logger by calling f(keyvals...).
func (f LoggerFunc) Log(keyvals ...interface{}) error {
	return f(keyvals...)
}

// DefaultLogger provides a default logger. The default implementation
// uses the Go standard library logger. To override, set this variable
// early in the program initialization.
var DefaultLogger Logger = StdLogger(1)

// StdLogger returns a logger that logs to the standard logger.
// Calldepth is the count of the number of frames to skip when
// computing the file name and line number if Llongfile or
// Lshortfile is set; a value of 1 will print the details for
// the caller of Log.
func StdLogger(calldepth int) Logger {
	return stdLogger{calldepth: calldepth + 2}
}

type stdLogger struct {
	calldepth int
}

func (logger stdLogger) Log(keyvals ...interface{}) error {
	msg, level, keyvals := flatten(keyvals...)
	var buf bytes.Buffer

	if level != "" {
		buf.WriteString(level)
	}
	if msg != "" {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		buf.WriteString(msg)
	}
	for len(keyvals) > 0 {
		b, err := logfmt.MarshalKeyvals(keyvals...)
		if err == nil {
			if buf.Len() > 0 {
				buf.WriteString(": ")
			}
			buf.Write(b)
			break
		}
		// cannot marshal keyvals, so keep removing keyvals until it works
		keyvals = keyvals[2:]
	}
	log.Output(logger.calldepth, buf.String())
	return nil
}

func flatten(v ...interface{}) (msg string, level string, keyvals []interface{}) {
	v = kv.Flatten(v)
	keyvals = make([]interface{}, 0, len(v))
	for i := 0; i < len(v); i += 2 {
		key := v[i].(string)
		value := v[i+1]
		switch key {
		case "msg":
			msg = toString(value)
		case "level", "lvl":
			level = toString(value)
		default:
			keyvals = append(keyvals, key, value)
		}
	}
	return msg, level, keyvals
}

func toString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
