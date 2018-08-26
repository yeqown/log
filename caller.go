package log

import (
	"runtime"
	"strings"
)

func findCaller(skip int) (file, function string, line int) {
	var (
		pc uintptr
	)

	for i := 0; i < 10; i++ {
		pc, file, line = runtimeCaller(skip + i)
		// println(file, line)
		if !strings.HasPrefix(file, "log") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return
}

func runtimeCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	// println(file)
	// n := 0
	// for i := len(file) - 1; i > 0; i-- {
	// 	if file[i] == '/' {
	// 		n++
	// 		if n >= 2 {
	// 			file = file[i+1:]
	// 			break
	// 		}
	// 	}
	// }

	splited := strings.Split(file, "/")
	file = strings.Join(splited[len(splited)-2:], "/")

	return pc, file, line
}
