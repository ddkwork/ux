package stream

import (
	"bufio"
	"bytes"
	"cmp"
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"iter"
	"math/big"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	cmp2 "github.com/google/go-cmp/cmp"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream/align"
)

// 全排列
// 矩阵置换
// 拓扑排序
// N叉树
// treeGrid
type (
	Type interface {
		string | HexString | HexDumpString | ~[]byte | cmp.Ordered | ~*bytes.Buffer | *big.Int | *Buffer
	}
	Buffer struct {
		path string
		*bytes.Buffer
	}
)

func GenMask() {
	for i := 1; i <= 8*4; i++ {
		mask := (1 << i) - 1
		fmt.Printf("%d 位掩码: 0x%X (二进制: %08b)\n", i, mask, mask)
	}
	mylog.Hex("0x000001000000007E&0xFFFFFFFF", uint64(0x000001000000007E&0xFFFFFFFF))
}

func NewBuffer[T Type](data T) *Buffer {
	switch b := any(data).(type) {
	case *Buffer:
		return b
	case []byte:
		return &Buffer{Buffer: bytes.NewBuffer(b)}
	case int8:
		return &Buffer{Buffer: bytes.NewBuffer([]byte{byte(b)})}
	case int16:
		return &Buffer{Buffer: bytes.NewBuffer(binary.LittleEndian.AppendUint16(nil, uint16(b)))}
	case int32:
		return &Buffer{Buffer: bytes.NewBuffer(binary.LittleEndian.AppendUint32(nil, uint32(b)))}
	case int, int64:
		return &Buffer{Buffer: bytes.NewBuffer(binary.LittleEndian.AppendUint64(nil, uint64(reflect.ValueOf(b).Int())))}
	case uint8:
		return &Buffer{Buffer: bytes.NewBuffer([]byte{b})}
	case uint16:
		return &Buffer{Buffer: bytes.NewBuffer(binary.LittleEndian.AppendUint16(nil, b))}
	case uint32:
		return &Buffer{Buffer: bytes.NewBuffer(binary.LittleEndian.AppendUint32(nil, b))}
	case uint, uint64:
		return &Buffer{Buffer: bytes.NewBuffer(binary.LittleEndian.AppendUint64(nil, reflect.ValueOf(b).Uint()))}
	case float32:
		return &Buffer{Buffer: bytes.NewBuffer(strconv.AppendFloat(nil, float64(b), 'f', -1, 32))}
	case float64:
		return &Buffer{Buffer: bytes.NewBuffer(strconv.AppendFloat(nil, b, 'f', -1, 64))}
	case bool:
		v := 0
		if b {
			v = 1
		}
		return &Buffer{Buffer: bytes.NewBuffer([]byte{byte(v)})}
	case *bytes.Buffer:
		return &Buffer{Buffer: b}
	case string:
		if IsFilePath(b) {
			return &Buffer{path: b, Buffer: bytes.NewBuffer(mylog.Check2(os.ReadFile(b)))} // todo ntddk.h d大文件不应该读取会影响性能
		}
		return &Buffer{Buffer: bytes.NewBufferString(b)}
	case HexString:
		return NewHexString(b)
	case HexDumpString:
		return NewHexDump(b)
	case *big.Int:
		return &Buffer{Buffer: bytes.NewBuffer(b.Bytes())}
	default:
		v := reflect.ValueOf(b)
		t := v.Type().Kind()
		switch t {
		case reflect.Slice, reflect.Array:
			if v.Index(0).Type().Kind() == reflect.Uint8 {
				return &Buffer{Buffer: bytes.NewBuffer(v.Bytes())}
			}
		}
		panic(fmt.Sprintf("unknown type %T", b))
	}
}

func NewHexString(s HexString) *Buffer {
	decodeString := mylog.Check2(hex.DecodeString(string(s)))
	return NewBuffer(decodeString)
}

type HexString string

func (b *Buffer) ReWriteSelf() {
	if b.path == "" {
		panic("path is empty")
	}
	WriteTruncate(b.path, b.Bytes())
}

func (b *Buffer) ReWriteSelfGo() {
	if b.path == "" {
		panic("path is empty")
	}
	WriteGoFile(b.path, b.Bytes())
}

func (b *Buffer) HexString() HexString      { return HexString(hex.EncodeToString(b.Bytes())) }
func (b *Buffer) HexStringUpper() HexString { return HexString(fmt.Sprintf("%#X", b.Bytes())[2:]) }
func (b *Buffer) Empty() bool               { return b.Len() == 0 }

type HexDumpString string

func NewHexDump(hexdumpStr HexDumpString) (data *Buffer) {
	hexdump := string(hexdumpStr)
	defer func() {
		mylog.CheckNil(data)
		cxx := fmt.Sprintf("%#v", data.Bytes())
		cxx = strings.ReplaceAll(cxx, "[]byte", "char data[]")
		cxx += ";\n"
		// mylog.Json("gen c++ code", cxx)
		// mylog.HexDump("recovery go buffer", data.Bytes())
	}()
	hexdump = strings.TrimSuffix(hexdump, newLine)
	hexdump = strings.TrimPrefix(hexdump, newLine)
	hexString := new(bytes.Buffer)
	switch {
	case !hasAddress(hexdump) && !strings.Contains(hexdump, sep):
		for s := range strings.FieldsSeq(hexdump) {
			hexString.WriteString(s)
		}
		decodeString := mylog.Check2(hex.DecodeString(hexString.String()))
		data = NewBuffer(decodeString)
	case strings.Contains(hexdump, sep):
		for s := range strings.Lines(hexdump) {
			if s == "" {
				continue
			}
			s = s[addressLen:strings.Index(s, sep)]
			hexString.WriteString(strings.ReplaceAll(s, " ", ""))
		}
		decodeString := mylog.Check2(hex.DecodeString(hexString.String()))
		data = NewBuffer(decodeString)
	default:
		for s := range strings.Lines(hexdump) {
			if s == "" {
				continue
			}
			indexAscii := strings.Index(s, "  ")
			if indexAscii > 0 {
				s = s[:indexAscii] // skip ascii
			}
			for field := range strings.FieldsSeq(s) {
				if len(field) > 2 {
					// 08A73200 57 61 72 68 61 6D 6D 65 72 20 34 30 2C 30 30 30  Warhammer 40,000
					// 跳过地址和ascii
					continue
				}
				hexString.WriteString(field)
			}
		}
		decodeString := mylog.Check2(hex.DecodeString(hexString.String()))
		data = NewBuffer(decodeString)
	}
	return
}

const (
	address    = "00000000  "
	sep        = "|"
	newLine    = "\n"
	addressLen = len(address)
)

func hasAddress(s string) bool {
	switch {
	case len(s) < len("00000000"):
		return false
	case strings.Contains(s, address):
		return true
	}
	return s[len("00000000")] == ' '
}

func WriteGoFile[T Type](name string, data T) {
	s := NewBuffer(data)
	source, e := format.Source(s.Bytes())
	mylog.CheckIgnore(e)
	if e != nil {
		join := filepath.Join("tmp", name)
		join = filepath.ToSlash(join)

		//318:65: missing ',' in parameter list
		errBody := e.Error()
		lineNember := "1"
		if strings.Contains(errBody, ":") {
			lineNember = strings.Split(errBody, ":")[0]
		}

		mylog.Warning("bad file ", join+":"+lineNember)
		write(join, false, s.Bytes())
		mylog.Json(cmp2.Diff(s.Bytes(), source))
		return
	}
	write(name, false, source)
}

func WriteAppend[T Type](name string, data T)     { write(name, true, data) }
func WriteTruncate[T Type](name string, data T)   { write(name, false, data) }
func WriteBinaryFile[T Type](name string, data T) { write(name, false, data) }

func write[T Type](name string, isAppend bool, data T) {
	mylog.Call(func() {
		if !CreatDirectory(filepath.Dir(name)) {
			mylog.Check(fmt.Errorf("create directory failed: %s", filepath.Dir(name)))
		}
		flag := os.O_WRONLY | os.O_CREATE | os.O_SYNC
		if isAppend {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		f, e := os.OpenFile(name, flag, 0o644)
		defer func() { mylog.Check(f.Close()) }()
		// mylog.CheckIgnore(e)
		if e != nil {
			write(name, isAppend, data)
		}
		mylog.Check2(f.Write(NewBuffer(data).Bytes()))
	})
}

func CreatDirectory(dir string) bool {
	dir = filepath.ToSlash(dir)
	info, e := os.Stat(dir)
	// mylog.CheckIgnore(e)
	if e == nil {
		if info.IsDir() {
			return true
		}
		mylog.Check("path exists but is not a directory " + dir)
	}
	switch {
	case os.IsExist(e):
		return info.IsDir()
	case os.IsNotExist(e):
		mylog.Check(os.MkdirAll(dir, os.ModePerm))
		return true
	default:
		mylog.Check(e)
	}
	return false
}

func (b *Buffer) ReaderGzip() *Buffer {
	buf := ReaderGzip(b.Bytes()).Bytes()
	b.Reset()
	mylog.Check2(b.Write(buf))
	return b
}

func ReaderGzip[T Type](data T) *Buffer {
	reader := mylog.Check2(gzip.NewReader(bytes.NewReader(NewBuffer(data).Bytes())))
	b := make([]byte, 1024*2)
	n := mylog.Check2(reader.Read(b))
	return NewBuffer(b[:n])
}

func (b *Buffer) WritePackageName() {
	b.WriteStringLn("package " + GetPackageName())
}

func GetPackageName() (pkgName string) {
	defer func() { pkgName = strings.ReplaceAll(pkgName, "-", "_") }()
	return filepath.Base(mylog.Check2(filepath.Abs(".")))
}

func (b *Buffer) CutWithIndex(x, y int) []byte { return b.Bytes()[x:y] }
func (b *Buffer) NewLine() *Buffer             { b.WriteString("\n"); return b }
func (b *Buffer) QuoteWith(s string) *Buffer   { b.WriteString(s); return b }
func (b *Buffer) WriteBinary(order binary.ByteOrder, data any) {
	mylog.Check(binary.Write(b, order, data))
}

func (b *Buffer) ReadBinary(order binary.ByteOrder) (data any) {
	mylog.Check(binary.Read(b, order, data))
	return b
}

func (b *Buffer) WriteBytesLn(buf []byte) *Buffer {
	b.Write(buf)
	b.NewLine()
	return b
}

func (b *Buffer) ReplaceAll(old, new string) *Buffer {
	s := strings.ReplaceAll(b.String(), old, new)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) Replace(old, new string, n int) *Buffer {
	s := strings.Replace(b.String(), old, new, n)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) TrimSuffix(suffix string) *Buffer {
	s := strings.TrimSuffix(b.String(), suffix)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) TrimPrefix(prefix string) *Buffer {
	s := strings.TrimPrefix(b.String(), prefix)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) TrimSpace() *Buffer {
	s := strings.TrimSpace(b.String())
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) WriteStringLn(s string) *Buffer {
	b.WriteString(s)
	b.NewLine()
	return b
}
func (b *Buffer) Quote()          { b.WriteByte('"') }
func (b *Buffer) ObjectBegin()    { b.WriteByte('{') }
func (b *Buffer) ObjectEnd()      { b.WriteByte('}') }
func (b *Buffer) SliceBegin()     { b.WriteByte('[') }
func (b *Buffer) SliceEnd()       { b.WriteByte(']') }
func (b *Buffer) Indent(deep int) { b.WriteString(strings.Repeat(" ", deep)) }

func (b *Buffer) Peek(n int) []byte {
	defer func() {
		for range n {
			mylog.Check(b.UnreadByte())
		}
	}()
	return b.ReadN(n)
}

func (b *Buffer) ReadN(n int) []byte {
	buf := make([]byte, n)
	size := mylog.Check2(b.Read(buf))
	if size != n {
		mylog.Check(fmt.Errorf("read %d bytes, but expected %d", size, n))
	}
	return buf[:size]
}

func (b *Buffer) Append(others ...*Buffer) *Buffer {
	for _, other := range others {
		b.Write(other.Bytes())
	}
	return b
}

func (b *Buffer) BigNumXorWithAlign(arg1, arg2 []byte, align int) []byte {
	xor := new(big.Int).Xor(new(big.Int).SetBytes(arg1), new(big.Int).SetBytes(arg2))
	if len(xor.Bytes()) != align {
		panic("not enough bytes for align")
	}
	return xor.Bytes()
}

func (b *Buffer) InsertString(index int, s string) *Buffer {
	buf := slices.Insert(b.Bytes(), index, []byte(s)...)
	b.Reset()
	b.Write(buf)
	return b
}

func (b *Buffer) InsertBytes(index int, insert []byte) {
	buf := slices.Insert(b.Bytes(), index, insert...)
	b.Reset()
	b.Write(buf)
}

func (b *Buffer) InsertByte(index int, ch byte) {
	buf := slices.Insert(b.Bytes(), index, ch)
	b.Reset()
	b.Write(buf)
}

func isAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (b *Buffer) Join(sep string, size int) string {
	mylog.Check(!isAllLetters(b.String()))
	result := ""
	for block := range slices.Chunk(b.Bytes(), size) {
		result += string(block) + sep
	}
	return strings.TrimSuffix(result, sep)
}

func (b *Buffer) AppendByteSlice(bytesSlice ...[]byte) []byte {
	for _, slice := range bytesSlice {
		b.Write(slice)
	}
	return b.Bytes()
}

func (b *Buffer) Contains(substr string) bool { return strings.Contains(b.String(), substr) }

func (b *Buffer) ToLines() (lines iter.Seq[string]) {
	return strings.Lines(b.String())
}

func FileLineCountIsMoreThan(path string, n int) bool {
	if !FileExists(path) {
		return false
	}
	f := mylog.Check2(os.Open(path))
	defer func() { mylog.Check(f.Close()) }()
	scanner := bufio.NewScanner(f)
	// scanner.Split(bufio.ScanLines)
	// scanner.Buffer(nil, 1024*1024)
	lineNumber := 1
	for scanner.Scan() {
		lineNumber++
		if lineNumber > n {
			return true
		}
	}
	mylog.Check(scanner.Err())
	return false
}

func ReadFileToLines(path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		f := mylog.Check2(os.Open(path))
		defer func() { mylog.Check(f.Close()) }()
		scanner := bufio.NewScanner(f)
		// scanner.Split(bufio.ScanLines)
		// scanner.Buffer(nil, 1024*1024)
		lineNumber := 1
		for scanner.Scan() {
			yield(scanner.Text())
			lineNumber++
		}
		mylog.Check(scanner.Err())
	}
}

func ReadFileToChunks(path string, n int) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		reader := mylog.Check2(os.Open(path))
		defer func() { mylog.Check(reader.Close()) }()
		r := bufio.NewReader(reader)
		buffer := make([]byte, n)
		for {
			size := mylog.Check2(r.Read(buffer))
			if size > 0 {
				// Yield the slice of the buffer that contains the data read
				if !yield(buffer[:size]) {
					break
				}
			}
		}
	}
}

// Lines returns an iterator over the newline-terminated lines in the string s.
// The lines yielded by the iterator include their terminating newlines.
// If s is empty, the iterator yields no lines at all.
// If s does not end in a newline, the final yielded line will not end in a newline.
// It returns a single-use iterator with both line number and line content.
func Lines(s string) iter.Seq2[int, string] { // 小文件，buffer，可以，大文件使用ReadFileToLines来提高性能
	return func(yield func(int, string) bool) {
		lineNumber := 1
		for len(s) > 0 {
			var line string
			if i := strings.IndexByte(s, '\n'); i >= 0 {
				line, s = s[:i+1], s[i+1:]
			} else {
				line, s = s, ""
			}
			if !yield(lineNumber, line) {
				return
			}
			lineNumber++
		}
	}
}

// LinesBytes returns an iterator over the newline-terminated lines in the byte slice s.
// The lines yielded by the iterator include their terminating newlines.
// If s is empty, the iterator yields no lines at all.
// If s does not end in a newline, the final yielded line will not end in a newline.
// It returns a single-use iterator.
func LinesBytes(s []byte) iter.Seq2[int, []byte] {
	return func(yield func(int, []byte) bool) {
		lineNumber := 1
		for len(s) > 0 {
			var line []byte
			if i := bytes.IndexByte(s, '\n'); i >= 0 {
				line, s = s[:i+1], s[i+1:]
			} else {
				line, s = s, nil
			}
			if !yield(lineNumber, line[:len(line):len(line)]) {
				return
			}
			lineNumber++
		}
	}
}

func (b *Buffer) Reverse() *Buffer {
	slices.Reverse(b.Bytes())
	return b
}

func IsZero(v reflect.Value) bool {
	if v.IsValid() {
		return true
	}
	return v.IsZero()
}

func ReflectVisibleFields(object any) iter.Seq2[int, reflect.StructField] { // todo reflect field value?
	return func(yield func(int, reflect.StructField) bool) {
		fields := reflect.VisibleFields(reflect.TypeOf(object))
		for i, field := range fields {
			if field.Tag.Get("table") == "-" || field.Tag.Get("json") == "-" {
				continue // todo
			}
			if !field.IsExported() {
				mylog.Trace("field name is not exported: ", field.Name) // 用于树形表格序列化json保存到文件，没有导出则json会失败
				continue
			}
			if !yield(i, field) {
				return
			}
		}
	}
}

func SwapAdjacent[T Type](data T) *Buffer { // 硬盘序列号交换字节
	b := NewBuffer(data)
	buf := b.Bytes()
	b.Reset()
	b.Write(swapBytes(buf))
	// 这种实现用于keygen，因为后续要xor什么的运算，转义方便
	//	uint16s := make([]uint16, 0)
	//	for bytes := range slices.Chunk(b, 2) {
	//		uint16s = append(uint16s, binary.BigEndian.Uint16(bytes))
	//	}
	return b
}

// swapBytes 交换字节切片中相邻的字节
func swapBytes(data []byte) []byte {
	for i := 0; i < len(data)-1; i += 2 {
		// 使用二进制操作交换相邻字节
		data[i], data[i+1] = data[i+1], data[i]
	}
	return data
}

func AlignString(s string, length int) (ss string) {
	runes := []rune(s)
	width := align.StringWidth[int](string(runes))
	if width < length {
		repeat := strings.Repeat(" ", length-width)
		ss = string(runes) + repeat
		return ss
	}
	if length <= len(runes) {
		return string(runes[:length])
	}
	return s
}

func MarshalJSON(v any) []byte {
	indent := mylog.Check2(json.MarshalIndent(v, "", " "))
	return indent
}

func MarshalJsonToFile(v any, name string) {
	indent := mylog.Check2(json.MarshalIndent(v, "", " "))
	ext := filepath.Ext(name)
	if ext != ".json" {
		name += ".json"
	}
	WriteTruncate(name, indent)
}

func JsonIndent(b []byte) string {
	buffer := new(bytes.Buffer)
	mylog.Check(json.Indent(buffer, b, "", " "))
	return buffer.String()
}

func RandomAnySlice[T any](slice []T) T {
	seed := rand.NewSource(time.Now().UnixNano())
	return slice[rand.New(seed).Intn(len(slice))]
}

func GenA2Z() iter.Seq[string] {
	return generateLatinAlphabets()
}

func generateLatinAlphabets() iter.Seq[string] {
	return func(yield func(string) bool) {
		for i := 'A'; i <= 'Z'; i++ {
			if !yield(string(i)) {
				return
			}
		}
	}
}

func GetWindowsLogicalDrives() iter.Seq[string] {
	return func(yield func(string) bool) {
		for s := range GenA2Z() {
			s += ":\\"
			_, e := os.Stat(s)
			if e == nil {
				if !yield(s) {
					break
				}
			}
		}
	}
}

func IsTermux() bool {
	r0, e := os.Stat("/data/data/com.termux/files/usr")
	if e == nil {
		return r0.IsDir()
	}
	return false
}
func IsAndroid() bool { return runtime.GOOS == "android" }
func IsDarwin() bool  { return runtime.GOOS == "darwin" }
func IsFreebsd() bool { return runtime.GOOS == "freebsd" }
func IsIos() bool     { return runtime.GOOS == "ios" }
func IsJs() bool      { return runtime.GOOS == "js" }
func IsLinux() bool   { return runtime.GOOS == "linux" }
func IsWindows() bool { return runtime.GOOS == "windows" }

const TimeLayout = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string { return t.Format(TimeLayout) }

func GetTimeNowString() string { return time.Now().Format("2006-01-02 15:04:05 ") }

func GetTimeStamp13Bits() int64 { return time.Now().UnixNano() / 1000000 }

func GetTimeStamp() string { return strconv.FormatInt(time.Now().UnixNano()/1000000, 10) }

func GetDaysDiff(dstTime string) string {
	t := mylog.Check2(time.Parse("2006-01-02", dstTime))
	now := time.Until(t)
	// now := t.Sub(time.Now())
	days := int(now.Hours() / 24)
	years := days / 365
	months := (days % 365) / 30
	remainingDays := (days % 365) % 30
	hours := int(now.Hours()) % 24
	minutes := int(now.Minutes()) % 60
	seconds := int(now.Seconds()) % 60

	s := NewBuffer("")
	s.WriteStringLn(fmt.Sprintf("相差天数 %d 天", days))
	s.WriteStringLn(fmt.Sprintf("相差年数 %d 年", years))
	s.WriteStringLn(fmt.Sprintf("相差月数 %d 月", months))
	s.WriteStringLn(fmt.Sprintf("相差时数 %d 时", hours))
	s.WriteStringLn(fmt.Sprintf("相差分数 %d 分", minutes))
	s.WriteStringLn(fmt.Sprintf("相差秒数 %d 秒", seconds))
	s.WriteStringLn(fmt.Sprintf("相差时间 %d 年 %d 月 %d 天 %d 时 %d 分 %d 秒",
		years, months, remainingDays, hours, minutes, seconds))
	return s.String()
}

func GetUserConfigDirs() (UserConfigDirs map[string]string) {
	UserConfigDirs = make(map[string]string)
	if runtime.GOOS == "windows" {
		dir := mylog.Check2(os.UserConfigDir())
		u := mylog.Check2(user.Current())
		UserConfigDirs[u.Username] = dir
	} else if IsTermux() {
		dir := mylog.Check2(os.UserConfigDir())
		u := mylog.Check2(user.Current())
		UserConfigDirs[u.Username] = dir
	} else {
		file := mylog.Check2(os.Open("/etc/passwd"))
		defer func() { mylog.Check(file.Close()) }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				username := parts[0]
				u, e := user.Lookup(username)
				mylog.CheckIgnore(e)
				if e != nil {
					continue
				}
				dir := u.HomeDir + "/.config"
				if strings.Contains(dir, "root") || strings.Contains(dir, "home") {
					UserConfigDirs[username] = dir
				}
			}
		}
	}
	return UserConfigDirs
}

var RegexpCenter = `(.+?)`

func RegexpWebBodyBlocks(tagName string) string {
	return `<` + tagName + `[^>]*?>[\w\W]*?<\/` + tagName + `>`
}

func IntegerToIP(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

var (
	RegexpIp     = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))`)
	RegexpIpPort = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))):([0-9]+)`)
)

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func CurrentDirName(path string) (currentDirName string) {
	if path == "" {
		path = mylog.Check2(os.Getwd())
	}
	split := strings.Split(path, "\\")
	if split == nil {
		return BaseName(filepath.Dir(path))
	}
	return split[len(split)-1]
}

func CopyDir(src, dst string) {
	mylog.Check(os.CopyFS(dst, os.DirFS(src)))
}

// Copy src to dst. src may be a directory, file, or symlink.
func Copy(src, dst string) { CopyWithMask(src, dst, 0o777) }

// CopyWithMask src to dst. src may be a directory, file, or symlink.
func CopyWithMask(src, dst string, mask fs.FileMode) {
	generalCopy(src, dst, mylog.Check2(os.Lstat(src)).Mode(), mask)
}

func generalCopy(src, dst string, srcMode, mask fs.FileMode) {
	switch {
	case srcMode&os.ModeSymlink != 0:
		linkCopy(src, dst)
	case srcMode.IsDir():
		dirCopy(src, dst, srcMode, mask)
	default:
		fileCopy(src, dst, srcMode, mask)
	}
}

func fileCopy(src, dst string, srcMode, mask fs.FileMode) {
	mylog.Check(os.MkdirAll(filepath.Dir(dst), 0o755&mask))
	f := mylog.Check2(os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, (srcMode&mask)|0o200))
	s := mylog.Check2(os.Open(src))
	defer mylog.Check(s.Close())
	mylog.Check2(io.Copy(f, s))
	mylog.Check(f.Close())
}

func dirCopy(srcDir, dstDir string, srcMode, mask fs.FileMode) {
	mylog.Check(os.MkdirAll(dstDir, srcMode&mask))
	list := mylog.Check2(os.ReadDir(srcDir))
	for _, one := range list {
		name := one.Name()
		generalCopy(filepath.Join(srcDir, name), filepath.Join(dstDir, name), one.Type(), mask)
	}
}

func linkCopy(src, dst string) {
	os.Symlink(mylog.Check2(os.Readlink(src)), dst)
}

func CopyFile(path, dstPath string) {
	mylog.Check(IsFilePathEx(path))
	WriteTruncate(dstPath, NewBuffer(path).Bytes())
}

// IsDir returns true if the specified path exists and is a directory.
func IsDir(path string) bool {
	fi, e := os.Stat(path)
	return e == nil && fi.IsDir()
}

// FileExists returns true if the path points to a regular file.
func FileExists(path string) bool {
	if fi, e := os.Stat(path); e == nil {
		mode := fi.Mode()
		return !mode.IsDir() && mode.IsRegular()
	}
	return false
}

func IsFilePathEx(path string) (ok bool) {
	return isFilePath(path, true)
}

func IsFilePath(path string) bool {
	return isFilePath(path, false)
}

func isFilePath(path string, debug bool) bool {
	pattern := []string{"://", "\n", "*", "?", "<", ">", "|"}
	for _, s := range pattern {
		if strings.Contains(path, s) {
			return false
		}
	}
	stat, e := os.Stat(path)
	if e != nil {
		if debug {
			mylog.CheckIgnore(e)
		}
		return false
	}
	// strings.Has(path, "/") || strings.Has(path, "\\")
	mode := stat.Mode()
	return !mode.IsDir() && mode.IsRegular()
}

func DirDepth(dirPath string) (depth int) {
	return strings.Count(dirPath, string(filepath.Separator))
}

func IsDirRoot(path string) bool {
	return !strings.Contains(path, string(filepath.Separator))
}

func IsDirEx(path string) (ok bool) {
	return isDir(path, true)
}

func isDir(path string, debug bool) bool {
	if strings.HasPrefix(path, ".") && IsDirRoot(path) {
		return true
	}
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return fi != nil && fi.IsDir()
}

func BaseName(path string) string {
	return TrimExtension(filepath.Base(mylog.Check2(filepath.Abs(path))))
}

func TrimExtension(path string) string {
	path = strings.ReplaceAll(path, "-", "_")
	return path[:len(path)-len(filepath.Ext(path))]
}
func JoinHomeDir(path string) (join string)  { return joinHome(path, true) }
func JoinHomeFile(path string) (join string) { return joinHome(path, false) }
func joinHome(path string, isDir bool) (join string) {
	join = filepath.Join(HomeDir(), path)
	if !IsFilePath(join) {
		switch isDir {
		case true:
			mylog.Check(os.MkdirAll(join, os.ModePerm))
		default:
			f := mylog.Check2(os.Create(join))
			mylog.Check(f.Close())
		}
	}
	return
}

func HomeDir() string {
	if u, e := user.Current(); e == nil {
		return u.HomeDir
	}
	if dir, e := os.UserHomeDir(); e == nil {
		return dir
	}
	return "."
}

func RunDir() string {
	return mylog.Check2(os.Getwd())
}

func GitProxy(isSetProxy bool) {
	mylog.Call(func() {
		s := NewBuffer("")
		if isSetProxy {
			// socks5
			// .py --mode socks5 -p 7890
			// git config --global http.sslVerify false
			s.WriteStringLn(`
[http]
    proxy = http://127.0.0.1:7890
    sslVerify = false
[https]
    proxy = http://127.0.0.1:7890
`)
		}
		s.WriteStringLn(`
[user]
	name = Admin
	email = 2762713521@qq.com
`)
		if IsWindows() {
			s.WriteStringLn(`
[core]
	autocrlf = false

[safe]
	directory = *

`)
		}

		path := JoinHomeFile(".gitconfig")
		WriteTruncate(path, s.String())
	})
}
