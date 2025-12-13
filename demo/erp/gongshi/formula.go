package main

import (
	"fmt"
	"strings"

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
	table := sdk.TableDemo()
	table.ToMarkdown("原始数据")

	// 4. 显示数据
	//fmt.Println("=== 基础数据展示 ===")
	//fmt.Printf("%-8s %-12s %-12s %-16s\n", "姓名", "女工日结", "男工车结", "女工实发工资")
	//fmt.Println("────────── ──────────── ──────────── ────────────────")

	//for _, row := range table.AllRows() {
	//	name := row.GetCell("姓名").Value
	//	day := row.GetCell("女工日结").Value
	//	car := row.GetCell("男工车结").Value
	//	salary := row.GetCell("女工实发工资").Value
	//	fmt.Printf("%-8v %-12v %-12v %-16v\n", name, day, car, salary)
	//}

	// 5. 排序演示
	//fmt.Println("\n=== 按女工日结降序排序 ===")
	//table.SortByColumn("女工日结", false)
	//for i, row := range table.AllRows() {
	//	name := row.GetCell("姓名").Value
	//	day := row.GetCell("女工日结").Value
	//	fmt.Printf("%d. %v: %v\n", i+1, name, day)
	//}
	//table.ToMarkdown("按女工日结排序")

	engine := NewYaegiEngine(table)
	for i := 0; i < table.RowCount(); i++ {
		engine.UpdateRowCell(i)
	}
	table.ToMarkdown("按公式列更新单元格数")
	//main2()//todo bug,need make uint test
}

func main2() {
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

	// 3. 显示原始数据
	fmt.Println("=== 原始数据 ===")
	printFlatTable(table)

	// 4. 按姓名分组
	fmt.Println("\n=== 按姓名分组后 ===")
	table.GroupBy("姓名")
	table.ToMarkdown("按姓名分组")

	printGroupedTable(table)

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

	// 7. 取消分组
	fmt.Println("\n=== 取消分组 ===")
	table.Ungroup()
	printFlatTable(table)
	table.ToMarkdown("取消分组")
}

func printFlatTable(table *sdk.TreeTable) {
	fmt.Printf("%-8s %-12s %-12s %-16s\n", "姓名", "女工日结", "男工车结", "女工实发工资")
	fmt.Println("────────── ──────────── ──────────── ────────────────")
	for _, row := range table.AllRows() {
		name := row.GetCell("姓名").Value
		day := row.GetCell("女工日结").Value
		car := row.GetCell("男工车结").Value
		salary := row.GetCell("女工实发工资").Value
		fmt.Printf("%-8v %-12v %-12v %-16v\n", name, day, car, salary)
	}
}

func printGroupedTable(table *sdk.TreeTable) {
	fmt.Println("树形结构:")
	for node := range table.Root.Walk() {
		indent := strings.Repeat("  ", node.Depth()-1)
		if node.IsContainer() {
			groupName := node.GroupKey
			if cell := node.GetCell("姓名"); cell != nil {
				groupName = fmt.Sprintf("%v", cell.Value)
			}
			fmt.Printf("%s📁 分组: %s (%d人)\n", indent, groupName, len(node.Children))
		} else {
			name := node.GetCell("姓名").Value
			day := node.GetCell("女工日结").Value
			fmt.Printf("%s👤 %v: %v\n", indent, name, day)
		}
	}
}
