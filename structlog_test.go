package structlog

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func ExampleStdLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	logger := StdLogger(1)
	logger.Log("msg", "the message")
	logger.Log("msg", "the message", "p1", 1, "lvl", "error")
	logger.Log("msg", "the message", "p1", 1, "lvl", "warn", "p2", "param 2")
	logger.Log("msg", "the message", "p1", os.Stderr, "p2", 2)
	logger.Log("msg", stringer{})
	logger.Log("msg", 123)

	// Output:
	// structlog_test.go:15: the message
	// structlog_test.go:16: error: the message: p1=1
	// structlog_test.go:17: warn: the message: p1=1 p2="param 2"
	// structlog_test.go:18: the message: p1="unsupported value type" p2=2
	// structlog_test.go:19: I'M A STRING
	// structlog_test.go:20: 123
}

func TestLoggerFunc(t *testing.T) {
	var s string
	f := func(v ...interface{}) error {
		s = fmt.Sprint(v...)
		return nil
	}

	logger := LoggerFunc(f)
	logger.Log("a", 1)
	var want = "a1"
	if got := s; got != want {
		t.Errorf("got=%q, want=%q", got, want)
	}
}

type stringer struct{}

func (s stringer) String() string {
	return "I'M A STRING"
}
