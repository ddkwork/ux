package mylog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ddkwork/golibrary/std/stream/align"
)

type (
	log struct {
		kind     kind             // type of log
		row      keyValue         // 不要使用map，因为允许重复key
		debug    bool             // display print and write log
		callBack func(row string) // for ux logView
	}
)

const (
	hexDumpIndentLen = 26
	separate         = ` │ `
)

func (l *log) textIndent(src string, isLeftAlign bool) string {
	spaceLen := hexDumpIndentLen - align.StringWidth[int](src)
	if src == "" {
		// spaceLen -= separateLen
	}
	spaceStr := ``
	if spaceLen > 0 {
		spaceStr = strings.Repeat(" ", spaceLen)
	}
	if isLeftAlign {
		return src + spaceStr
	}
	return spaceStr + src
}

type keyValue struct {
	key   string
	value string // 重点是处理换行逻辑
}

func (k keyValue) Value() string {
	b := strings.Builder{}
	for s := range strings.Lines(k.value) {
		if s == "" || s == "\n" {
			continue
		}
		b.WriteString(s)
	}
	k.value = b.String()
	return k.value
}

func (k keyValue) String() string {
	return k.key + k.Value() // 为了兼容这些布局，分隔符需要手动处理 layout struct { //buffer 结构体 堆栈信息，普通kv行
}

func (l *log) SetCallBack(callBack func(row string)) {
	l.callBack = callBack
}

func SetCallBack(callBack func(row string)) {
	l.SetCallBack(callBack)
}

func logPath() (path string) {
	return filepath.Join(dataDir(), "log.log")
}

func dataDir() string {
	switch {
	// case IsAndroid():
	// return Check2(app.dataDir())
	case IsTermux():
		return "/data/data/com.termux/files/usr" // todo choose another dir
	default: // windows,linux
		return "."
	}
}

func New() *log {
	return &log{
		kind:     0,
		row:      keyValue{},
		debug:    true,
		callBack: nil,
	}
}

var l = New()

func init() {
	if IsAndroid() || IsTermux() {
		return
	}
	if FileLineCountIsMoreThan(logPath(), 2) {
		CheckIgnore(os.Truncate(logPath(), io.SeekStart))
	}
	Trace("--------- key ---------", "------------------ value ------------------") // android not work,why?
	if IsWindows() {
		/*
				cmd/link: enable ASLR by default on Windows
			ASLR（Address Space Layout Randomization）是一种安全技术，用于防止缓冲区溢出攻击。它通过随机化内存中关键数据结构的位置，
			使得攻击者难以预测和利用这些位置来进行攻击。在Windows操作系统中，启用ASLR可以提高系统的安全性。

					// loadDll堆栈溢出
					// go build -buildmode=exe
					// go env -w GOFLAGS="-buildmode=exe"
					// https://github.com/golang/go/issues/42593
		*/
		Check2(exec.Command("go", "env", "-w", "GOFLAGS=-buildmode=exe").CombinedOutput())
	}
	Check2(exec.Command("go", "env", "-w", "GOPROXY=https://goproxy.cn").CombinedOutput())
	// GitProxy(true)
}

var (
	GithubWorkspace = os.Getenv("GITHUB_WORKSPACE")
	IsAction        = GithubWorkspace != ""
)

func ChdirToGithubWorkspace() {
	if IsAction {
		Check(os.Chdir(GithubWorkspace))
	}
	Info("GITHUB_WORKSPACE", Check2(os.Getwd()))
}

func HexDump[K keyType, V []byte | *bytes.Buffer](title K, buf V) {
	l.hexDump(fmt.Sprint(title), DumpHex(buf))
}

func Todo(body any) {
	Warning("TODO", body)
}

type keyType interface{ string | Integer }

func Hex[K keyType, V Unsigned](title K, v V) string {
	return l.hex(fmt.Sprint(title), FormatInteger(v))
}
func Info[K keyType](title K, msg ...any)             { l.Info(fmt.Sprint(title), msg...) }
func Trace[K keyType](title K, msg ...any)            { l.Trace(fmt.Sprint(title), msg...) }
func Warning[K keyType](title K, msg ...any)          { l.Warning(fmt.Sprint(title), msg...) }
func MarshalJson[K keyType](title K, msg any)         { l.MarshalJson(fmt.Sprint(title), msg) }
func Json[K keyType](title K, msg ...any)             { l.Json(fmt.Sprint(title), msg...) }
func Success[K keyType](title K, msg ...any)          { l.Success(fmt.Sprint(title), msg...) }
func Struct(object any)                               { l.Struct(reflect.TypeOf(object).String(), object) } // log any type or struct
func SetDebug(debug bool)                             { l.debug = debug }
func Request(Request *http.Request, body bool)        { l.Request(Request, body) }
func Response(Response *http.Response, body bool)     { l.Response(Response, body) }
func DumpRequest(req *http.Request, body bool) string { return l.DumpRequest(req, body) }
func DumpResponse(resp *http.Response, body bool) string {
	return l.DumpResponse(resp, body)
}
