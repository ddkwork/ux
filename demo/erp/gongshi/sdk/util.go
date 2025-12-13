package sdk

import (
	"fmt"
	"iter"
	"strconv"
	"time"

	"github.com/ddkwork/golibrary/std/stream/uuid"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

// ToFloat converts a value to float64.
func ToFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		if f, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64); err == nil {
			return f, true
		}
		return 0, false
	}
}

// TransposeMatrix 把行切片矩阵置换为列切片,用于计算最大列宽的参数
// https://github.com/BaseMax/SparseMatrixLinkedListGo
func TransposeMatrix[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for _, row := range rootRows {
			for j, v := range row {
				if !yield(j, v) { // j 是列索引，作为转置后的行键
					return
				}
			}
		}
	}
}

func newID() uuid.ID {
	return uuid.New('n')
}

// getDefaultValue returns the default value for a field type.
func getDefaultValue(ft field.FieldType) any {
	switch ft {
	case field.NumberType, field.CurrencyType, field.PercentType:
		return 0.0
	case field.CheckboxType:
		return false
	case field.DateTimeType:
		return time.Now()
	case field.SingleSelectType, field.MultipleSelectType:
		return ""
	case field.UrlType:
		return "https://"
	case field.EmailType:
		return "example@domain.com"
	case field.PhoneType:
		return "+1234567890"
	default:
		return ""
	}
}
