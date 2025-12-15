package main

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

func TestFormula(t *testing.T) {
	table := sdk.TableDemo()
	engine := NewYaegiEngine(table)
	for i := 0; i < table.RowCount(); i++ {
		engine.UpdateRowCell(i)
	}
	md := `# Tree Table Structure

| 姓名 | 女工日结 | 男工车结 | 女工实发工资 | 
|--------|--------|--------|--------|
| 📄 三人组 | 2966.3 | 1104.2 | 0 |
| 📄 房东 | 442.4 | 196.8 | 442.4 |
| 📄 二人组 | 5913.6 | 2248.6 | 3945.566666666667 |
| 📄 杨萍 | 3744.9 | 1465.2 | 4733.666666666667 |
| 📄 拼车 | 406.9 | 175 | 0 |
`

	assert.Equal(t, md, table.ToMarkdown("按公式列更新单元格数"))
	assert.Equal(t, "三人组", table.GetCellByRowIndex(0, "姓名").AsString())
	assert.Equal(t, 0, table.GetCellByRowIndex(0, "女工实发工资").AsFloat())
	assert.Equal(t, 442.4, table.GetCellByRowIndex(1, "女工实发工资").AsFloat())
	assert.Equal(t, 3945.566666666667, table.GetCellByRowIndex(2, "女工实发工资").AsFloat())
	assert.Equal(t, 4733.666666666667, table.GetCellByRowIndex(3, "女工实发工资").AsFloat())
	assert.Equal(t, 0, table.GetCellByRowIndex(4, "女工实发工资").AsFloat())

}

func TestSort(t *testing.T) {
	table := sdk.TableDemo()
	table.SortByColumn("女工日结", false)
	md := `# Tree Table Structure

| 姓名 | 女工日结 | 男工车结 | 女工实发工资 | 
|--------|--------|--------|--------|
| 📄 二人组 | 5913.6 | 2248.6 | 0 |
| 📄 杨萍 | 3744.9 | 1465.2 | 0 |
| 📄 三人组 | 2966.3 | 1104.2 | 0 |
| 📄 房东 | 442.4 | 196.8 | 0 |
| 📄 拼车 | 406.9 | 175 | 0 |
`
	assert.Equal(t, md, table.ToMarkdown("按女工日结排序"))
}

func TestGroupBy(t *testing.T) {
	table := sdk.NewTreeTable()
	data := sdk.TableData{
		Columns: []sdk.ColumnConfig{
			{Name: "姓名", Type: field.TextType},
			{Name: "女工日结", Type: field.NumberType},
			{Name: "男工车结", Type: field.NumberType},
			{Name: "女工实发工资", Type: field.NumberType},
		},
		Rows: [][]any{
			{"三人组", 2966.30, 1104.20, 0.0},
			{"房东", 442.40, 196.80, 442.4},
			{"二人组", 5913.60, 2248.60, 3945.57},
			{"杨萍", 3744.90, 1465.20, 4733.67},
			{"拼车", 406.90, 175.00, 0.0},
			{"三人组", 3000.00, 1200.00, 0.0}, // 另一个三人组
			{"房东", 500.00, 200.00, 500.0},   // 另一个房东
		},
	}

	table.LoadTableData(data)
	table.ToMarkdown("原始数据")

	table.GroupBy("姓名")
	md := `# Tree Table Structure

| 姓名 | 女工日结 | 男工车结 | 女工实发工资 | 
|--------|--------|--------|--------|
| 📂 三人组 (2) | 5966.3 | 2304.2 | 0 |
| &nbsp;&nbsp;&nbsp;📄 三人组 | 2966.3 | 1104.2 | 0 |
| &nbsp;&nbsp;&nbsp;📄 三人组 | 3000 | 1200 | 0 |
| 📂 二人组 (1) | 5913.6 | 2248.6 | 3945.57 |
| &nbsp;&nbsp;&nbsp;📄 二人组 | 5913.6 | 2248.6 | 3945.57 |
| 📂 房东 (2) | 942.4 | 396.8 | 942.4 |
| &nbsp;&nbsp;&nbsp;📄 房东 | 442.4 | 196.8 | 442.4 |
| &nbsp;&nbsp;&nbsp;📄 房东 | 500 | 200 | 500 |
| 📂 拼车 (1) | 406.9 | 175 | 0 |
| &nbsp;&nbsp;&nbsp;📄 拼车 | 406.9 | 175 | 0 |
| 📂 杨萍 (1) | 3744.9 | 1465.2 | 4733.67 |
| &nbsp;&nbsp;&nbsp;📄 杨萍 | 3744.9 | 1465.2 | 4733.67 |
`
	assert.Equal(t, md, table.ToMarkdown("按姓名分组集合"))

	aggResult := table.Aggregate("姓名", "女工日结", "sum")
	for group, sum := range aggResult {
		fmt.Printf("sum %s 女工日结: %.2f\n", group, sum)
	}

	table.Ungroup()
	table.ToMarkdown("取消分组")
}

func TestLoadXlsx(t *testing.T) {
	table := sdk.NewTreeTable()
	mylog.Check(table.LoadXlsx("(数据表)日结流水.xlsx"))
	table.ToMarkdown("TestLoadXlsx")
}
