package exrpc

import (
	"fmt"
	"runtime"
)

func Tracks() []string {
	var list []string
	var i int
	for {
		if i >= 20 {
			break
		}
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		i++
		list = append(list, fmt.Sprintf("%s:%d", file, line))
	}
	return list
}
