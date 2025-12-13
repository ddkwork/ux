package sdk

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddkwork/golibrary/std/stream"
	"github.com/xuri/excelize/v2"
)

// TableData 表格数据（核心！）
type TableData struct {
	Columns []ColumnConfig // 列定义
	Rows    [][]any        // 行数据（每行对应一个数据记录）
}

// LoadTableData 一键加载表格数据（超直观！）
func (t *TreeTable) LoadTableData(data TableData) {
	// 清空现有数据
	t.Root = NewContainerNode("root", nil)
	t.Columns = make([]ColumnDefinition, len(data.Columns))
	t.columnMap = make(map[string]*ColumnDefinition)

	// 设置列定义
	for i, colCfg := range data.Columns {
		if colCfg.Formula != "" && !stream.IsAndroid() {
			stream.WriteGoFile(filepath.Join("tmp", colCfg.Name+"_column_script.go"), colCfg.Formula)
		}
		t.Columns[i] = ColumnDefinition{
			Name:         colCfg.Name,
			Type:         colCfg.Type,
			Formula:      colCfg.Formula,
			Options:      colCfg.Options,
			Width:        colCfg.Width,
			IsDisabled:   colCfg.Disabled,
			DefaultValue: getDefaultValue(colCfg.Type),
			Values:       nil,
		}
		t.columnMap[colCfg.Name] = &t.Columns[i]
	}

	// 添加行数据
	for _, rowData := range data.Rows {
		var cells []CellData
		for colIdx, cellValue := range rowData {
			if colIdx < len(data.Columns) {
				colName := data.Columns[colIdx].Name
				cells = append(cells, CellData{
					ColumnName: colName,
					Value:      cellValue,
					Type:       data.Columns[colIdx].Type,
				})
			}
		}
		t.Root.AddChild(NewNode(cells))
	}
}

func parseCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	return csvReader.ReadAll()
}

func parseXLS(filename string) ([][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func loadFile(filename string) ([][]string, error) {
	if strings.HasSuffix(filename, ".csv") {
		return parseCSV(filename)
	} else if strings.HasSuffix(filename, ".xls") || strings.HasSuffix(filename, ".xlsx") {
		return parseXLS(filename)
	}
	return nil, fmt.Errorf("unsupported file type")
}
