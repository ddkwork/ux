package ux

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/mylog"
)

var want = packet{
	Scheme:        "tcp",
	Method:        http.MethodGet,
	Host:          "www.baidu.com",
	Path:          "/cmsocket/",
	ContentType:   "json",
	ContentLength: 20,
	Status:        http.StatusText(http.StatusOK),
	Note:          "this is a note",
	Process:       "steam.exe",
	PadTime:       4,
}

func TestStructCodec(t *testing.T) {
	cells := MarshalRow(want, func(key string, field any) (value string) {
		// return ""
		switch key {
		case "ContentLength":
			return strconv.Itoa(field.(int)) // 日了，要断言
		case "PadTime":
			return want.PadTime.String()
		default:
			return ""
		}
	})
	for _, cell := range cells {
		mylog.Trace(cell.Key, cell.Value)
	}
	get := UnmarshalRow[packet](cells, func(key, value string) (field any) {
		switch key {
		case "ContentLength":
			return mylog.Check2(strconv.ParseInt(value, 10, 64))
		case "PadTime":
			return mylog.Check2(time.ParseDuration(value))
		default:
			return nil
		}
		return nil
	})
	mylog.Struct(get)
	assert.Equal(t, want, get)
}

// mylog.Check2(strconv.Atoi(value))
type packet struct {
	Scheme        string        // 请求协议
	Method        string        // 请求方式
	Host          string        // 请求主机
	Path          string        // 请求路径
	ContentType   string        // 收发都有
	ContentLength int           // 收发都有
	Status        string        // 返回的状态码文本
	Note          string        // 注释
	Process       string        // 进程
	PadTime       time.Duration // 请求到返回耗时
}
