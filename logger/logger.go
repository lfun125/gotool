package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	With(args ...interface{}) Interface
	Kind(v string) Interface
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
}

type Logger struct {
	Dir         string
	TimeFormat  string
	args        map[interface{}]interface{}
	loggerStore *sync.Map
}

func NewLogger(dir, timeFormat string) *Logger {
	l := new(Logger)
	l.Dir = dir
	l.TimeFormat = timeFormat
	l.args = map[interface{}]interface{}{}
	l.loggerStore = &sync.Map{}
	return l
}

func (l *Logger) clone() *Logger {
	n := new(Logger)
	n.Dir = l.Dir
	n.TimeFormat = l.TimeFormat
	n.args = map[interface{}]interface{}{}
	for k, v := range l.args {
		n.args[k] = v
	}
	n.loggerStore = l.loggerStore
	return n
}

func (l *Logger) clearArgs() *Logger {
	l.args = map[interface{}]interface{}{}
	return l
}

func (l *Logger) With(args ...interface{}) Interface {
	n := l.clone()
	var k interface{}
	for i, v := range args {
		if i%2 == 0 {
			n.args[v] = nil
			k = v
		} else {
			n.args[k] = v
		}
	}
	return n
}

func (l *Logger) getArgs() []interface{} {
	var data []interface{}
	for i, v := range l.args {
		data = append(data, i, v)
	}
	return data
}

func (l *Logger) Kind(v string) Interface {
	return l.clone().clearArgs().With("kind", v)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(args ...interface{}) {
	l.getZapLogger().With(l.getArgs()...).Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(args ...interface{}) {
	l.getZapLogger().With(l.getArgs()...).Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(args ...interface{}) {
	l.getZapLogger().With(l.getArgs()...).Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...interface{}) {
	l.getZapLogger().With(l.getArgs()...).Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.getZapLogger().With(l.getArgs()...).Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(args ...interface{}) {
	l.getZapLogger().With(l.getArgs()...).Fatal(args...)
}

func (l *Logger) getZapLogger() *zap.SugaredLogger {
	timeFile := time.Now().Format(l.TimeFormat)
	if v, ok := l.loggerStore.Load(timeFile); ok {
		return v.(*zap.SugaredLogger)
	}
	actual, ok := l.loggerStore.LoadOrStore(timeFile, func() *zap.SugaredLogger {
		filename := fmt.Sprintf("%s/%s.log", l.Dir, timeFile)
		encodingCfg := zap.NewProductionEncoderConfig()
		encodingCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.000Z"))
		}
		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encodingCfg),
				zapcore.AddSync(l.getWriter(filename)),
				zap.NewAtomicLevelAt(zap.InfoLevel),
			),
		)
		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		return zapLogger.Sugar()
	}())
	if !ok {
		l.loggerStore.Range(func(key, value interface{}) bool {
			if key != timeFile {
				l.loggerStore.Delete(key)
			}
			return true
		})
	}
	l.loggerStore.Range(func(key, value interface{}) bool {
		return true
	})
	return actual.(*zap.SugaredLogger)
}

func (l Logger) getWriter(filename string) (w io.Writer) {
	w = os.Stderr
	dir := path.Dir(filename)
	if dir == "/dev/stderr" {
		return os.Stderr
	} else if dir == "/dev/stdout" {
		return os.Stdout
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("create log file err: %v", err)
		return
	}
	if file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
		log.Printf("create log file err: %v", err)
		return
	} else {
		w = file
	}
	return
}
