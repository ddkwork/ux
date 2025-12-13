package main

import (
	"fmt"

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

// 添加示例节点（含公式列依赖）
func addExampleNodes() *sdk.Node {
	// 员工1 - 杨萍
	emp1 := sdk.NewNode([]sdk.CellData{
		{Name: "姓名", Value: "杨萍", Type: "text"},
		{Name: "出生年份", Value: 1990, Type: "number"},
		{Name: "年龄", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				return 2024 - int(ctx["出生年份"].(float64))
			}
		`},
		{Name: "女工日结", Value: 150.0, Type: "number"},
		{Name: "计算结果", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				name := ctx["姓名"].(string)
				nvGong := ctx["女工日结"].(float64)
				sanRenZuSum := table.SumIf("姓名", "三人组", "女工日结")
				
				if name == "杨萍" {
					return (sanRenZuSum/3.0) + nvGong
				}
				return 0.0
			}
		`},
		{Name: "入职日期", Value: "2020-03-15", Type: "date"},
		{Name: "状态", Value: "在职", Type: "select"},
	})

	// 员工2 - 房东
	emp2 := sdk.NewNode([]sdk.CellData{
		{Name: "姓名", Value: "房东", Type: "text"},
		{Name: "出生年份", Value: 1985, Type: "number"},
		{Name: "年龄", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				return 2024 - int(ctx["出生年份"].(float64))
			}
		`},
		{Name: "女工日结", Value: 200.0, Type: "number"},
		{Name: "计算结果", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				name := ctx["姓名"].(string)
				nvGong := ctx["女工日结"].(float64)
				
				if name == "房东" {
					return nvGong
				}
				return 0.0
			}
		`},
		{Name: "入职日期", Value: "2019-07-01", Type: "date"},
		{Name: "状态", Value: "在职", Type: "select"},
	})

	// 三人组
	sanRenZu := sdk.NewNode([]sdk.CellData{
		{Name: "姓名", Value: "三人组", Type: "text"},
		{Name: "出生年份", Value: 0, Type: "number"},
		{Name: "年龄", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				return 0
			}
		`},
		{Name: "女工日结", Value: 300.0, Type: "number"},
		{Name: "计算结果", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				return 0
			}
		`},
		{Name: "入职日期", Value: "", Type: "date"},
		{Name: "状态", Value: "在职", Type: "select"},
	})

	// 容器节点（部门）
	dept := sdk.NewContainerNode("部门", []sdk.CellData{
		{Name: "姓名", Value: "技术部", Type: "text"},
		{Name: "出生年份", Value: 0, Type: "number"},
		{Name: "年龄", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				return 0
			}
		`},
		{Name: "女工日结", Value: 0.0, Type: "number"},
		{Name: "计算结果", Type: "formula", Formula: `
			func(ctx map[string]interface{}, node *Node, table *TreeTable) interface{} {
				return 0
			}
		`},
		{Name: "入职日期", Value: "", Type: "date"},
		{Name: "状态", Value: "在职", Type: "select"},
	})
	dept.AddChildren([]*sdk.Node{emp1, emp2, sanRenZu})
	return dept
}

// ------------------------------ 示例用法（含公式计算演示） ------------------------------
func main() {
	// 创建表格并添加示例节点
	t := sdk.NewTreeTable()
	nodes := addExampleNodes()
	t.Root.AddChild(nodes)

	stream.WriteTruncate("tmp/1.md", t.ToMarkdown())

	// 打印Markdown（自动计算公式列值）
	//fmt.Println("=== 树形表格（含公式列）===")
	//fmt.Println(t.ToMarkdown())

	// 修改依赖列值，观察公式列自动更新
	if len(t.Root.Children) > 0 && len(t.Root.Children[0].Children) > 0 {
		empNode := t.Root.Children[0].Children[0] // 获取第一个员工节点
		empNode.SetCellValue("女工日结", 250.0, t)    // 修改女工日结

		// 手动触发公式计算
		for i := range empNode.RowCells {
			if empNode.RowCells[i].Type == "formula" {
				//t.calculateFormulaCell(empNode, &empNode.RowCells[i])
			}
		}

		//fmt.Println("\n=== 修改女工日结后 ===")
	}

	// 导出JSON（含公式定义）
	//if jsonData, err := t.ToJSON(); err == nil {
	//	fmt.Println("\n=== JSON导出（含公式）===")
	//	fmt.Println(string(jsonData))
	//}

	engine := NewYaegiEngine(t)
	for i := 0; i < t.RowCount(); i++ {
		if err := engine.CalculateRow(i); err != nil {
			fmt.Printf("行 %d 错误: %v\n", i, err)
		}
	}
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
