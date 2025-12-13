package sdk

import "github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"

func tableDemo() *TreeTable {
	table := NewTreeTable()

	// 2. 直观的表格数据定义
	data := TableData{
		Columns: []ColumnConfig{
			{Name: "姓名", Type: field.TextType},
			{Name: "女工日结", Type: field.NumberType},
			{Name: "男工车结", Type: field.NumberType},
			{Name: "女工实发工资", Type: field.FormulaType, Formula: `
package main

import (
	"fmt"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
)

func RunScript(t *TreeTable, rowIndex int) {
	nameVal, ok := t.GetCellByRowIndex(rowIndex, "姓名")
	if !ok {
		return
	}
	name := fmt.Sprintf("%v", nameVal)

	nvGongVal, ok := t.GetCellByRowIndex(rowIndex, "女工日结")
	if !ok {
		return
	}
	nvGong, _ := ToFloat(nvGongVal)

	sanRenZuSum := t.SumIf("姓名", "三人组", "女工日结")
	switch name {
	case "拼车", "三人组":
		t.SetCellValue(rowIndex, "计算结果", 0.0)
	case "房东":
		t.SetCellValue(rowIndex, "计算结果", nvGong)
	case "杨萍":
		t.SetCellValue(rowIndex, "计算结果", (sanRenZuSum/3.0)+nvGong)
	case "二人组":
		t.SetCellValue(rowIndex, "计算结果", (sanRenZuSum/3.0)+(nvGong/2.0))
	default:
		t.SetCellValue(rowIndex, "计算结果", 0.0)
	}
}
`},
		},
		Rows: [][]any{
			{"三人组", 2966.30, 1104.20, 0},
			{"房东", 442.40, 196.80, 0},
			{"二人组", 5913.60, 2248.60, 0},
			{"杨萍", 3744.90, 1465.20, 0},
			{"拼车", 406.90, 175.00, 0},
		},
	}
	//姓名	女工日结	男工车结	女工实发工资
	//三人组	2966.30 1104.20 	0
	//房东	442.40 	196.80 	    442.4
	//二人组	5913.60 2248.60 	3945.566666667
	//杨萍	3744.90 1465.20 	4733.666666667
	//拼车	406.90 	175.00 	    0

	// 3. 一键设置数据
	table.LoadTableData(data)
	return table
}
