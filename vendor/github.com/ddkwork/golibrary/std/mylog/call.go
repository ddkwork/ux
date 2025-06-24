package mylog

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ddkwork/golibrary/std/stream/align"
)

func Call(f func()) {
	callWithHandler(f, func(err error) { l.errorCall(err) })
}

func callWithHandler(f func(), errHandler func(err error)) {
	defer recovery(errHandler)
	f()
}

type recoveryHandler func(error)

func recovery(handler recoveryHandler) {
	if recovered := recover(); recovered != nil && handler != nil {
		e, ok := recovered.(error)
		if !ok {
			e = fmt.Errorf("%+v", recovered)
		}
		defer recovery(nil)
		handler(e)
	}
}

func (l *log) errorCall(err any) bool {
	if err == nil {
		return true
	}
	*l = log{
		callBack: l.callBack,
		kind:     errorKind,
		row: keyValue{
			key:   "",
			value: reason(err),
		},
		debug: l.debug,
	}
	l.printAndWrite2()
	return false
}

func reason(err any) string {
	switch e := err.(type) {
	case error:
		return e.Error()
	case string:
		return strings.TrimSuffix(err.(string), "\n")
	}
	return ""
}

func layoutStack(k kind, value string, child bool) string {
	if child {
		leftIndent := align.StringWidth[int](GetTimeNowString()+k.String()) + hexDumpIndentLen
		return strings.Repeat(" ", leftIndent) + " │ " + value + "\n"
	}
	leftIndent := hexDumpIndentLen // - align.StringWidth[int](key)
	if strings.Contains(value, "\n") {
		b := strings.Builder{}
		b.WriteString("\n")
		for s := range strings.Lines(value) {
			if s == "" || s == "\n" {
				continue
			}
			s = strings.TrimSuffix(s, "\r\n")
			s = strings.TrimSuffix(s, "\n")
			indent := strings.Repeat(" ", leftIndent+len(GetTimeNowString()+k.String())) + " │ " + s
			b.WriteString(indent)
			b.WriteString("\n")
		}
		value = b.String()
	}
	value = strings.TrimSuffix(value, "\n")
	return GetTimeNowString() + k.String() + strings.Repeat(" ", leftIndent) + " │ " + value + "\n"
}

func (l *log) printAndWrite2() {
	lock.Lock()
	defer lock.Unlock()
	if IsTermux() {
		return
	}
	stackChildren := make([]string, 0)
	frames := runtime.CallersFrames(callStack())
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		skip := false
		for _, s := range RuntimePrefixesToFilter {
			if strings.HasPrefix(frame.Function, s) || strings.HasSuffix(frame.Function, ".func1") {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		s := fmt.Sprintf("%s+0x%x %s:%d", frame.Function, frame.PC-frame.Entry, frame.File, frame.Line)
		stackChildren = append(stackChildren, s+"\n")
	}

	b := strings.Builder{}
	if l.row.Value() == "" {
		l.row.value = `""`
	}
	b.WriteString(layoutStack(l.kind, l.row.Value(), false))
	for _, child := range stackChildren {
		b.WriteString(layoutStack(l.kind, child, true))
	}
	l.row = keyValue{
		key:   "",
		value: b.String(),
	}
	s := l.row.String()
	s += "\n"
	l.printColorBody(s)
	if l.callBack != nil {
		l.callBack(s)
	}
	if IsAndroid() {
		println("android log is not support yet")
		return
	}
	WriteAppend(logPath(), s)
}

var RuntimePrefixesToFilter = []string{
	"runtime.",
	"testing.",
	"github.com/ddkwork/golibrary/std/mylog.callWithHandler",
	"github.com/ddkwork/golibrary/std/mylog.Call",
}

func callStack() []uintptr {
	var pcs [512]uintptr
	n := runtime.Callers(6, pcs[:])
	cs := make([]uintptr, n)
	copy(cs, pcs[:n])
	return cs
}
