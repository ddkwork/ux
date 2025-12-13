package main

import (
	"fmt"
	"strings"

	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
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

func (e *YaegiEngine) CalculateRow(rowIndex int) error {
	row := e.table.GetRow(rowIndex)
	if row == nil {
		return fmt.Errorf("行不存在")
	}

	i := interp.New(interp.Options{
		GoPath:       "./",
		Unrestricted: true,
	})
	i.Use(stdlib.Symbols)

	err := i.Use(Symbols)

	if err != nil {
		return fmt.Errorf("导出失败: %v", err)
	}

	userScript := `
package main

import (
	"fmt"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
)

func RunScript(t *sdk.TreeTable, rowIndex int) {
	nameVal, ok := t.GetCellValue(rowIndex, "姓名")
	if !ok {
		return
	}
	name := fmt.Sprintf("%v", nameVal)

	nvGongVal, ok := t.GetCellValue(rowIndex, "女工日结")
	if !ok {
		return
	}
	nvGong, _ := sdk.ToFloat(nvGongVal)

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

`
	stream.WriteGoFile("tmp/user_script.go", userScript)

	_, err = i.Eval(userScript)
	if err != nil {
		return fmt.Errorf("执行脚本失败: %v", err)
	}
	runScriptVal, err := i.Eval("RunScript")
	if err != nil {
		return fmt.Errorf("获取 RunScript 失败: %v", err)
	}

	runScript := runScriptVal.Interface().(func(*sdk.TreeTable, int))
	runScript(e.table, rowIndex)

	return nil
}

//func main() {
//	table := sdk.NewTreeTable()
//
//	table.AddRow(map[string]any{"姓名": "拼车", "女工日结": 0.0})
//	table.AddRow(map[string]any{"姓名": "三人组", "女工日结": 900.0})
//	table.AddRow(map[string]any{"姓名": "房东", "女工日结": 350.0})
//	table.AddRow(map[string]any{"姓名": "杨萍", "女工日结": 200.0})
//	table.AddRow(map[string]any{"姓名": "二人组", "女工日结": 600.0})
//
//	engine := NewYaegiEngine(table)
//
//	fmt.Println("=== 计算结果 ===")
//	fmt.Printf("%-10s | %-10s | %-10s\n", "姓名", "女工日结", "计算结果")
//	fmt.Println("-----------|------------|------------")
//
//	for i := 0; i < table.RowCount(); i++ {
//		if err := engine.CalculateRow(i); err != nil {
//			fmt.Printf("行 %d 错误: %v\n", i, err)
//		}
//
//		row := table.GetRow(i)
//		result := row["计算结果"]
//		if result == nil {
//			result = 0.0
//		}
//		fmt.Printf("%-10s | %-10.0f | %-10.2f\n",
//			row["姓名"], row["女工日结"], result)
//	}
//
//	sanRenZuSum := table.SumIf("姓名", "三人组", "女工日结")
//	fmt.Printf("\n验证: 三人组总和=%.0f\n", sanRenZuSum)
//	fmt.Printf("杨萍应得: %.0f/3 + 200 = %.0f\n", sanRenZuSum, sanRenZuSum/3+200)
//	fmt.Printf("二人组应得: %.0f/3 + 600/2 = %.0f\n", sanRenZuSum, sanRenZuSum/3+300)
//}
 
func main() {
	// 1. 创建表格
	table := sdk.NewTreeTable()

	// 2. 直观的表格数据定义
	data := sdk.TableData{
		Columns: []sdk.ColumnConfig{
			{Name: "姓名", Type: sdk.FieldTypeSingleLineText},
			{Name: "女工日结", Type: sdk.FieldTypeNumber},
			{Name: "男工车结", Type: sdk.FieldTypeNumber},
			{Name: "女工实发工资", Type: sdk.FieldTypeFormula, Formula: "{{女工日结}} * 0.8 + {{男工车结}} * 0.5"},
		},
		Rows: [][]any{
			{"三人组", 2966.30, 1104.20, 0},
			{"房东", 442.40, 196.80, 442.4},
			{"二人组", 5913.60, 2248.60, 3945.566666667},
			{"杨萍", 3744.90, 1465.20, 4733.666666667},
			{"拼车", 406.90, 175.00, 0},
		},
	}

	// 3. 一键设置数据
	table.LoadTableData(data)

	// 4. 显示数据
	fmt.Println("=== 基础数据展示 ===")
	fmt.Printf("%-8s %-12s %-12s %-16s\n", "姓名", "女工日结", "男工车结", "女工实发工资")
	fmt.Println("────────── ──────────── ──────────── ────────────────")

	for _, row := range table.AllRows() {
		name := row.GetCell("姓名", table).Value
		day := row.GetCell("女工日结", table).Value
		car := row.GetCell("男工车结", table).Value
		salary := row.GetCell("女工实发工资", table).Value
		fmt.Printf("%-8v %-12v %-12v %-16v\n", name, day, car, salary)
	}

	// 5. 排序演示
	fmt.Println("\n=== 按女工日结降序排序 ===")
	table.SortByColumn("女工日结", false)
	for i, row := range table.AllRows() {
		name := row.GetCell("姓名", table).Value
		day := row.GetCell("女工日结", table).Value
		fmt.Printf("%d. %v: %v\n", i+1, name, day)
	}
}

func main2() {
	// 1. 创建表格
	table := sdk.NewTreeTable()

	// 2. 设置数据（包含重复姓名用于分组）
	data := sdk.TableData{
		Columns: []sdk.ColumnConfig{
			{Name: "姓名", Type: sdk.FieldTypeSingleLineText},
			{Name: "女工日结", Type: sdk.FieldTypeNumber},
			{Name: "男工车结", Type: sdk.FieldTypeNumber},
			{Name: "女工实发工资", Type: sdk.FieldTypeNumber},
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

	// 3. 显示原始数据
	fmt.Println("=== 原始数据 ===")
	printFlatTable(table)

	// 4. 按姓名分组
	fmt.Println("\n=== 按姓名分组后 ===")
	table.GroupBy("姓名")
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
		if cell := group.GetCell("姓名", table); cell != nil {
			groupName = fmt.Sprintf("%v", cell.Value)
		}
		fmt.Printf("\n📁 分组: %s (%d人)\n", groupName, len(group.Children))

		for _, member := range group.Children {
			if name := member.GetCell("姓名", table); name != nil {
				if day := member.GetCell("女工日结", table); day != nil {
					fmt.Printf("  👤 %v: %v\n", name.Value, day.Value)
				}
			}
		}
	}

	// 7. 取消分组
	fmt.Println("\n=== 取消分组 ===")
	table.Ungroup()
	printFlatTable(table)
}

func printFlatTable(table *sdk.TreeTable) {
	fmt.Printf("%-8s %-12s %-12s %-16s\n", "姓名", "女工日结", "男工车结", "女工实发工资")
	fmt.Println("────────── ──────────── ──────────── ────────────────")
	for _, row := range table.AllRows() {
		name := row.GetCell("姓名", table).Value
		day := row.GetCell("女工日结", table).Value
		car := row.GetCell("男工车结", table).Value
		salary := row.GetCell("女工实发工资", table).Value
		fmt.Printf("%-8v %-12v %-12v %-16v\n", name, day, car, salary)
	}
}

func printGroupedTable(table *sdk.TreeTable) {
	fmt.Println("树形结构:")
	for node := range table.Root.Walk() {
		indent := strings.Repeat("  ", node.Depth()-1)
		if node.IsContainer() {
			groupName := node.GroupKey
			if cell := node.GetCell("姓名", table); cell != nil {
				groupName = fmt.Sprintf("%v", cell.Value)
			}
			fmt.Printf("%s📁 分组: %s (%d人)\n", indent, groupName, len(node.Children))
		} else {
			name := node.GetCell("姓名", table).Value
			day := node.GetCell("女工日结", table).Value
			fmt.Printf("%s👤 %v: %v\n", indent, name, day)
		}
	}
}
