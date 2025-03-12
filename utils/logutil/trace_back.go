package logutil

import (
	"fmt"
	"runtime"
)

func getTraceBack() *[]string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var trace []string

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		trace = append(trace, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	return &trace
}
