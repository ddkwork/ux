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
	table  *sdk.MemoryTable
}

func NewYaegiEngine(table *sdk.MemoryTable) *YaegiEngine {
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

import "github.com/ddkwork/ux/demo/erp/gongshi/sdk"

func RunScript(t *sdk.MemoryTable, index int) {
	name := t.GetRow(index)["姓名"].(string)
	nvGong := t.GetRow(index)["女工日结"].(float64)
	sanRenZuSum := t.SumIf("姓名", "三人组", "女工日结")

	switch name {
	case "拼车", "三人组":
		t.SetValue(index, "计算结果", 0.0)
	case "房东":
		t.SetValue(index, "计算结果", nvGong)
	case "杨萍":
		t.SetValue(index, "计算结果", (sanRenZuSum/3.0)+nvGong)
	case "二人组":
		t.SetValue(index, "计算结果", (sanRenZuSum/3.0)+(nvGong/2.0))
	default:
		t.SetValue(index, "计算结果", 0.0)
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

	runScript := runScriptVal.Interface().(func(*sdk.MemoryTable, int))
	runScript(e.table, rowIndex)

	return nil
}

func main() {
	table := sdk.NewMemoryTable()

	table.AddRow(map[string]interface{}{"姓名": "拼车", "女工日结": 0.0})
	table.AddRow(map[string]interface{}{"姓名": "三人组", "女工日结": 900.0})
	table.AddRow(map[string]interface{}{"姓名": "房东", "女工日结": 350.0})
	table.AddRow(map[string]interface{}{"姓名": "杨萍", "女工日结": 200.0})
	table.AddRow(map[string]interface{}{"姓名": "二人组", "女工日结": 600.0})

	engine := NewYaegiEngine(table)

	fmt.Println("=== 计算结果 ===")
	fmt.Printf("%-10s | %-10s | %-10s\n", "姓名", "女工日结", "计算结果")
	fmt.Println("-----------|------------|------------")

	for i := 0; i < table.RowCount(); i++ {
		if err := engine.CalculateRow(i); err != nil {
			fmt.Printf("行 %d 错误: %v\n", i, err)
		}

		row := table.GetRow(i)
		result := row["计算结果"]
		if result == nil {
			result = 0.0
		}
		fmt.Printf("%-10s | %-10.0f | %-10.2f\n",
			row["姓名"], row["女工日结"], result)
	}

	sanRenZuSum := table.SumIf("姓名", "三人组", "女工日结")
	fmt.Printf("\n验证: 三人组总和=%.0f\n", sanRenZuSum)
	fmt.Printf("杨萍应得: %.0f/3 + 200 = %.0f\n", sanRenZuSum, sanRenZuSum/3+200)
	fmt.Printf("二人组应得: %.0f/3 + 600/2 = %.0f\n", sanRenZuSum, sanRenZuSum/3+300)
}
