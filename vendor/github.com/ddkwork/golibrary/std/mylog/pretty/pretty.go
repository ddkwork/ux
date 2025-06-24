package pretty

import (
	"bytes"
	"fmt"
	"io"
	"os"
	r "reflect"
	"strconv"
	"strings"
)

const (
	DEFAULT_INDENT = "  "
	DEFAULT_NIL    = "nil"
)

type Pretty struct {
	// indent string
	Indent string
	// output recipient
	Out io.Writer
	// string for nil
	NilString string
	// compact empty array and struct
	Compact bool
	// Maximum nesting level
	MaxLevel int
}

func Print(i any) {
	PrintTo(os.Stdout, i, true)
}

func Format(i any) string {
	var out bytes.Buffer
	PrintTo(&out, i, false)
	return out.String()
}

func PrintTo(out io.Writer, i any, nl bool) {
	p := &Pretty{Indent: DEFAULT_INDENT, Out: out, NilString: DEFAULT_NIL}
	if nl {
		p.Println(i)
	} else {
		p.Print(i)
	}
}

func (p *Pretty) Print(i any) {
	p.PrintValue(r.ValueOf(i), 0)
}

func (p *Pretty) Println(i any) {
	p.PrintValue(r.ValueOf(i), 0)
	io.WriteString(p.Out, "\n")
}

func (p *Pretty) checkStringer(val r.Value) {
	if !val.CanInterface() {
		return // todo test not exported todo get filed name
	}
	m, ok := val.Interface().(fmt.Stringer)
	if ok {
		io.WriteString(p.Out, " //"+strconv.Quote(m.String()))
	}
}

func (p *Pretty) PrintValue(val r.Value, level int) {
	if !val.IsValid() {
		io.WriteString(p.Out, p.NilString)
		return
	}

	cur := strings.Repeat(p.Indent, level)
	next := strings.Repeat(p.Indent, level+1)

	newLine := "\n"
	if len(p.Indent) == 0 {
		newLine = " "
	}

	if p.MaxLevel > 0 && level >= p.MaxLevel {
		io.WriteString(p.Out, val.String())
		return
	}

	switch val.Kind() {
	case r.Int, r.Int64:
		io.WriteString(p.Out, FormatInteger(val.Int()))
		p.checkStringer(val)
	case r.Int8:
		io.WriteString(p.Out, FormatInteger(int8(val.Int())))
		p.checkStringer(val)
	case r.Int16:
		io.WriteString(p.Out, FormatInteger(int16(val.Int())))
		p.checkStringer(val)
	case r.Int32:
		io.WriteString(p.Out, FormatInteger(int32(val.Int())))
		p.checkStringer(val)
	case r.Uint, r.Uint64, r.Uintptr: // Uintptr for Windows api use etc
		io.WriteString(p.Out, FormatInteger(val.Uint()))
		p.checkStringer(val)
	case r.Uint8:
		io.WriteString(p.Out, FormatInteger(byte(val.Uint())))
		p.checkStringer(val)
	case r.Uint16:
		io.WriteString(p.Out, FormatInteger(uint16(val.Uint())))
		p.checkStringer(val)
	case r.Uint32:
		io.WriteString(p.Out, FormatInteger(uint32(val.Uint())))
		p.checkStringer(val)
	case r.Float32, r.Float64:
		io.WriteString(p.Out, strconv.FormatFloat(val.Float(), 'f', -1, 64))
		p.checkStringer(val)
	case r.String:
		io.WriteString(p.Out, strconv.Quote(val.String()))
		p.checkStringer(val)
	case r.Bool:
		io.WriteString(p.Out, strconv.FormatBool(val.Bool()))
		p.checkStringer(val)
	case r.Map:
		l := val.Len()

		io.WriteString(p.Out, "{"+newLine)
		for i, k := range val.MapKeys() {
			io.WriteString(p.Out, next)
			io.WriteString(p.Out, strconv.Quote(k.String()))
			io.WriteString(p.Out, ": ")
			p.PrintValue(val.MapIndex(k), level+1)
			if i < l-1 {
				io.WriteString(p.Out, ","+newLine)
			} else {
				io.WriteString(p.Out, newLine)
			}
		}
		io.WriteString(p.Out, cur)
		io.WriteString(p.Out, "}")

	case r.Array, r.Slice:
		l := val.Len()
		if ValueIsBytesType(val) {
			if val.Kind() == r.Array {
				if !val.CanAddr() {
					panic("如果结构体字段是数组，那么字段格式化不可寻址，请使用指针类型,在外部 &object传进来,字段类型保持不变") // panic("reflect.Value.Bytes of unaddressable byte array")
				}
			}
			io.WriteString(p.Out, strings.Replace(DumpHex(val.Bytes()), "}", "},", 1))
			p.checkStringer(val)
			return
		}

		if p.Compact && l == 0 {
			io.WriteString(p.Out, "[]")
		} else {
			io.WriteString(p.Out, "["+newLine)
			for i := range l {

				io.WriteString(p.Out, next)

				io.WriteString(p.Out, "<")
				io.WriteString(p.Out, strconv.Itoa(i))
				io.WriteString(p.Out, "> ")

				p.PrintValue(val.Index(i), level+1)
				if i < l-1 {
					io.WriteString(p.Out, ","+newLine)
				} else {
					io.WriteString(p.Out, newLine)
				}
			}
			io.WriteString(p.Out, cur)
			io.WriteString(p.Out, "]")
		}

	case r.Interface, r.Ptr:
		p.PrintValue(val.Elem(), level)

	case r.Struct:
		if val.CanInterface() {
			i := val.Interface()
			if i, ok := i.(fmt.Stringer); ok {
				io.WriteString(p.Out, i.String())
			} else {
				l := val.NumField()
				sOpen := val.Type().String() + " {"
				if p.Compact {
					sOpen = "{"
				}
				if p.Compact && l == 0 {
					io.WriteString(p.Out, "{}")
				} else {
					io.WriteString(p.Out, sOpen+newLine)
					// 计算最大字段名长度，以便对齐
					maxKeyLen := 0
					for i := range l {
						keyLen := len(val.Type().Field(i).Name)
						if keyLen > maxKeyLen {
							maxKeyLen = keyLen
						}
					}
					for i := range l {
						io.WriteString(p.Out, next)
						fieldName := val.Type().Field(i).Name

						// 计算当前键名后需要的空格数
						spaces := strings.Repeat(" ", maxKeyLen-len(fieldName))
						io.WriteString(p.Out, spaces) // 填充空格实现对齐
						io.WriteString(p.Out, fieldName)
						io.WriteString(p.Out, ": ")
						p.PrintValue(val.Field(i), level+1)
						if i < l-1 {
							io.WriteString(p.Out, ","+newLine) // todo 整数格式化后已经有 , 了，需要移除
						} else {
							io.WriteString(p.Out, newLine)
						}
					}
					io.WriteString(p.Out, cur)
					io.WriteString(p.Out, "}")
				}
			}
		} else {
			io.WriteString(p.Out, "protected")
		}
	default:
		io.WriteString(p.Out, "unsupported:")
		io.WriteString(p.Out, val.String())
	}
}
