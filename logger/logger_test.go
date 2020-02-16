package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger_Info(t *testing.T) {
	l := NewLogger(".", "20060102150405")
	for {
		l.With("a", 1).With("c", 2).Info("haha")
		l.Error("cc")
		time.Sleep(time.Second)
	}
}

func TestLogger_Kind(t *testing.T) {
	l := NewLogger(".", "20060102")
	l1 := l.Kind("1")
	l2 := l.Kind("2")
	fmt.Printf("%+v\n", l1)
	fmt.Printf("%+v\n", l2)
}

func TestLogW(t *testing.T) {
	l := NewLogger(".", "20060102150405")
	l.With("a", 1).Info("aaa")
}
