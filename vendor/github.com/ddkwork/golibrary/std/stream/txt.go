package stream

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream/caseconv"
)

func ParseFloat(sizeStr string) (size float64) {
	return mylog.Check2(strconv.ParseFloat(sizeStr, 64))
}

func Float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func Float64Cut(value float64, bits int) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(bits)+"f", value), 64)
}

func ParseInt(s string) int64 {
	return mylog.Check2(strconv.ParseInt(s, 10, 64))
}

func ParseUint(s string) uint64 {
	return mylog.Check2(strconv.ParseUint(s, 10, 64))
}

func Atoi(s string) int {
	return mylog.Check2(strconv.Atoi(s))
}

func ToCamel(data string) string {
	return strings.TrimSpace(fmt.Sprintf("%-50s", caseconv.ToCamel(data)))
}

func ToCamelUpper(s string) string {
	camel := ToCamel(s)
	return strings.ToUpper(string(camel[0])) + camel[1:]
}

func ToCamelToLower(s string) string {
	camel := ToCamel(s)
	return strings.ToLower(string(camel[0])) + camel[1:]
}

// FirstToUpper converts the first character to upper case.
func FirstToUpper(in string) string {
	if in == "" {
		return in
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		return in
	}
	return string(unicode.ToUpper(r)) + in[size:]
}

// FirstToLower converts the first character to lower case.
func FirstToLower(in string) string {
	if in == "" {
		return in
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		return in
	}
	return string(unicode.ToLower(r)) + in[size:]
}

// CapitalizeWords capitalizes the first letter of each word in a string.
func CapitalizeWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = FirstToUpper(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}

// IsTruthy returns true for "truthy" values, i.e. ones that should be interpreted as true.
func IsTruthy(in string) bool {
	in = strings.ToLower(in)
	return in == "1" || in == "true" || in == "yes" || in == "on"
}

// StripBOM removes the BOM marker from UTF-8 data, if present.
func StripBOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return b[3:]
	}
	return b
}

// UnquoteBytes strips up to one set of surrounding double quotes from the bytes and returns them as a string. For a
// more capable version that supports different quoting types and unescaping, consider using strconv.Unquote().
func UnquoteBytes(text []byte) []byte {
	if len(text) > 1 {
		if ch, _ := utf8.DecodeRune(text); ch == '"' {
			if ch, _ = utf8.DecodeLastRune(text); ch == '"' {
				text = text[1 : len(text)-1]
			}
		}
	}
	return text
}

// Unquote strips up to one set of surrounding double quotes from the bytes and returns them as a string. For a more
// capable version that supports different quoting types and unescaping, consider using strconv.Unquote().
func Unquote_(text string) string {
	if len(text) > 1 {
		if ch, _ := utf8.DecodeRuneInString(text); ch == '"' {
			if ch, _ = utf8.DecodeLastRuneInString(text); ch == '"' {
				text = text[1 : len(text)-1]
			}
		}
	}
	return text
}

func Unquote(line string) string {
	begin := strings.Index(line, `\"`)
	if begin < 0 {
		return line
	}
	ss := NewBuffer("")
	for s := range strings.SplitSeq(line, `\"`) {
		if strings.Contains(s, `"`) {
			s = strings.ReplaceAll(s, `"`, ``)
			ss.WriteString(s)
		} else {
			ss.WriteString("strconv.Quote(")
			ss.WriteString(strconv.Quote(s))
			ss.WriteString(")")
		}
	}
	return ss.String()
}

// Wrap text to a certain length, giving it an optional prefix on each line. Words will not be broken, even if they
// exceed the maximum column size and instead will extend past the desired length.
func Wrap(prefix, text string, maxColumns int) string {
	var buffer strings.Builder
	for i, line := range strings.Split(text, "\n") {
		if i != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(prefix)
		avail := maxColumns - len(prefix)
		for j, token := range strings.Fields(line) {
			if j != 0 {
				if 1+len(token) > avail {
					buffer.WriteByte('\n')
					buffer.WriteString(prefix)
					avail = maxColumns - len(prefix)
				} else {
					buffer.WriteByte(' ')
					avail--
				}
			}
			buffer.WriteString(token)
			avail -= len(token)
		}
	}
	return buffer.String()
}
