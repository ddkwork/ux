package sdk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
	"github.com/xuri/excelize/v2"
)

// ===================== 类型推断系统 =====================
type SmartTypeInferencer struct {
	columnStats map[string]*columnStatistics
}

type columnStatistics struct {
	nonEmptyValues []string
	typeCounts     map[field.FieldType]int
}

func NewSmartTypeInferencer() *SmartTypeInferencer {
	return &SmartTypeInferencer{
		columnStats: make(map[string]*columnStatistics),
	}
}

func (s *SmartTypeInferencer) UpdateColumnStats(colName string, value string) {
	stats, exists := s.columnStats[colName]
	if !exists {
		stats = &columnStatistics{
			nonEmptyValues: make([]string, 0),
			typeCounts:     make(map[field.FieldType]int),
		}
		s.columnStats[colName] = stats
	}

	if strings.TrimSpace(value) == "" {
		return
	}

	stats.nonEmptyValues = append(stats.nonEmptyValues, value)
	inferredType := inferType2(value)
	stats.typeCounts[inferredType]++
}

func (s *SmartTypeInferencer) GetBestType(colName string) field.FieldType {
	stats, exists := s.columnStats[colName]
	if !exists || len(stats.nonEmptyValues) == 0 {
		return field.TextType
	}

	maxCount := 0
	bestType := field.TextType
	for typ, count := range stats.typeCounts {
		if count > maxCount || (count == maxCount && typPriority(typ) > typPriority(bestType)) {
			maxCount = count
			bestType = typ
		}
	}
	return bestType
}

func typPriority(t field.FieldType) int {
	switch t {
	case field.NumberType, field.CurrencyType, field.PercentType:
		return 4
	case field.DateTimeType:
		return 3
	case field.CheckboxType:
		return 2
	default:
		return 1
	}
}

func inferType2(value string) field.FieldType {
	value = strings.TrimSpace(value)
	if value == "" {
		return field.TextType
	}

	if isDate(value) {
		return field.DateTimeType
	}

	if isNumber(value) {
		if isCurrency(value) {
			return field.CurrencyType
		}
		if isPercentage(value) {
			return field.PercentType
		}
		return field.NumberType
	}

	if isBoolean(value) {
		return field.CheckboxType
	}

	if strings.Contains(value, "\n") {
		return field.MultiLineTextType
	}

	return field.TextType
}

func isDate(value string) bool {
	formats := []string{
		"2006-01-02", "2006/01/02", "01/02/2006", "02/01/2006",
		"2006-01-02 15:04:05", "2006/01/02 15:04:05",
		"2006年01月02日", "2006年1月2日",
	}
	for _, format := range formats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

func isNumber(value string) bool {
	_, err := strconv.ParseFloat(cleanNumber(value), 64)
	return err == nil
}

func isCurrency(value string) bool {
	return strings.ContainsAny(value, "¥$€£₽") && isNumber(cleanNumber(value))
}

func isPercentage(value string) bool {
	return strings.HasSuffix(value, "%") && isNumber(strings.TrimSuffix(value, "%"))
}

func isBoolean(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" ||
		lower == "yes" || lower == "no" ||
		lower == "是" || lower == "否" ||
		lower == "1" || lower == "0"
}

func cleanNumber(value string) string {
	clean := strings.ReplaceAll(value, "¥", "")
	clean = strings.ReplaceAll(clean, "$", "")
	clean = strings.ReplaceAll(clean, "€", "")
	clean = strings.ReplaceAll(clean, "£", "")
	clean = strings.ReplaceAll(clean, "₽", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, " ", "")
	return clean
}

// ===================== Excel 解析功能 =====================
func (t *TreeTable) LoadXlsx(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取工作表失败: %w", err)
	}

	if len(rows) < 1 {
		return fmt.Errorf("Excel文件为空")
	}

	// 初始化智能推断器
	inferencer := NewSmartTypeInferencer()
	headers := rows[0]

	// 收集列统计信息
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		for colIdx, header := range headers {
			if colIdx < len(row) {
				cellValue := row[colIdx]
				inferencer.UpdateColumnStats(header, cellValue)
			}
		}
	}

	// 创建列定义
	columns := make([]ColumnConfig, len(headers))
	for i, header := range headers {
		header = strings.TrimSpace(header)
		columns[i] = ColumnConfig{
			Name: header,
			Type: inferencer.GetBestType(header),
		}
	}

	// 解析数据行
	var dataRows [][]any
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		rowData := make([]any, len(headers))

		for colIdx := 0; colIdx < len(headers); colIdx++ {
			if colIdx < len(row) {
				cellValue := row[colIdx]
				colType := columns[colIdx].Type
				rowData[colIdx] = parseAnyValue(cellValue, colType)
			} else {
				rowData[colIdx] = nil
			}
		}

		dataRows = append(dataRows, rowData)
	}

	// 转换为TableData并加载
	data := TableData{
		Columns: columns,
		Rows:    dataRows,
	}
	t.LoadTableData(data)
	return nil
}

func parseAnyValue(value string, typ field.FieldType) any {
	if value == "" {
		return nil
	}

	switch typ {
	case field.DateTimeType:
		formats := []string{
			"2006-01-02", "2006/01/02", "01/02/2006", "02/01/2006",
			"2006-01-02 15:04:05", "2006/01/02 15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, value); err == nil {
				return t.Format("2006-01-02")
			}
		}
		return value

	case field.NumberType, field.CurrencyType, field.PercentType:
		clean := cleanNumber(value)
		if num, err := strconv.ParseFloat(clean, 64); err == nil {
			return num
		}
		return 0.0

	case field.CheckboxType:
		lower := strings.ToLower(value)
		if lower == "true" || lower == "yes" || lower == "是" || lower == "1" {
			return true
		}
		if lower == "false" || lower == "no" || lower == "否" || lower == "0" {
			return false
		}
		return false

	case field.MultiLineTextType:
		return value

	default: // TextType
		return value
	}
}
