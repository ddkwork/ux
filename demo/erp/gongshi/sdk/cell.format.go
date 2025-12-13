package sdk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

func (c *CellData) FormatValue() string {
	switch c.Type {
	case field.TextType, field.MultiLineTextType:
		return c.AsString()
	case field.NumberType:
		return formatNumber(c.Value)
	case field.SingleSelectType:
		return formatSelect(c.Value, false)
	case field.MultipleSelectType:
		return formatSelect(c.Value, true)
	case field.DateTimeType:
		return formatDateTime(c.Value)
	case field.FormulaType:
		return fmt.Sprintf("=%v", c.Value)
	case field.AttachmentType:
		return formatAttachments(c.Value)
	case field.LinkType:
		return formatLink(c.Value)
	case field.UserType:
		return formatUsers(c.Value)
	case field.PhoneType:
		return formatPhone(c.Value)
	case field.EmailType:
		return c.AsString()
	case field.CheckboxType:
		return formatCheckbox(c.Value)
	case field.UrlType:
		return formatUrl(c.Value)
	case field.CurrencyType:
		return formatCurrency(c.Value)
	case field.PercentType:
		return formatPercent(c.Value)
	default:
		return fmt.Sprintf("%v", c.Value)
	}
}

// 辅助函数
func formatNumber(v any) string {
	switch num := v.(type) {
	case int:
		return strconv.Itoa(num)
	case float64:
		if num == float64(int(num)) {
			return strconv.Itoa(int(num))
		}
		return strconv.FormatFloat(num, 'f', -1, 64)
	case string:
		if f, err := strconv.ParseFloat(num, 64); err == nil {
			if f == float64(int(f)) {
				return strconv.Itoa(int(f))
			}
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
		return num
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatSelect(v any, multi bool) string {
	switch val := v.(type) {
	case []string:
		return strings.Join(val, ", ")
	case string:
		if multi {
			return strings.ReplaceAll(val, ";", ", ")
		}
		return val
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatDateTime(v any) string {
	switch t := v.(type) {
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	case string:
		if parsed, err := time.Parse("2006-01-02 15:04:05", t); err == nil {
			return parsed.Format("2006-01-02 15:04:05")
		}
		return t
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatCheckbox(v any) string {
	if b, ok := v.(bool); ok {
		return map[bool]string{true: "✓", false: ""}[b]
	}
	if s, ok := v.(string); ok {
		switch strings.ToLower(s) {
		case "true", "1", "yes", "是":
			return "✓"
		}
	}
	return ""
}

func formatCurrency(v any) string {
	numStr := formatNumber(v)
	if strings.Contains(numStr, ".") {
		return "¥" + numStr
	}
	return "¥" + numStr + ".00"
}

func formatPercent(v any) string {
	switch num := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f%%", num*100)
	case int:
		return fmt.Sprintf("%d%%", num)
	default:
		return fmt.Sprintf("%v%%", v)
	}
}

func formatUrl(v any) string {
	if url, ok := v.(string); ok {
		return fmt.Sprintf("<%s>", url)
	}
	return fmt.Sprintf("%v", v)
}

func formatPhone(v any) string {
	if phone, ok := v.(string); ok {
		return strings.ReplaceAll(phone, "-", "")
	}
	return fmt.Sprintf("%v", v)
}

func formatLink(v any) string {
	if link, ok := v.(map[string]string); ok {
		return fmt.Sprintf("[%s](%s)", link["text"], link["url"])
	}
	return fmt.Sprintf("%v", v)
}

func formatUsers(v any) string {
	switch users := v.(type) {
	case []string:
		return strings.Join(users, ", ")
	case string:
		return users
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatAttachments(v any) string {
	switch att := v.(type) {
	case []string:
		return fmt.Sprintf("%d个附件", len(att))
	case string:
		return att
	default:
		return fmt.Sprintf("%v", v)
	}
}
