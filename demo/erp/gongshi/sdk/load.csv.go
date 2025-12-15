package sdk

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func (t *TreeTable) LoadCSV(filePath string) error {
	// 打开CSV文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开CSV文件失败: %w", err)
	}
	defer file.Close()

	// 创建CSV阅读器 - 使用制表符分隔
	reader := csv.NewReader(file)
	reader.Comma = '\t'         // 使用制表符分隔
	reader.FieldsPerRecord = -1 // 允许可变字段数

	// 读取所有行
	allRows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("解析CSV内容失败: %w", err)
	}

	if len(allRows) < 1 {
		return fmt.Errorf("CSV文件为空")
	}

	// 分离表头和数据行
	headers := allRows[0]   // 第一行
	dataRows := allRows[1:] // 剩余行是数据行

	// 初始化智能推断器
	inferencer := NewSmartTypeInferencer()

	// 关键修正：遍历数据行收集统计信息
	for _, row := range dataRows {
		for colIdx := 0; colIdx < len(headers); colIdx++ {
			if colIdx < len(row) {
				// 使用列名作为标识收集统计信息
				colName := headers[colIdx]
				cellValue := row[colIdx]
				inferencer.UpdateColumnStats(colName, cellValue)
			}
		}
	}

	// 创建列定义 - 使用收集到的统计信息
	columns := make([]ColumnConfig, len(headers))
	for colIdx, colName := range headers {
		trimmedName := strings.TrimSpace(colName)

		// 关键修正：使用列名获取最佳类型（基于数据行）
		bestType := inferencer.GetBestType(trimmedName)

		columns[colIdx] = ColumnConfig{
			Name: trimmedName,
			Type: bestType,
		}
	}

	// 解析数据行
	var parsedRows [][]any
	for _, row := range dataRows {
		parsedRow := make([]any, len(headers))
		for colIdx := 0; colIdx < len(headers); colIdx++ {
			if colIdx < len(row) {
				cellValue := row[colIdx]
				colType := columns[colIdx].Type
				parsedRow[colIdx] = parseAnyValue(cellValue, colType)
			} else {
				parsedRow[colIdx] = nil
			}
		}
		parsedRows = append(parsedRows, parsedRow)
	}

	// 转换为TableData并加载
	data := TableData{
		Columns: columns,
		Rows:    parsedRows,
	}
	t.LoadTableData(data)
	return nil
}
