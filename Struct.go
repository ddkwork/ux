package ux

import (
	"fmt"
	"image/color"
	"reflect"
	"strconv"
	"time"

	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
)

// Table 是row布局，grid是rows和columns布局，网格是处理所有单元格的，row的话是通过结构体处理行单元格的
// column布局只有一个场景，复制列数据到剪切板，没有别的场景了。不需要支持

// 布局说明:
// structView:key-value形式的结构体字段渲染成一行，然后flex垂直排列
// 树形表格: key-value形式的结构体字段渲染成一行，然后flex水平排列，key用于每一列的标题，value用于每一行的每个单元格值
// 两种场景都支持自定义格式化和反序列化字段，通过回调函数实现，
// 对于序列化，回调返回空则使用默认的fmt包格式化字段值，否则使用回调返回值
// 对于反序列化，回调返回nil则使用默认的反序列化方式，否则使用回调返回值更新字段值

// filed或者cell都是key-value键值对，对于form和structView，key是字段名，value是字段值，对于树形表格，key是列名，value是单元格值
// 树形节点编辑调用structView布局
// 适用场景：树形表格，structView，form 布局

type (
	MarshalRowCallback   func(key string, field any) (value string) // 元数据结构体的字段名称是key的时候，执行自定义格式化元数据结构体的同名字段返回value给structView存储到rows []cellData和显示
	UnmarshalRowCallback func(key, value string) (field any)        // 元数据结构体的字段名称是key的时候，执行自定义反序列化已经格式化的value赋值给元数据结构体的同名字段，实现更新字段值
)

func MarshalRow[T any](data T, callback MarshalRowCallback) (rows []CellData) { // flex水平排列
	return MarshalFields[T](data, callback) // 获得字段键值对后渲染树形表格的row是flex水平排列，节点编辑是structView的flex垂直排列
}

func UnmarshalRow[T any](rows []CellData, callback UnmarshalRowCallback) T { // 树形节点编辑,flex垂直排列,布局成StructView
	return UnmarshalFields[T](rows, callback) // 从StructView的每一行的键值对反序列化成结构体,结果是更新树形表格的rowCells
}

// MarshalFields 序列化字段,想象成structView的flex垂直布局的每一行，每行都是一个键值对，key是字段名，value是字段值,callback理解成structView的每一行的自定义序列化的callback
func MarshalFields[T any](data T, callback MarshalRowCallback) (fields []CellData) { // flex垂直排列
	fields = make([]CellData, 0)
	if reflect.TypeOf(data).Kind() != reflect.Struct {
		panic("data must be a struct not a pointer")
	}
	visibleFields := reflect.VisibleFields(reflect.TypeOf(data))
	value := reflect.Indirect(reflect.ValueOf(data))
	if len(visibleFields) != value.NumField() {
		panic("wrong number of visibleFields")
	}
	for i, field := range visibleFields {
		if field.Tag.Get("table") == "-" || field.Tag.Get("json") == "-" {
			// mylog.Trace("field is ignored: ", field.Name) // 用于树形表格序列化json保存到文件，标签table或json为-则忽略
			continue // todo test
		}
		if !field.IsExported() {
			mylog.Trace("field name is not exported: ", field.Name) // 用于树形表格序列化json保存到文件，没有导出则json会失败
			continue
		}

		// field.Tag.Get("table") 1= "" // todo 编解码通过这个获取中文表头似乎不合适，具体用例的时候测试一波

		v := value.Field(i).Interface() // todo test for struct2table MapKey and MapValue
		if callback != nil {
			v2 := callback(field.Name, v)
			if v2 != "" {
				v = v2
			}
		}
		vv := fmt.Sprint(v)
		// if vv == "" {
		// 	vv = strconv.Quote("") // for struct2table
		// }
		fields = append(fields, CellData{
			Key:       field.Name,
			Value:     vv,
			Tooltip:   "",
			Icon:      nil,
			FgColor:   color.NRGBA{},
			IsNasm:    false,
			Disabled:  false,
			Clickable: widget.Clickable{},
			RichText:  RichText{},
		})
	}
	return
}

// UnmarshalFields 反序列化字段,fields理解成struckView或者form布局的所有行，每行都是一个键值对，key是字段名，value是字段值,callback理解成structView的每一行的自定义反序列化的callback
func UnmarshalFields[T any](fields []CellData, callback UnmarshalRowCallback) T {
	var zero T
	data := reflect.New(reflect.TypeOf(zero)).Interface()
	valueOf := reflect.ValueOf(data)
	if valueOf.Kind() != reflect.Ptr || valueOf.Elem().Kind() != reflect.Struct {
		panic("Data must be a pointer to struct")
	}
	valueOf = valueOf.Elem()
	for _, f := range fields {
		field, ok := valueOf.Type().FieldByName(f.Key)
		if !ok {
			continue
		}
		if !field.IsExported() ||
			field.Tag.Get("table") == "-" ||
			field.Tag.Get("json") == "-" {
			continue
		}
		fieldValue := valueOf.FieldByName(f.Key)
		if !fieldValue.CanSet() {
			panic("field is not settable: " + f.Key)
		}
		if callback != nil {
			if val := callback(f.Key, f.Value); val != nil {
				if rv := reflect.ValueOf(val); rv.Type().AssignableTo(fieldValue.Type()) {
					if s, ok := val.(string); ok && s == "" {
						panic("反序列化回调内应该返回nil或者非空字符串，请检查回调内的default返回值")
					}
					fieldValue.Set(rv)
					continue
				}
			}
		}
		switch field.Type.Kind() {
		case reflect.TypeFor[time.Duration]().Kind():
			fieldValue.Set(reflect.ValueOf(mylog.Check2(time.ParseDuration(f.Value))))
		case reflect.TypeFor[time.Time]().Kind():
			fieldValue.Set(reflect.ValueOf(mylog.Check2(time.Parse(time.RFC3339, f.Value))))
		case reflect.String:
			// if f.Value == strconv.Quote("") {
			// 	f.Value = mylog.Check2(strconv.Unquote(f.Value)) // todo test for struct2table
			// }
			fieldValue.SetString(f.Value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue.SetInt(mylog.Check2(strconv.ParseInt(f.Value, 10, 64))) // 如果是其他进制和位数呢？这就需要回调函数了
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// todo panic: strconv.ParseUint: parsing "string": invalid syntax in sturct2table
			// 编辑节点，点击应用按钮触发
			fieldValue.SetUint(mylog.Check2(strconv.ParseUint(f.Value, 10, 64)))
		case reflect.Float32, reflect.Float64:
			fieldValue.SetFloat(mylog.Check2(strconv.ParseFloat(f.Value, 64)))
		case reflect.Bool:
			fieldValue.SetBool(mylog.Check2(strconv.ParseBool(f.Value)))
		default:
			mylog.Check("unsupported field type: " + field.Type.Kind().String())
		}
	}
	// 老外的做法:和返回一个new(T)没区别吧？
	//	reflect.ValueOf(&e.beforeData).Elem().Set(reflect.New(reflect.TypeOf(e.beforeData).Elem()))
	//	e.beforeData.CopyFrom(target)
	//
	//	reflect.ValueOf(&e.editorData).Elem().Set(reflect.New(reflect.TypeOf(e.editorData).Elem()))
	//	e.editorData.CopyFrom(target)
	return reflect.Indirect(reflect.ValueOf(data)).Interface().(T)
}
