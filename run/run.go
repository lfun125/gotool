package run

import (
	"fmt"
	"runtime"
	"time"

	"github.com/lfun125/gotool/logger"
)

func Tracks() []string {
	var list []string
	var i int
	for {
		if i >= 20 {
			break
		}
		_, file, line, ok := runtime.Caller(i)
		// dir, filename := filepath.Split(file)
		// file = fmt.Sprintf("%s/%s", filepath.Base(dir), filename)
		if !ok {
			break
		}
		i++
		list = append(list, fmt.Sprintf("%s:%d", file, line))
	}
	return list
}

func GO(l *logger.Logger, f func(), number int) {
	for i := 0; i < number; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					tracks := Tracks()
					l.With("track_list", tracks).Error("panic", err)
				}
			}()
			f()
		}()
	}
}

func Hold(l *logger.Logger, handler func() error, sleepDuration, errSleepDuration time.Duration) func() {
	return func() {
		for {
			if err := handler(); err != nil {
				l.With("exec hold function error").Error(err)
				if errSleepDuration != 0 {
					time.Sleep(errSleepDuration)
				}
			} else {
				if sleepDuration != 0 {
					time.Sleep(sleepDuration)
				}
			}
		}
	}
}
