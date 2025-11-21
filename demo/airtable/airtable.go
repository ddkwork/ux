package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 生成唯一ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// FieldType 字段类型枚举
type FieldType int

const (
	FieldTypeSingleLineText FieldType = iota
	FieldTypeNumber
	FieldTypeSingleSelect
	FieldTypeMultipleSelect
	FieldTypeDateTime
	FieldTypeFormula
	FieldTypeAttachment
	FieldTypeLink
	FieldTypeUser
)

// Field 字段定义
type Field struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Options     []string  `json:"options,omitempty"` // 用于选择型字段
	FormulaExpr string    `json:"formula_expr,omitempty"`
	IsPrimary   bool      `json:"is_primary"`
}

// Record 记录/行
type Record struct {
	ID          string                 `json:"id"`
	Values      map[string]interface{} `json:"values"` // fieldId -> value
	CreatedTime int64                  `json:"created_time"`
	UpdatedTime int64                  `json:"updated_time"`
}

// Table 表结构
type Table struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Fields  []*Field  `json:"fields"`
	Records []*Record `json:"records"`
	// 用于快速查找字段
	fieldMap map[string]*Field
}

// GroupResult 分组结果
type GroupResult struct {
	GroupValue interface{}        `json:"group_value"`
	Records    []*Record          `json:"records"`
	Sums       map[string]float64 `json:"sums"` // 字段ID -> 求和结果
	Count      int                `json:"count"`
}

// 初始化表
func NewTable(name string) *Table {
	return &Table{
		ID:       generateID(),
		Name:     name,
		Fields:   make([]*Field, 0),
		Records:  make([]*Record, 0),
		fieldMap: make(map[string]*Field),
	}
}

// 添加字段
func (t *Table) AddField(name string, fieldType FieldType, options ...string) (*Field, error) {
	field := &Field{
		ID:      generateID(),
		Name:    name,
		Type:    fieldType,
		Options: options,
	}

	t.Fields = append(t.Fields, field)
	t.fieldMap[field.ID] = field
	return field, nil
}

// 获取字段
func (t *Table) GetField(fieldID string) *Field {
	return t.fieldMap[fieldID]
}

// 查找字段
func (t *Table) FindFieldByName(name string) *Field {
	for _, field := range t.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

// 验证字段值
func (t *Table) validateFieldValue(fieldID string, value interface{}) error {
	field, exists := t.fieldMap[fieldID]
	if !exists {
		return fmt.Errorf("field with ID %s does not exist", fieldID)
	}

	// 根据字段类型验证值
	switch field.Type {
	case FieldTypeSingleSelect, FieldTypeMultipleSelect:
		if strVal, ok := value.(string); ok {
			// 检查值是否在选项中
			found := false
			for _, option := range field.Options {
				if option == strVal {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("value %s is not a valid option for field %s", strVal, field.Name)
			}
		} else {
			return fmt.Errorf("field %s expects a string value", field.Name)
		}
	case FieldTypeNumber:
		if _, ok := value.(float64); !ok {
			if intVal, ok := value.(int); ok {
				intVal = intVal
				// 允许整型值，将在存储时转换为float64
			} else {
				return fmt.Errorf("field %s expects a numeric value", field.Name)
			}
		}
	}

	return nil
}

// 添加记录
func (t *Table) AddRecord(values map[string]interface{}) (*Record, error) {
	record := &Record{
		ID:          generateID(),
		Values:      make(map[string]interface{}),
		CreatedTime: getCurrentTimestamp(),
		UpdatedTime: getCurrentTimestamp(),
	}

	// 验证并设置值
	for fieldID, value := range values {
		if err := t.validateFieldValue(fieldID, value); err != nil {
			return nil, err
		}

		// 转换整型值为浮点型
		if intVal, ok := value.(int); ok {
			if field := t.GetField(fieldID); field != nil && field.Type == FieldTypeNumber {
				record.Values[fieldID] = float64(intVal)
			} else {
				record.Values[fieldID] = value
			}
		} else {
			record.Values[fieldID] = value
		}
	}

	// 计算公式字段
	if err := t.calculateFormulaFields(record); err != nil {
		return nil, err
	}

	t.Records = append(t.Records, record)
	return record, nil
}

// 计算公式字段
func (t *Table) calculateFormulaFields(record *Record) error {
	engine := NewFormulaEngine(t)

	for _, field := range t.Fields {
		if field.Type == FieldTypeFormula && field.FormulaExpr != "" {
			result, err := engine.Evaluate(field.FormulaExpr, record)
			if err != nil {
				return err
			}
			record.Values[field.ID] = result
		}
	}

	return nil
}

// 更新记录
func (t *Table) UpdateRecord(recordID string, values map[string]interface{}) error {
	var record *Record
	for _, r := range t.Records {
		if r.ID == recordID {
			record = r
			break
		}
	}

	if record == nil {
		return fmt.Errorf("record with ID %s not found", recordID)
	}

	// 验证并更新值
	for fieldID, value := range values {
		if err := t.validateFieldValue(fieldID, value); err != nil {
			return err
		}
		record.Values[fieldID] = value
	}

	// 重新计算公式字段
	if err := t.calculateFormulaFields(record); err != nil {
		return err
	}

	record.UpdatedTime = getCurrentTimestamp()
	return nil
}

// FormulaEngine 公式引擎
type FormulaEngine struct {
	table *Table
}

func NewFormulaEngine(table *Table) *FormulaEngine {
	return &FormulaEngine{table: table}
}

// Evaluate 计算公式
func (e *FormulaEngine) Evaluate(expr string, record *Record) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	// 处理字段引用: {fieldName}
	if strings.HasPrefix(expr, "{") && strings.HasSuffix(expr, "}") {
		fieldName := strings.Trim(expr, "{}")
		return e.getFieldValue(fieldName, record)
	}

	// 处理基本运算
	if strings.Contains(expr, "+") {
		return e.evaluateBinaryOp(expr, "+", record, func(a, b float64) float64 { return a + b })
	}
	if strings.Contains(expr, "-") {
		return e.evaluateBinaryOp(expr, "-", record, func(a, b float64) float64 { return a - b })
	}
	if strings.Contains(expr, "*") {
		return e.evaluateBinaryOp(expr, "*", record, func(a, b float64) float64 { return a * b })
	}
	if strings.Contains(expr, "/") {
		return e.evaluateBinaryOp(expr, "/", record, func(a, b float64) float64 {
			if b == 0 {
				return 0 // 避免除零错误
			}
			return a / b
		})
	}

	// 处理函数调用
	if strings.Contains(expr, "(") && strings.Contains(expr, ")") {
		return e.evaluateFunction(expr, record)
	}

	// 默认返回字符串
	return expr, nil
}

// 获取字段值
func (e *FormulaEngine) getFieldValue(fieldName string, record *Record) (interface{}, error) {
	for _, field := range e.table.Fields {
		if field.Name == fieldName {
			if value, exists := record.Values[field.ID]; exists {
				return value, nil
			}
			return nil, nil
		}
	}
	return nil, nil
}

// 计算二元运算
func (e *FormulaEngine) evaluateBinaryOp(expr, op string, record *Record, opFunc func(float64, float64) float64) (float64, error) {
	parts := strings.Split(expr, op)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid binary operation: %s", expr)
	}

	left, err := e.Evaluate(strings.TrimSpace(parts[0]), record)
	if err != nil {
		return 0, err
	}

	right, err := e.Evaluate(strings.TrimSpace(parts[1]), record)
	if err != nil {
		return 0, err
	}

	leftNum, err := toNumber(left)
	if err != nil {
		return 0, err
	}

	rightNum, err := toNumber(right)
	if err != nil {
		return 0, err
	}

	return opFunc(leftNum, rightNum), nil
}

// 计算函数
func (e *FormulaEngine) evaluateFunction(expr string, record *Record) (interface{}, error) {
	funcName := strings.Split(expr, "(")[0]
	funcName = strings.ToUpper(strings.TrimSpace(funcName))

	// 提取参数
	start := strings.Index(expr, "(")
	end := strings.LastIndex(expr, ")")
	argsStr := expr[start+1 : end]
	args := strings.Split(argsStr, ",")
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	switch funcName {
	case "SUM":
		return e.sum(args, record)
	case "AVERAGE":
		return e.average(args, record)
	case "MAX":
		return e.max(args, record)
	case "MIN":
		return e.min(args, record)
	case "IF":
		return e.ifFunc(args, record)
	case "CONCAT":
		return e.concat(args, record)
	default:
		return nil, fmt.Errorf("unknown function: %s", funcName)
	}
}

// SUM 函数实现
func (e *FormulaEngine) sum(args []string, record *Record) (float64, error) {
	var sum float64
	for _, arg := range args {
		val, err := e.Evaluate(arg, record)
		if err != nil {
			return 0, err
		}
		num, err := toNumber(val)
		if err != nil {
			return 0, err
		}
		sum += num
	}
	return sum, nil
}

// AVERAGE 函数实现
func (e *FormulaEngine) average(args []string, record *Record) (float64, error) {
	sum, err := e.sum(args, record)
	if err != nil {
		return 0, err
	}
	return sum / float64(len(args)), nil
}

// MAX 函数实现
func (e *FormulaEngine) max(args []string, record *Record) (float64, error) {
	maxVal := math.Inf(-1)
	for _, arg := range args {
		val, err := e.Evaluate(arg, record)
		if err != nil {
			return 0, err
		}
		num, err := toNumber(val)
		if err != nil {
			return 0, err
		}
		if num > maxVal {
			maxVal = num
		}
	}
	return maxVal, nil
}

// MIN 函数实现
func (e *FormulaEngine) min(args []string, record *Record) (float64, error) {
	minVal := math.Inf(1)
	for _, arg := range args {
		val, err := e.Evaluate(arg, record)
		if err != nil {
			return 0, err
		}
		num, err := toNumber(val)
		if err != nil {
			return 0, err
		}
		if num < minVal {
			minVal = num
		}
	}
	return minVal, nil
}

// IF 函数实现
func (e *FormulaEngine) ifFunc(args []string, record *Record) (interface{}, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("IF requires 3 arguments")
	}

	condition, err := e.Evaluate(args[0], record)
	if err != nil {
		return nil, err
	}

	// 检查条件是否为真
	if isTruthy(condition) {
		return e.Evaluate(args[1], record)
	}
	return e.Evaluate(args[2], record)
}

// CONCAT 函数实现
func (e *FormulaEngine) concat(args []string, record *Record) (string, error) {
	var result strings.Builder
	for _, arg := range args {
		val, err := e.Evaluate(arg, record)
		if err != nil {
			return "", err
		}
		result.WriteString(fmt.Sprintf("%v", val))
	}
	return result.String(), nil
}

// 转换为数字
func toNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to number", value)
	}
}

// 检查值是否为真
func isTruthy(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v != "" && v != "false"
	default:
		return false
	}
}

// 分组记录
func (t *Table) GroupBy(fieldID string, sumFieldIDs ...string) ([]*GroupResult, error) {
	// 检查字段是否存在
	groupField := t.GetField(fieldID)
	if groupField == nil {
		return nil, fmt.Errorf("field with ID %s does not exist", fieldID)
	}

	// 检查求和字段是否存在
	for _, sumFieldID := range sumFieldIDs {
		if sumField := t.GetField(sumFieldID); sumField == nil {
			return nil, fmt.Errorf("sum field with ID %s does not exist", sumFieldID)
		}
	}

	// 按分组字段的值分组记录
	groups := make(map[interface{}]*GroupResult)

	for _, record := range t.Records {
		groupValue, exists := record.Values[fieldID]
		if !exists {
			continue
		}

		// 初始化分组结果
		if _, exists := groups[groupValue]; !exists {
			groups[groupValue] = &GroupResult{
				GroupValue: groupValue,
				Records:    make([]*Record, 0),
				Sums:       make(map[string]float64),
				Count:      0,
			}
		}

		// 添加记录到分组
		group := groups[groupValue]
		group.Records = append(group.Records, record)
		group.Count++

		// 计算求和字段
		for _, sumFieldID := range sumFieldIDs {
			if value, exists := record.Values[sumFieldID]; exists {
				if num, err := toNumber(value); err == nil {
					group.Sums[sumFieldID] += num
				}
			}
		}
	}

	// 转换为切片并排序
	result := make([]*GroupResult, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	// 按分组值排序
	sort.Slice(result, func(i, j int) bool {
		return fmt.Sprintf("%v", result[i].GroupValue) < fmt.Sprintf("%v", result[j].GroupValue)
	})

	return result, nil
}

// 打印表格
func (t *Table) Print() {
	// 打印表头
	headers := make([]string, len(t.Fields))
	for i, field := range t.Fields {
		headers[i] = field.Name
	}
	fmt.Println(strings.Join(headers, "\t"))

	// 打印分隔线
	separators := make([]string, len(t.Fields))
	for i := range separators {
		separators[i] = "--------"
	}
	fmt.Println(strings.Join(separators, "\t"))

	// 打印记录
	for _, record := range t.Records {
		values := make([]string, len(t.Fields))
		for i, field := range t.Fields {
			if value, exists := record.Values[field.ID]; exists {
				values[i] = fmt.Sprintf("%v", value)
			} else {
				values[i] = ""
			}
		}
		fmt.Println(strings.Join(values, "\t"))
	}
}

// 打印分组结果
func PrintGroupResults(results []*GroupResult, groupFieldName string, table *Table) {
	for _, group := range results {
		fmt.Printf("\n%s: %v (Count: %d)\n", groupFieldName, group.GroupValue, group.Count)

		// 打印求和结果
		if len(group.Sums) > 0 {
			fmt.Println("Sums:")
			for fieldID, sum := range group.Sums {
				if field := table.GetField(fieldID); field != nil {
					fmt.Printf("  %s: %.2f\n", field.Name, sum)
				}
			}
		}

		// 打印分组中的记录
		fmt.Println("Records:")
		for _, record := range group.Records {
			values := make([]string, 0)
			for _, field := range table.Fields {
				if value, exists := record.Values[field.ID]; exists {
					values = append(values, fmt.Sprintf("%s: %v", field.Name, value))
				}
			}
			fmt.Printf("  - %s\n", strings.Join(values, ", "))
		}
	}
}

// 导出为JSON
func (t *Table) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func tableEntry() {
	// 创建新表
	table := NewTable("销售记录")

	// 添加字段
	productField, _ := table.AddField("产品", FieldTypeSingleLineText)
	categoryField, _ := table.AddField("类别", FieldTypeSingleSelect, "电子产品", "家具", "服装")
	quantityField, _ := table.AddField("数量", FieldTypeNumber)
	priceField, _ := table.AddField("单价", FieldTypeNumber)

	// 添加公式字段
	totalField, _ := table.AddField("总价", FieldTypeFormula)
	totalField.FormulaExpr = "{数量} * {单价}"

	// 添加记录
	table.AddRecord(map[string]interface{}{
		productField.ID:  "iPhone 13",
		categoryField.ID: "电子产品",
		quantityField.ID: 5,
		priceField.ID:    799.99,
	})

	table.AddRecord(map[string]interface{}{
		productField.ID:  "MacBook Pro",
		categoryField.ID: "电子产品",
		quantityField.ID: 3,
		priceField.ID:    1999.99,
	})

	table.AddRecord(map[string]interface{}{
		productField.ID:  "沙发",
		categoryField.ID: "家具",
		quantityField.ID: 2,
		priceField.ID:    899.99,
	})

	table.AddRecord(map[string]interface{}{
		productField.ID:  "餐桌",
		categoryField.ID: "家具",
		quantityField.ID: 1,
		priceField.ID:    499.99,
	})

	table.AddRecord(map[string]interface{}{
		productField.ID:  "T恤",
		categoryField.ID: "服装",
		quantityField.ID: 10,
		priceField.ID:    29.99,
	})

	// 打印表格
	fmt.Println("销售记录表:")
	table.Print()

	// 按类别分组并计算数量和总价的总和
	fmt.Println("\n\n按类别分组结果:")
	groups, err := table.GroupBy(categoryField.ID, quantityField.ID, totalField.ID)
	if err != nil {
		fmt.Printf("分组错误: %v\n", err)
		return
	}

	PrintGroupResults(groups, "类别", table)

	// 导出为JSON
	jsonStr, err := table.ToJSON()
	if err == nil {
		fmt.Println("\nJSON表示:")
		fmt.Println(jsonStr)
	}
}
