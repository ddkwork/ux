package mylog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	logrusPackage      string
	minimumCallerDepth = 1
	callerInitOnce     sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func caller() string {
	return long()
	return short()
}

func short() string {
	frame := getCaller()
	CheckNil(frame)
	return fmt.Sprintf("%s %s:%d", filepath.Base(frame.Function), filepath.Base(frame.File), frame.Line)
}

func long() string {
	frame := getCaller()
	CheckNil(frame)
	return fmt.Sprintf("%s+0x%x %s:%d", frame.Function, frame.PC-frame.Entry, frame.File, frame.Line)
}

func getCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		runtime.Callers(0, pcs)
		for i := range maximumCallerDepth {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}
		minimumCallerDepth = knownLogrusFrames
	})
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if pkg != logrusPackage {
			return &f
		}
	}
	return nil
}
