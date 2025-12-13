package sdk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

// CellData represents a cell in the table.
type CellData struct {
	ColumnName string          // Associated column name (key! find column definition via this field)
	Value      any             // Cell value
	Type       field.FieldType // Data type (can inherit from column definition)
}

// GetCellByRowIndex 通过行索引和列名获取单元格值
func (t *TreeTable) GetCellByRowIndex(rowIndex int, colName string) *CellData {
	row := t.GetRow(rowIndex)
	return row.GetCell(colName)
}

// GetCell retrieves a cell by column name.
func (n *Node) GetCell(colName string) *CellData {
	for i := range n.RowCells {
		if n.RowCells[i].ColumnName == colName {
			return &n.RowCells[i]
		}
	}
	panic("GetCell fail")
}

// SetCellValue sets the value of a cell by column name.
func (n *Node) SetCellValue(colName string, value any, table *TreeTable) {
	for i := range n.RowCells {
		if n.RowCells[i].ColumnName == colName {
			cell := &n.RowCells[i]
			cell.Value = value
			if cell.Type.Valid() {
				cell.Type = inferType(value)
			}
			return
		}
	}
	// Column does not exist, add it
	colDef := table.GetColumnDefinition(colName)
	if colDef != nil {
		n.RowCells = append(n.RowCells, CellData{
			ColumnName: colName,
			Value:      value,
			Type:       colDef.Type,
		})
	} else {
		n.RowCells = append(n.RowCells, CellData{
			ColumnName: colName,
			Value:      value,
			Type:       inferType(value),
		})
	}
}

// SetCellValue 通过行索引和列名设置单元格值
func (t *TreeTable) SetCellValue(rowIndex int, colName string, value any) bool {
	row := t.GetRow(rowIndex)
	if row == nil {
		return false
	}
	row.SetCellValue(colName, value, t)
	return true
}

// AsString returns the string representation of the cell value.
func (c *CellData) AsString() string {
	return c.Value.(string)
}

// AsInt returns the integer representation of the cell value.
func (c *CellData) AsInt() int {
	switch v := c.Value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case int64:
		return int(v)
	default:
		if i, err := strconv.Atoi(fmt.Sprintf("%v", c.Value)); err == nil {
			return i
		}
		panic("AsInt fail")
	}
}

// AsFloat returns the float64 representation of the cell value.
func (c *CellData) AsFloat() float64 {
	switch v := c.Value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		if f, err := strconv.ParseFloat(fmt.Sprintf("%v", c.Value), 64); err == nil {
			return f
		}
		panic("AsFloat fail")
	}
}

// AsBool returns the boolean representation of the cell value.
func (c *CellData) AsBool() bool {
	v, ok := c.Value.(bool)
	if ok {
		return v
	}

	// String representations of boolean
	if s, ok := c.Value.(string); ok {
		switch strings.ToLower(s) {
		case "true", "1", "yes", "是":
			return true
		case "false", "0", "no", "否":
			return false
		}
	}
	panic("AsBool fail")
}

// AsTime returns the time.Time representation of the cell value.
func (c *CellData) AsTime() time.Time {
	switch v := c.Value.(type) {
	case time.Time:
		return v
	case string:
		// Try multiple time formats
		formats := []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"01/02/2006",
			"2006/01/02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t
			}
		}
	}
	panic("AsTime fail")
}

// IsFormula checks if the cell is a formula column.
func (c *CellData) IsFormula() bool {
	return c.Type == field.FormulaType
}

// IsSelect checks if the cell is a select type (single or multiple).
func (c *CellData) IsSelect() bool {
	return c.Type == field.SingleSelectType || c.Type == field.MultipleSelectType
}

// IsAttachment checks if the cell is an attachment type.
func (c *CellData) IsAttachment() bool {
	return c.Type == field.AttachmentType
}

// IsLink checks if the cell is a link type.
func (c *CellData) IsLink() bool {
	return c.Type == field.LinkType
}

// DetectValueType detects the field type of the cell value.
func (c *CellData) DetectValueType() field.FieldType {
	return inferType(c.Value)
}

func inferType(v any) field.FieldType {
	if v == nil {
		return field.TextType
	}

	switch val := v.(type) {
	case string:
		return detectTypeFromString(val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return field.NumberType
	case float32, float64:
		return field.NumberType
	case bool:
		return field.CheckboxType
	case time.Time:
		return field.DateTimeType
	default:
		return detectTypeFromString(fmt.Sprintf("%v", v))
	}
}

// 从字符串值探测数据类型
func detectTypeFromString(s string) field.FieldType {
	// 尝试解析为布尔值
	if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") ||
		s == "1" || s == "0" || s == "是" || s == "否" {
		return field.CheckboxType
	}

	// 尝试解析为整数
	if _, err := strconv.Atoi(s); err == nil {
		return field.NumberType
	}

	// 尝试解析为浮点数
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return field.NumberType
	}

	// 尝试解析为日期时间 (RFC3339格式)
	if _, err := time.Parse(time.RFC3339, s); err == nil {
		return field.DateTimeType
	}

	// 尝试解析为简单日期格式 (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", s); err == nil {
		return field.DateTimeType
	}

	// 尝试解析为时间格式 (HH:MM:SS)
	if _, err := time.Parse("15:04:05", s); err == nil {
		return field.DateTimeType
	}

	// 尝试解析为日期时间组合 (YYYY-MM-DD HH:MM:SS)
	if _, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return field.DateTimeType
	}

	// 检查是否是URL
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return field.UrlType
	}

	// 检查是否是电子邮件
	if strings.Contains(s, "@") && strings.Contains(s, ".") {
		return field.EmailType
	}

	// 检查是否是电话号码 (简单验证)
	if len(s) >= 7 && len(s) <= 15 {
		allDigits := true
		for _, r := range s {
			if r < '0' || r > '9' {
				if r != '+' && r != '-' && r != '(' && r != ')' && r != ' ' {
					allDigits = false
					break
				}
			}
		}
		if allDigits {
			return field.PhoneType
		}
	}

	// 检查是否是多行文本 (包含换行符)
	if strings.Contains(s, "\n") {
		return field.MultiLineTextType
	}

	// 默认作为单行文本
	return field.TextType
}
