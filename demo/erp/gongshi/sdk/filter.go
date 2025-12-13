package sdk

import (
	"fmt"
	"strings"
	"time"
)

// ------------------------------ 过滤系统 ------------------------------

// FilterOperator 定义过滤操作符
type FilterOperator string

const (
	OpEqual        FilterOperator = "="
	OpNotEqual     FilterOperator = "!="
	OpGreaterThan  FilterOperator = ">"
	OpLessThan     FilterOperator = "<"
	OpGreaterEqual FilterOperator = ">="
	OpLessEqual    FilterOperator = "<="
	OpContains     FilterOperator = "contains"
	OpStartsWith   FilterOperator = "startswith"
	OpEndsWith     FilterOperator = "endswith"
	OpIn           FilterOperator = "in"
	OpNotIn        FilterOperator = "notin"
	OpIsEmpty      FilterOperator = "isempty"
	OpIsNotEmpty   FilterOperator = "isnotempty"
	OpIsTrue       FilterOperator = "istrue"
	OpIsFalse      FilterOperator = "isfalse"
)

// FilterCondition 定义单个过滤条件
type FilterCondition struct {
	Column   string         // 列名
	Operator FilterOperator // 操作符
	Value    any            // 比较值
}

// FilterGroup 定义过滤条件组（支持AND/OR组合）
type FilterGroup struct {
	Conditions []FilterCondition // 条件列表
	SubGroups  []*FilterGroup    // 子条件组
	Logic      string            // "AND" 或 "OR"
}

// FilterOptions 过滤配置选项
type FilterOptions struct {
	CaseSensitive bool             // 是否区分大小写
	DatePrecision string           // 日期精度: "day", "month", "year"
	CustomEval    func(*Node) bool // 自定义过滤函数
}

// 默认过滤选项
var DefaultFilterOptions = FilterOptions{
	CaseSensitive: false,
	DatePrecision: "day",
}

// ------------------------------ 过滤方法实现 ------------------------------

// 应用过滤条件到单个节点
func (n *Node) ApplyFilter(condition FilterCondition, options FilterOptions, table *TreeTable) bool {
	cell := n.GetCell(condition.Column, table)
	if cell == nil {
		return false
	}

	// 处理特殊操作符
	switch condition.Operator {
	case OpIsEmpty:
		return cell.Value == nil || fmt.Sprint(cell.Value) == ""
	case OpIsNotEmpty:
		return cell.Value != nil && fmt.Sprint(cell.Value) != ""
	case OpIsTrue:
		if b, ok := cell.AsBool(); ok {
			return b
		}
		return false
	case OpIsFalse:
		if b, ok := cell.AsBool(); ok {
			return !b
		}
		return false
	}

	// 获取比较值
	compareValue := condition.Value
	if !options.CaseSensitive && compareValue != nil {
		if s, ok := compareValue.(string); ok {
			compareValue = strings.ToLower(s)
		}
	}

	// 获取单元格值
	cellValue := cell.Value
	if cellValue == nil {
		return false
	}

	// 处理字符串比较
	if s, ok := cellValue.(string); ok {
		if !options.CaseSensitive {
			s = strings.ToLower(s)
		}
		return applyStringFilter(s, compareValue, condition.Operator)
	}

	// 处理数值比较
	if num, ok := ToFloat(cellValue); ok {
		if compareNum, ok := ToFloat(compareValue); ok {
			return applyNumericFilter(num, compareNum, condition.Operator)
		}
	}

	// 处理日期比较
	if dateStr, ok := cellValue.(string); ok {
		if parsedDate, err := parseDate(dateStr); err == nil {
			if compareDate, ok := compareValue.(time.Time); ok {
				return applyDateFilter(parsedDate, compareDate, condition.Operator, options.DatePrecision)
			}
			if compareStr, ok := compareValue.(string); ok {
				if parsedCompare, err := parseDate(compareStr); err == nil {
					return applyDateFilter(parsedDate, parsedCompare, condition.Operator, options.DatePrecision)
				}
			}
		}
	}

	// 默认比较
	switch condition.Operator {
	case OpEqual:
		return fmt.Sprint(cellValue) == fmt.Sprint(compareValue)
	case OpNotEqual:
		return fmt.Sprint(cellValue) != fmt.Sprint(compareValue)
	case OpIn:
		return applyInFilter(cellValue, compareValue, false, options.CaseSensitive)
	case OpNotIn:
		return applyInFilter(cellValue, compareValue, true, options.CaseSensitive)
	default:
		return false
	}
}

// 应用字符串过滤
func applyStringFilter(cellValue, compareValue any, op FilterOperator) bool {
	strVal := fmt.Sprint(cellValue)
	compVal := fmt.Sprint(compareValue)

	switch op {
	case OpContains:
		return strings.Contains(strVal, compVal)
	case OpStartsWith:
		return strings.HasPrefix(strVal, compVal)
	case OpEndsWith:
		return strings.HasSuffix(strVal, compVal)
	case OpEqual:
		return strVal == compVal
	case OpNotEqual:
		return strVal != compVal
	default:
		return false
	}
}

// 应用数值过滤
func applyNumericFilter(cellValue, compareValue float64, op FilterOperator) bool {
	switch op {
	case OpEqual:
		return cellValue == compareValue
	case OpNotEqual:
		return cellValue != compareValue
	case OpGreaterThan:
		return cellValue > compareValue
	case OpLessThan:
		return cellValue < compareValue
	case OpGreaterEqual:
		return cellValue >= compareValue
	case OpLessEqual:
		return cellValue <= compareValue
	default:
		return false
	}
}

// 应用日期过滤
func applyDateFilter(cellDate, compareDate time.Time, op FilterOperator, precision string) bool {
	// 根据精度调整日期
	if precision == "month" {
		cellDate = time.Date(cellDate.Year(), cellDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		compareDate = time.Date(compareDate.Year(), compareDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else if precision == "year" {
		cellDate = time.Date(cellDate.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		compareDate = time.Date(compareDate.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}

	switch op {
	case OpEqual:
		return cellDate.Equal(compareDate)
	case OpNotEqual:
		return !cellDate.Equal(compareDate)
	case OpGreaterThan:
		return cellDate.After(compareDate)
	case OpLessThan:
		return cellDate.Before(compareDate)
	case OpGreaterEqual:
		return cellDate.After(compareDate) || cellDate.Equal(compareDate)
	case OpLessEqual:
		return cellDate.Before(compareDate) || cellDate.Equal(compareDate)
	default:
		return false
	}
}

// 应用IN/NOT IN过滤
func applyInFilter(cellValue, compareValue any, negate bool, caseSensitive bool) bool {
	// 处理比较值（应为切片或数组）
	compSlice, ok := convertToSlice(compareValue)
	if !ok {
		return false
	}

	// 获取单元格值字符串表示
	cellStr := fmt.Sprint(cellValue)
	if !caseSensitive {
		cellStr = strings.ToLower(cellStr)
	}

	// 检查是否在列表中
	found := false
	for _, item := range compSlice {
		itemStr := fmt.Sprint(item)
		if !caseSensitive {
			itemStr = strings.ToLower(itemStr)
		}
		if itemStr == cellStr {
			found = true
			break
		}
	}

	return found != negate // XOR操作
}

// 将任意值转换为切片
func convertToSlice(v any) ([]any, bool) {
	switch slice := v.(type) {
	case []any:
		return slice, true
	case []string:
		result := make([]any, len(slice))
		for i, s := range slice {
			result[i] = s
		}
		return result, true
	case []int:
		result := make([]any, len(slice))
		for i, num := range slice {
			result[i] = num
		}
		return result, true
	case []float64:
		result := make([]any, len(slice))
		for i, num := range slice {
			result[i] = num
		}
		return result, true
	default:
		// 尝试使用反射处理其他类型
		return nil, false
	}
}

// 解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}

// 评估过滤组
func evaluateFilterGroup(node *Node, group *FilterGroup, options FilterOptions, table *TreeTable) bool {
	// 评估当前组的条件
	currentResult := true
	if len(group.Conditions) > 0 {
		for _, cond := range group.Conditions {
			condResult := node.ApplyFilter(cond, options, table)
			if group.Logic == "OR" {
				currentResult = currentResult || condResult
				if currentResult {
					break // OR逻辑下有一个为真即为真
				}
			} else { // AND
				currentResult = currentResult && condResult
				if !currentResult {
					break // AND逻辑下有一个为假即为假
				}
			}
		}
	} else {
		currentResult = true // 无条件组默认为真
	}

	// 评估子组
	if len(group.SubGroups) > 0 {
		subResult := false
		for i, subGroup := range group.SubGroups {
			groupResult := evaluateFilterGroup(node, subGroup, options, table)
			if i == 0 {
				subResult = groupResult
			} else {
				if group.Logic == "OR" {
					subResult = subResult || groupResult
				} else {
					subResult = subResult && groupResult
				}
			}
		}

		// 组合当前结果和子组结果
		if group.Logic == "OR" {
			currentResult = currentResult || subResult
		} else {
			currentResult = currentResult && subResult
		}
	}

	return currentResult
}

// ------------------------------ TreeTable 过滤方法 ------------------------------

// FilterNodes 根据过滤条件返回匹配的节点
func (t *TreeTable) FilterNodes(group *FilterGroup, options FilterOptions) []*Node {
	var results []*Node

	// 使用自定义过滤函数（如果存在）
	if options.CustomEval != nil {
		for node := range t.Root.Walk() {
			if options.CustomEval(node) {
				results = append(results, node)
			}
		}
		return results
	}

	// 使用过滤组
	for node := range t.Root.Walk() {
		if evaluateFilterGroup(node, group, options, t) {
			results = append(results, node)
		}
	}

	return results
}

// FilterByColumn 单列简单过滤
func (t *TreeTable) FilterByColumn(column string, operator FilterOperator, value any, options FilterOptions) []*Node {
	group := &FilterGroup{
		Conditions: []FilterCondition{
			{Column: column, Operator: operator, Value: value},
		},
		Logic: "AND",
	}
	return t.FilterNodes(group, options)
}

// FilterByText 文本搜索（在所有列中搜索）
func (t *TreeTable) FilterByText(searchText string, options FilterOptions) []*Node {
	group := &FilterGroup{Logic: "OR"}

	// 为每个文本列添加条件
	for _, col := range t.Columns {
		if col.Type == FieldTypeSingleLineText ||
			col.Type == FieldTypeMultiLineText ||
			col.Type == FieldTypeEmail ||
			col.Type == FieldTypeURL {
			group.Conditions = append(group.Conditions, FilterCondition{
				Column:   col.Name,
				Operator: OpContains,
				Value:    searchText,
			})
		}
	}

	// 如果没有文本列，使用所有列
	if len(group.Conditions) == 0 {
		for _, col := range t.Columns {
			group.Conditions = append(group.Conditions, FilterCondition{
				Column:   col.Name,
				Operator: OpContains,
				Value:    searchText,
			})
		}
	}

	return t.FilterNodes(group, options)
}

// FilterAdvanced 高级过滤（示例）
func (t *TreeTable) FilterAdvanced() []*Node {
	// 示例：查找年龄在30岁以上且在职的员工
	group := &FilterGroup{
		Logic: "AND",
		Conditions: []FilterCondition{
			{Column: "年龄", Operator: OpGreaterThan, Value: 30},
			{Column: "状态", Operator: OpEqual, Value: "在职"},
		},
	}
	return t.FilterNodes(group, DefaultFilterOptions)
}

// FilterByDateRange 日期范围过滤
func (t *TreeTable) FilterByDateRange(column string, start, end time.Time, options FilterOptions) []*Node {
	group := &FilterGroup{
		Logic: "AND",
		Conditions: []FilterCondition{
			{Column: column, Operator: OpGreaterEqual, Value: start},
			{Column: column, Operator: OpLessEqual, Value: end},
		},
	}
	return t.FilterNodes(group, options)
}

// FilterByList 列表过滤（值在列表中）
func (t *TreeTable) FilterByList(column string, allowedValues []string, options FilterOptions) []*Node {
	group := &FilterGroup{
		Conditions: []FilterCondition{
			{Column: column, Operator: OpIn, Value: allowedValues},
		},
		Logic: "AND",
	}
	return t.FilterNodes(group, options)
}

// ClearFilters 清除所有过滤条件
func (t *TreeTable) ClearFilters() []*Node {
	return t.AllRows()
}

// ------------------------------ 实用工具函数 ------------------------------

// 创建过滤组辅助函数
func NewFilterGroup(logic string, conditions []FilterCondition, subGroups ...*FilterGroup) *FilterGroup {
	return &FilterGroup{
		Conditions: conditions,
		SubGroups:  subGroups,
		Logic:      logic,
	}
}

// 创建过滤条件辅助函数
func NewFilterCondition(column string, operator FilterOperator, value any) FilterCondition {
	return FilterCondition{
		Column:   column,
		Operator: operator,
		Value:    value,
	}
}

// 示例用法
func ExampleUsage() {
	table := NewTreeTable()

	// 简单过滤：查找姓名为"张三"的员工
	_ = table.FilterByColumn("姓名", OpEqual, "张三", DefaultFilterOptions)

	// 组合过滤：年龄在25-35岁之间且状态为在职
	ageGroup := NewFilterGroup("AND", []FilterCondition{
		NewFilterCondition("年龄", OpGreaterEqual, 25),
		NewFilterCondition("年龄", OpLessEqual, 35),
	})
	statusGroup := NewFilterGroup("AND", []FilterCondition{
		NewFilterCondition("状态", OpEqual, "在职"),
	})

	combinedGroup := NewFilterGroup("AND", nil, ageGroup, statusGroup)
	_ = table.FilterNodes(combinedGroup, DefaultFilterOptions)

	// 文本搜索：在所有文本列中搜索"技术"
	_ = table.FilterByText("技术", DefaultFilterOptions)

	// 日期范围：入职日期在2020年之后
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()
	_ = table.FilterByDateRange("入职日期", startDate, endDate, DefaultFilterOptions)
}
