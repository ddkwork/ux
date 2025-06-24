package assert

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Complex interface {
	~complex64 | ~complex128
}

type Ordered interface {
	Integer | Float | ~string
}

func DumpHex[T []byte | *bytes.Buffer](buf T) (dump string) {
	var b []byte
	switch v := any(buf).(type) {
	case []byte:
		b = v
	case *bytes.Buffer:
		b = v.Bytes()
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}

	if len(b) == 0 {
		return "[]bytes{}"
	}

	length := len(b)
	switch {
	case length == 0:
		return "[]bytes{}"
	case length < 16+1:
		// 结构体字段格式打印优先
		dump += fmt.Sprintf("%#v", b) // 兼容结构体字段打印样式,复制到单元测试方便，todo 输入c语法样式,目前感觉太占空间了
		dump += "\t//"
		dump += hex.EncodeToString(b) // 方便复制到rsa解密工具测试
		dump += hex.Dump(b)
		if length < 9 {
			dump = strings.ReplaceAll(dump, "                           |", "  ")
		}
		dump = strings.TrimSuffix(dump, "\n")
		return dump
	default:
		if length > 4096 { // for x64dbg big packet
			fmt.Println("big data", length)
			b = b[:4096]
		}
		dump += formatBytesAsGoCode(b)
		dump += makeMultiLineComment(strings.NewReplacer(
			"[]byte", "unsigned char b[]=",
			"}", "};\n",
		).Replace(formatBytesAsGoCode(b)))
		dump += "\n//" + hex.EncodeToString(b) // 方便复制到rsa解密工具测试
		dump += makeMultiLineComment(hex.Dump(b))
		return
	}
}

func makeMultiLineComment(data string) string {
	s := "\n"
	s += "/*"
	s += "\n"
	s += data
	s += "*/"
	return s
}

func formatBytesAsGoCode(data []byte) string {
	var buffer bytes.Buffer
	buffer.WriteString("[]byte{\n")
	// 使用 slices.Chunk 将数据分为每8个字节一组
	for chunk := range slices.Chunk(data, 8) {
		buffer.WriteString("\t") // 添加一层缩进
		for i, b := range chunk {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(fmt.Sprintf("0x%s", hex.EncodeToString([]byte{b}))) // todo 1字节需要对齐
		}
		buffer.WriteString(",\n")
	}
	buffer.WriteString("}")
	return buffer.String()
}

func FormatInteger[T Integer](data T) string {
	// return fmt.Sprintf("%d", reflect.ValueOf(data).Interface()) + "[" + FormatIntegerHex0x(data) + "]"
	return FormatIntegerHex0x(data) + " # " + fmt.Sprintf("%d", reflect.ValueOf(data).Interface())
}

func FormatIntegerHex0x[T Integer](data T) string {
	return "0x" + FormatIntegerHex(data)
}

func FormatIntegerHex[T Integer](data T) string {
	format := ""
	switch any(data).(type) {
	case int, int64, uint, uint64, uintptr:
		format = "%016X"
	case int8, uint8:
		format = "%02X"
	case int16, uint16:
		format = "%04X"
	case int32, uint32:
		format = "%08X"
	default:
		switch reflect.TypeOf(data).Kind() {
		case reflect.Int, reflect.Uint, reflect.Uint64, reflect.Uintptr, reflect.Int64:
			format = "%016X"
		case reflect.Int8, reflect.Uint8:
			format = "%02X"
		case reflect.Int16, reflect.Uint16:
			format = "%04X"
		case reflect.Int32, reflect.Uint32:
			format = "%08X"
		default:
			panic("unsupported type ---> " + reflect.TypeOf(data).Name())
		}
	}
	return fmt.Sprintf(format, data)
}

func Integer2Bool[T Integer](value T) bool {
	switch v := any(value).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		if v == 1 {
			return true
		}
	}
	return false
}

func Bool2Integer[T Integer](b bool) T {
	var zero T
	switch any(zero).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		if b {
			zero = 1
		}
	}
	return zero
}

func IsIncludeLine(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "#") {
		return false
	}
	s = strings.TrimPrefix(s, "#")
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "include")
}

func ValueIsBytesType(v reflect.Value) bool {
	return v.Type().Elem().Kind() == reflect.Uint8
}

// func isASCIILower(c byte) bool { return 'a' <= c && c <= 'z' }
// func isASCIIUpper(c byte) bool { return 'A' <= c && c <= 'Z' }
func isASCIIDigit(c byte) bool { return '0' <= c && c <= '9' }

// func IsASCIIAlpha(s string) bool {
//	for i := 0; i < len(s); i++ {
//		c := s[i] // 直接按字节获取
//		if !isASCIILower(c) && !isASCIIUpper(c) {
//			return false
//		}
//	}
//	return true
// }

func IsASCIIDigit(s string) bool {
	for i := range len(s) {
		if !isASCIIDigit(s[i]) {
			return false
		}
	}
	return len(s) > 0 // 确保字符串非空
}

// func IsAlphanumeric(s string) bool {
//	for i := 0; i < len(s); i++ {
//		c := s[i]
//		if !isASCIIDigit(c) && !isASCIILower(c) && !isASCIIUpper(c) {
//			return false
//		}
//	}
//	return len(s) > 0 // 确保字符串非空
// }

// func isOneByteInteger(n int) bool {
//	return n >= -128 && n <= 127 // 检查有符号整数
//	// return n >= 0 && n <= 255 // 可以用于无符号整数
// }
