package log

import (
	"runtime"
	"strings"
)

// FIXME: could not get caller correctly
// find the first caller to use log.Log
func findCaller(skip int) (file, fn string, line int) {
	var pc uintptr

	for i := 0; i < 10; i++ {
		pc, file, line = runtimeCaller(skip + i)
		if !strings.HasSuffix(file, "log/caller.go") {
			break
		}
	}

	if pc != 0 {
		fn = runtime.FuncForPC(pc).Name()
	}

	return
}

// runtimeCaller report the caller line and function
func runtimeCaller(skip int) (pc uintptr, file string, line int) {
	var ok bool

	if pc, file, line, ok = runtime.Caller(skip); !ok {
		return 0, "", 0
	}

	return pc, file, line
}
