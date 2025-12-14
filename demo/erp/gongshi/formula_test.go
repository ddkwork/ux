package main

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
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
	fmt.Println("\n=== 按女工日结降序排序 ===")
	table.SortByColumn("女工日结", false)
	for i, row := range table.AllRows() {
		name := row.GetCell("姓名").Value
		day := row.GetCell("女工日结").Value
		fmt.Printf("%d. %v: %v\n", i+1, name, day)
	}
	table.ToMarkdown("按女工日结排序")
}
