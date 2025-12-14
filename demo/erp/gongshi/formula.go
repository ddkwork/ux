package main

import (
	"fmt"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ddkwork/ux/demo/erp/gongshi/sdk
var Symbols = interp.Exports{}

type YaegiEngine struct {
	interp *interp.Interpreter
	table  *sdk.TreeTable
}

func NewYaegiEngine(table *sdk.TreeTable) *YaegiEngine {
	i := interp.New(interp.Options{
		GoPath:       "./",
		Unrestricted: true,
	})
	i.Use(stdlib.Symbols)

	engine := &YaegiEngine{interp: i, table: table}
	return engine
}

func (e *YaegiEngine) UpdateRowCell(rowIndex int) {
	row := e.table.GetRow(rowIndex)
	if row == nil {
		panic("行不存在")
	}

	i := interp.New(interp.Options{
		GoPath:       "./",
		Unrestricted: true,
	})
	mylog.Check(i.Use(stdlib.Symbols))
	mylog.Check(i.Use(Symbols))

	for _, cell := range row.RowCells {
		if cell.IsFormula() {
			for _, column := range e.table.Columns {
				if cell.ColumnName == column.Name {
					mylog.Check2(i.Eval(column.Formula))
					runScript := mylog.Check2(i.Eval("RunScript")).Interface().(func(*sdk.TreeTable, int))
					runScript(e.table, rowIndex)
				}
			}
		}
	}
}

func main() {
	// 1. 创建表格
	table := sdk.NewTreeTable()

	// 2. 设置数据（包含重复姓名用于分组）
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
			{"房东", 500.00, 200.00, 500.0},  // 另一个房东
		},
	}

	table.LoadTableData(data)
	table.ToMarkdown("原始数据")

	table.GroupBy("姓名")
	table.ToMarkdown("按姓名分组")
	return

	// 5. 聚合计算
	fmt.Println("\n=== 分组聚合结果 ===")
	aggResult := table.Aggregate("姓名", "女工日结", "sum")
	for group, sum := range aggResult {
		fmt.Printf("%s 组女工日结总和: %.2f\n", group, sum)
	}

	// 6. 显示每个分组详情
	fmt.Println("\n=== 分组详情 ===")
	for _, group := range table.GetGroups() {
		groupName := group.GroupKey
		if cell := group.GetCell("姓名"); cell != nil {
			groupName = fmt.Sprintf("%v", cell.Value)
		}
		fmt.Printf("\n📁 分组: %s (%d人)\n", groupName, len(group.Children))

		for _, member := range group.Children {
			if name := member.GetCell("姓名"); name != nil {
				if day := member.GetCell("女工日结"); day != nil {
					fmt.Printf("  👤 %v: %v\n", name.Value, day.Value)
				}
			}
		}
	}
	table.ToMarkdown("按姓名分组集合")

	table.Ungroup()
	table.ToMarkdown("取消分组")
}
